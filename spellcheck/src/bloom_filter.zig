const std = @import("std");
const files = @import("./files.zig");

const BloomFilter = struct {
    m: u64,
    k: u64,

    const Writer = std.io.Writer(
        *BloomFilter,
        error{},
        write,
    );

    /// Initialises the BloomFilter.
    /// This **must** be done before writing to it.
    fn init(self: *BloomFilter) void {
        // TODO:
        self;

        // creates hash functions
        // allocates bitarray
    }

    fn write(self: *BloomFilter, data: []const u8) error{}!usize {
        // TODO:

        // throws error if `init` has not been called

        self;
        data;
        return 0;
    }

    fn writer(self: *BloomFilter) Writer {
        return .{ .context = self };
    }
};

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
pub fn from_reader(approx_word_count: u64, reader: std.io.AnyReader) BloomFilter {
    const m_k = get_m_k(approx_word_count, p_default);

    const bloom_filter = BloomFilter{ .m = m_k.m, .k = m_k.k };

    // TODO: does this only write the first line???
    reader.streamUntilDelimiter(bloom_filter.writer(), "\n", null);

    return bloom_filter;
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
