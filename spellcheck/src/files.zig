const std = @import("std");
const bf = @import("./bloom_filter.zig");

const FileError = error{
    AllocFailed,
    ReadError,
};

const sample_ratio = 100;
const sample_limit = 2048;

/// Computes the approximate number of words in a `File` by reading a sample of the file.
///
/// To open a file, you may use `std.fs.cwd().openFile(path: []const u8)`.
pub fn approx_word_count(allocator: std.mem.Allocator, file: std.fs.File) FileError!u64 {
    const stats = file.stat() catch {
        return FileError.ReadError;
    };
    const size = stats.size;
    const sample_size = @min(@min(@max(size / sample_ratio, sample_ratio), sample_limit), size);

    if (size == 0) {
        return 0;
    }

    const sample = allocator.alloc(u8, sample_size) catch {
        return FileError.AllocFailed;
    };
    defer allocator.free(sample);

    _ = file.reader().readAtLeast(sample, sample_size) catch {
        return FileError.ReadError;
    };

    var sample_count: u64 = 0;
    for (sample) |byte| {
        if (byte == '\n') {
            sample_count += 1;
        }
    }

    const approx_count = sample_count * size / sample_size;

    return approx_count;
}

/// Builds the dictionary from the source file to a destination file.
pub fn build(allocator: std.mem.Allocator, source: std.fs.File, _: std.fs.File) (FileError || bf.BloomFilterError)!void {
    const wc = try approx_word_count(allocator, source);

    // reset file pointer
    source.seekTo(0) catch {
        return FileError.ReadError;
    };

    const source_reader = source.reader();

    var bloom_filter = try bf.from_reader(allocator, wc, source_reader);
    defer bloom_filter.deinit();

    // TODO: remove
    var m: u64 = 0;
    while (m < bloom_filter.m) : (m += 1) {
        std.debug.print("{}: {} ", .{ m, bloom_filter.bit_set.isSet(m) });
    }
}

test "approx_word_count" {
    // ugly hard-coding of a test file
    const file = try std.fs.cwd().openFile("dict.txt", .{});

    // const tmp_dir = std.testing.tmpDir(.{});
    // const file = try tmp_dir.dir.createFile("dict.txt", .{ .read = true });
    // _ = try file.write("word\nword\nword\n");

    const count = approx_word_count(std.testing.allocator, file) catch |err| {
        return err;
    };

    try std.testing.expectEqual(272768, count);
}

test "build" {
    const tmp_dir = std.testing.tmpDir(.{});

    const dict_file = try tmp_dir.dir.createFile("dict.txt", .{});
    try dict_file.writer().writeAll("these\nare\nwords\n");
    try dict_file.sync();
    dict_file.close();

    const source = try tmp_dir.dir.openFile("dict.txt", .{});

    const dest = try tmp_dir.dir.createFile("dict.sc", .{});

    try build(std.testing.allocator, source, dest);
}
