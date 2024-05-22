const std = @import("std");

const sample_ratio = 100;
const sample_limit = 2048;

/// Computes the approximate number of words in a `File` by reading a sample of the file.
///
/// To open a file, you may use `std.fs.cwd().openFile(path: []const u8)`.
pub fn approx_word_count(allocator: std.mem.Allocator, file: std.fs.File) !u64 {
    // TODO: https://ziglang.org/documentation/master/#toc-Choosing-an-Allocator

    const stats = try file.stat();
    const size = stats.size;
    const sample_size = @min(@max(size / sample_ratio, sample_ratio), sample_limit);

    if (size == 0) {
        return 0;
    }

    const sample = try allocator.alloc(u8, sample_size);
    defer allocator.free(sample);

    _ = try file.reader().readAtLeast(sample, sample_size);

    var sample_count: u64 = 0;
    for (sample) |byte| {
        // for (try sample.toOwnedSlice()) |byte| {
        if (byte == '\n') {
            sample_count += 1;
        }
    }

    const approx_count = sample_count * size / sample_size;

    return approx_count;
}

test "approx_word_count" {
    // ugly hard-coding of a test file
    const file = try std.fs.cwd().openFile("dict.txt", .{});

    // const tmp_dir = std.testing.tmpDir(.{});
    // const file = try tmp_dir.dir.createFile("dict.txt", .{ .read = true });
    // _ = try file.write("word\nword\nword\n");

    const count = approx_word_count(std.testing.allocator, file) catch |err| {
        std.debug.print("error calling approx_word_count {}\n", .{err});
        return err;
    };

    try std.testing.expectEqual(272768, count);
}
