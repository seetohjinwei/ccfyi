const std = @import("std");
const files = @import("./files.zig");

const IDENTIFIER = "CCBF";
const VERSION = "01";

pub const Error = error{
    AllocFailed,
    HashFailed,
    WordIsTooLong,
    InvalidFormat,
};

/// BloomFilter is an implementation of a standard bloom filter.
///
/// `m`: size of bitset
/// `k`: number of hash functions
pub const BloomFilter = struct {
    allocator: std.mem.Allocator,
    bit_set: std.bit_set.DynamicBitSetUnmanaged,
    m: u64,
    k: u64,

    /// Initialises the BloomFilter.
    fn init(allocator: std.mem.Allocator, m: u64, k: u64) Error!BloomFilter {
        const bit_set = std.bit_set.DynamicBitSetUnmanaged.initEmpty(allocator, m) catch {
            return Error.AllocFailed;
        };

        return BloomFilter{ .allocator = allocator, .bit_set = bit_set, .m = m, .k = k };
    }

    /// Adds `data` to the BloomFilter.
    fn add(self: *BloomFilter, data: []const u8) Error!void {
        var k: u64 = 0;
        while (k < self.k) : (k += 1) {
            const index = hash(self.allocator, self.m, k, data) catch {
                return Error.HashFailed;
            };
            self.bit_set.set(index);
        }
    }

    /// Returns true if `data` is probably in the BloomFilter.
    /// Returns false if `data` is definitely not in the BloomFilter.
    pub fn has(self: *BloomFilter, data: []const u8) Error!bool {
        var k: u64 = 0;
        while (k < self.k) : (k += 1) {
            const index = hash(self.allocator, self.m, k, data) catch {
                return Error.HashFailed;
            };
            if (!self.bit_set.isSet(index)) {
                return false;
            }
        }

        return true;
    }

    /// Returns the list of words that are definitely not in the BloomFilter.
    /// The caller owns the returned memory. Free it with `data.deinit()`.
    pub fn has_many(self: *BloomFilter, allocator: std.mem.Allocator, data: std.ArrayList([]const u8)) Error!std.ArrayList([]const u8) {
        var no_match = std.ArrayList([]const u8).init(allocator);

        for (data.items) |item| {
            if (!try self.has(item)) {
                no_match.append(item) catch {
                    return Error.AllocFailed;
                };
            }
        }

        return no_match;
    }

    fn bit_set_to_bytes(self: *BloomFilter, bytes: *std.ArrayList(u8)) Error!void {
        // The actual bytes of the bitset.
        const num_bytes = std.math.divCeil(u64, self.m, 8) catch unreachable;
        var byte_index: usize = 0;
        while (byte_index < num_bytes) : (byte_index += 1) {
            var byte: u8 = 0;

            // doesn't compile without `inline` lol
            inline for (0..8) |bit_index| {
                const bit_set_pos = byte_index * 8 + bit_index;
                if (bit_set_pos >= self.bit_set.capacity()) {
                    break;
                }

                if (self.bit_set.isSet(bit_set_pos)) {
                    byte |= @as(u8, 1) << bit_index;
                }
            }

            bytes.append(byte) catch {
                return Error.AllocFailed;
            };
        }
    }

    /// Returns bytes reprsenting the BloomFilter. The BloomFilter can be recreated using these bytes.
    /// The caller owns the returned memory. Free it with `allocator.free(bytes)`.
    pub fn to_bytes(self: *BloomFilter) Error![]u8 {
        // alloc on stack (and reuse this)
        var buf = [_]u8{0} ** 4;
        const buf_u16 = buf[0..@sizeOf(u16)];

        var bytes = std.ArrayList(u8).init(self.allocator);
        defer bytes.deinit();

        // The first four bytes will be an identifier, we’ll use CCBF.
        bytes.appendSlice(IDENTIFIER) catch {
            return Error.AllocFailed;
        };
        // The next two bytes will be a version number to describe the version number of the file.
        bytes.appendSlice(VERSION) catch {
            return Error.AllocFailed;
        };
        // The next two bytes will be the number of hash functions used.
        std.mem.writeInt(u16, buf_u16, @intCast(self.k), std.builtin.Endian.big);
        bytes.appendSlice(buf_u16) catch {
            return Error.AllocFailed;
        };
        // The next four bytes will be the number of bits used for the filter.
        std.mem.writeInt(u32, &buf, @intCast(self.m), std.builtin.Endian.big);
        bytes.appendSlice(&buf) catch {
            return Error.AllocFailed;
        };

        // The actual bytes of the bitset.
        try self.bit_set_to_bytes(&bytes);

        const data = bytes.toOwnedSlice() catch {
            return Error.AllocFailed;
        };

        return data;
    }

    /// Deinits the BloomFilter.
    pub fn deinit(self: *BloomFilter) void {
        self.bit_set.deinit(self.allocator);
    }
};

