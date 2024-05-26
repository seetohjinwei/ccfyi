const std = @import("std");
const args = @import("./args.zig");
const files = @import("./files.zig");

const default_dict_file = "dict.sc";

fn display_misspelled_words(allocator: std.mem.Allocator, words: std.ArrayList([]const u8)) !void {
    if (words.items.len == 0) {
        return;
    }

    const stdout = std.io.getStdOut();

    try stdout.writeAll("Misspelt:\n");

    for (words.items) |word| {
        const f = try std.fmt.allocPrint(allocator, "  {s}\n", .{word});
        try stdout.writeAll(f);
    }
}

pub fn main() !void {
    // this is a CLI program :)
    // https://ziglang.org/documentation/master/#toc-Choosing-an-Allocator

    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();

    const allocator = arena.allocator();

    var arguments = args.parse(allocator) catch {
        std.process.exit(1);
    };
    defer arguments.deinit();

    const dict_path = arguments.dict_path orelse default_dict_file;

    if (arguments.build_path != null) {
        const source = try std.fs.cwd().openFile(arguments.build_path.?, .{});
        defer source.close();
        const dest = try std.fs.cwd().createFile(dict_path, .{});
        defer dest.close();
        try files.build(allocator, source, dest);
    }

    const dict_file = try std.fs.cwd().openFile(dict_path, .{});
    var bloom_filter = try files.read_dict(allocator, dict_file);

    // TODO: accept piped input too
    const misspelled_words = try bloom_filter.has_many(allocator, arguments.words);
    try display_misspelled_words(allocator, misspelled_words);
    defer misspelled_words.deinit();
}
