const std = @import("std");
const files = @import("./files.zig");

const VERSION = "01";

pub const BloomFilterError = error{
    AllocFailed,
    HashFailed,
    WordIsTooLong,
};

/// BloomFilter is an implementation of a standard bloom filter.
///
/// `m`: size of bitset
/// `k`: number of hash functions
const BloomFilter = struct {
    allocator: std.mem.Allocator,
    bit_set: std.bit_set.DynamicBitSetUnmanaged,
    m: u64,
    k: u64,

    /// Initialises the BloomFilter.
    fn init(allocator: std.mem.Allocator, m: u64, k: u64) BloomFilterError!BloomFilter {
        const bit_set = std.bit_set.DynamicBitSetUnmanaged.initEmpty(allocator, m) catch {
            return BloomFilterError.AllocFailed;
        };

        return BloomFilter{ .allocator = allocator, .bit_set = bit_set, .m = m, .k = k };
    }

    /// Adds `data` to the BloomFilter.
    fn add(self: *BloomFilter, data: []const u8) BloomFilterError!void {
        var k: u64 = 0;
        while (k < self.k) : (k += 1) {
            const index = hash(self.allocator, self.m, k, data) catch {
                return BloomFilterError.HashFailed;
            };
            self.bit_set.set(index);
        }
    }

    /// Returns true if `data` is probably in the BloomFilter.
    /// Returns false if `data` is definitely not in the BloomFilter.
    fn has(self: *BloomFilter, data: []const u8) BloomFilterError!bool {
        var k: u64 = 0;
        while (k < self.k) : (k += 1) {
            const index = hash(self.allocator, self.m, k, data) catch {
                return BloomFilterError.HashFailed;
            };
            if (!self.bit_set.isSet(index)) {
                return false;
            }
        }

        return true;
    }

    /// Returns bytes reprsenting the BloomFilter. The BloomFilter can be recreated using these bytes.
    /// The caller owns the returned memory. Free it with `allocator.free(bytes)`.
    pub fn to_bytes(self: *BloomFilter) BloomFilterError![]u8 {
        // alloc on stack (and reuse this)
        var buf = [_]u8{ 0, 0, 0, 0 };
        const buf_u16 = buf[0..@sizeOf(u16)];

        var bytes = std.ArrayList(u8).init(self.allocator);

        // The first four bytes will be an identifier, we’ll use CCBF.
        bytes.appendSlice("CCBF") catch {
            return BloomFilterError.AllocFailed;
        };
        // The next two bytes will be a version number to describe the version number of the file.
        bytes.appendSlice(VERSION) catch {
            return BloomFilterError.AllocFailed;
        };
        // The next two bytes will be the number of hash functions used.
        std.mem.writeInt(u16, buf_u16, @intCast(self.k), std.builtin.Endian.big);
        bytes.appendSlice(buf_u16) catch {
            return BloomFilterError.AllocFailed;
        };
        // The next four bytes will be the number of bits used for the filter.
        std.mem.writeInt(u32, &buf, @intCast(self.m), std.builtin.Endian.big);
        bytes.appendSlice(&buf) catch {
            return BloomFilterError.AllocFailed;
        };

        // TODO: write the bitset -> bytes conversion here
        // probably just loop through the bits?

        defer bytes.deinit();

        const data = bytes.toOwnedSlice() catch {
            return BloomFilterError.AllocFailed;
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
pub fn from_reader(allocator: std.mem.Allocator, approx_word_count: u64, reader: anytype) BloomFilterError!BloomFilter {
    const m_k = get_m_k(approx_word_count, p_default);

    var bloom_filter = try BloomFilter.init(allocator, m_k.m, m_k.k);

    while (reader.readUntilDelimiterOrEofAlloc(allocator, '\n', 1000) catch {
        return BloomFilterError.WordIsTooLong;
    }) |word| {
        try bloom_filter.add(word);
        allocator.free(word);
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

    // var m: u64 = 0;
    // while (m < bloom_filter.m) : (m += 1) {
    //     std.debug.print("{}: {} ", .{ m, bloom_filter.bit_set.isSet(m) });
    // }

    try std.testing.expect(try bloom_filter.has("b"));
    try std.testing.expectEqual(false, try bloom_filter.has("34"));
    try std.testing.expectEqual(false, try bloom_filter.has("jrk"));
    try std.testing.expectEqual(false, try bloom_filter.has("421"));
}