fn hash(allocator: std.mem.Allocator, limit: u64, index: u64, data: []const u8) !u64 {
    var key = [_]u8{0} ** 8;
    std.mem.writeInt(u64, &key, index, std.builtin.Endian.big);

    const buf = try allocator.alloc(u8, key.len + data.len);
    @memcpy(buf[0..key.len], key[0..]);
    @memcpy(buf[key.len..], data);
    defer allocator.free(buf);

    return std.hash.Fnv1a_64.hash(buf) % limit;
}

const p_default = 0.0001;

/// Computes `m` and `k` from `n` and `p`..
/// formulas from: https://hur.st/bloomfilter
fn get_m_k(n: u64, p: f64) struct { m: u64, k: u64 } {
    const n_float = @as(f64, @floatFromInt(n));

    const m = @ceil((n_float * @log(p)) / @log(1 / std.math.pow(f64, 2, @log(2.0))));
    const k = @round((m / n_float) * @log(2.0));

    return .{ .m = @intFromFloat(m), .k = @intFromFloat(k) };
}

/// Creates a BloomFilter from a Reader.
/// `approx_word_count` is an approximation of the number of words in the reader.
/// To generate `approx_word_count`, you may use `files.approx_word_count`.
pub fn from_reader(allocator: std.mem.Allocator, approx_word_count: u64, reader: anytype) Error!BloomFilter {
    const m_k = get_m_k(approx_word_count, p_default);

    var bloom_filter = try BloomFilter.init(allocator, m_k.m, m_k.k);

    while (reader.readUntilDelimiterOrEofAlloc(allocator, '\n', 1000) catch {
        return Error.WordIsTooLong;
    }) |word| {
        try bloom_filter.add(word);
        allocator.free(word);
    }

    return bloom_filter;
}

pub fn from_sc(allocator: std.mem.Allocator, reader: anytype) Error!BloomFilter {
    var buf = [_]u8{0} ** 4;
    const buf_2 = buf[0..2];
    const buf_4 = buf[0..4];

    // The first four bytes will be an identifier, we’ll use CCBF.
    const identifier_size = reader.readAll(buf_4) catch {
        return Error.InvalidFormat;
    };
    if (identifier_size != 4) {
        return Error.InvalidFormat;
    }
    if (!std.mem.eql(u8, buf_4, IDENTIFIER)) {
        return Error.InvalidFormat;
    }
    // The next two bytes will be a version number to describe the version number of the file.
    const version_size = reader.readAll(buf_2) catch {
        return Error.InvalidFormat;
    };
    if (version_size != 2) {
        return Error.InvalidFormat;
    }
    if (!std.mem.eql(u8, buf_2, VERSION)) {
        return Error.InvalidFormat;
    }
    // The next two bytes will be the number of hash functions used.
    const k_size = reader.readAll(buf_2) catch {
        return Error.InvalidFormat;
    };
    if (k_size != 2) {
        return Error.InvalidFormat;
    }
    const k = std.mem.readInt(u16, buf_2, std.builtin.Endian.big);
    // The next four bytes will be the number of bits used for the filter.
    const m_size = reader.readAll(buf_4) catch {
        return Error.InvalidFormat;
    };
    if (m_size != 4) {
        return Error.InvalidFormat;
    }
    const m = std.mem.readInt(u32, buf_4, std.builtin.Endian.big);

    var bloom_filter = try BloomFilter.init(allocator, m, k);

    // The actual bytes of the bitset.
    const num_bytes = std.math.divCeil(u64, m, 8) catch unreachable;
    // std.debug.print("num_bytes={}\n", .{num_bytes});

    for (0..num_bytes) |byte_index| {
        const byte = reader.readByte() catch {
            // std.debug.print("ERR={}\n", .{err});
            return Error.InvalidFormat;
        };

        // std.debug.print("byte {}={b}\n", .{ byte_index, byte });

        inline for (0..8) |bit_index| {
            const bit_set_pos = byte_index * 8 + bit_index;
            if (bit_set_pos >= m) {
                break;
            }

            // std.debug.print("bit {}={b}\n", .{ bit_index, byte & (@as(u8, 1) << bit_index) });

            if ((byte & (@as(u8, 1) << bit_index)) > 0) {
                bloom_filter.bit_set.set(bit_set_pos);
            }
        }
    }

    return bloom_filter;
}

test "hash" {
    try std.testing.expectEqual(65, try hash(std.testing.allocator, 100, 1, "key"));
    try std.testing.expectEqual(36000358, try hash(std.testing.allocator, 42141512, 2, "key"));
}

test "get_m_k" {
    const cases = [_]struct {
        n: u64,
        p: f64,
        m: u64,
        k: u64,
    }{
        .{ .n = 240000, .p = p_default, .m = 4600829, .k = 13 },
        .{ .n = 134, .p = p_default, .m = 2569, .k = 13 },
        .{ .n = 240000, .p = 1.0 / 100.0, .m = 2300415, .k = 7 },
        .{ .n = 134, .p = 1.0 / 12.0, .m = 694, .k = 4 },
    };

    for (cases) |case| {
        const res = get_m_k(case.n, case.p);
        try std.testing.expectEqual(case.m, res.m);
        try std.testing.expectEqual(case.k, res.k);
    }
}

test "bloom_filter" {
    const data = "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n";
    var stream = std.io.fixedBufferStream(data);
    const reader = stream.reader();

    // const file = try std.fs.cwd().openFile("dict.txt", .{});
    // const reader = file.reader();

    var bloom_filter = from_reader(std.testing.allocator, 10, reader) catch |err| {
        std.debug.print("error calling from_reader {}\n", .{err});
        return err;
    };
    defer bloom_filter.deinit();

    try std.testing.expect(try bloom_filter.has("b"));
    try std.testing.expectEqual(false, try bloom_filter.has("34"));
    try std.testing.expectEqual(false, try bloom_filter.has("jrk"));
    try std.testing.expectEqual(false, try bloom_filter.has("421"));
}

test "bloom_filter.bit_set_to_bytes" {
    var bloom_filter = try BloomFilter.init(std.testing.allocator, 10, 1);
    defer bloom_filter.deinit();

    const indexes = [_]usize{ 0, 2, 5, 6, 9 };

    for (indexes) |i| {
        bloom_filter.bit_set.set(i);
    }

    var bytes = std.ArrayList(u8).init(std.testing.allocator);
    defer bytes.deinit();

    try bloom_filter.bit_set_to_bytes(&bytes);

    try std.testing.expectEqual(2, bytes.items.len);
    try std.testing.expectEqual(0b10, bytes.items[1]);
    try std.testing.expectEqual(0b01100101, bytes.items[0]);
}

test "from_sc" {
    const data = [_]u8{ 'C', 'C', 'B', 'F', '0', '1', 0, 0x5, 0, 0, 0, 0xa } ++ [2]u8{ 0b01001010, 0x00 };
    // const data = "CCBF010AA493";
    var stream = std.io.fixedBufferStream(&data);
    const reader = stream.reader();

    var bloom_filter = try from_sc(std.testing.allocator, reader);
    defer bloom_filter.deinit();

    const expected = [10]bool{ false, true, false, true, false, false, true, false, false, false };

    for (expected, 0..10) |is_set, i| {
        try std.testing.expectEqual(is_set, bloom_filter.bit_set.isSet(i));
    }
}
