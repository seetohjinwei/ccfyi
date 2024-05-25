const std = @import("std");
const args = @import("./args.zig");

pub fn main() !void {
    // this is a CLI program :)
    // https://ziglang.org/documentation/master/#toc-Choosing-an-Allocator

    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();

    const allocator = arena.allocator();

    const arguments = args.parse(allocator) catch {
        std.process.exit(1);
    };

    std.debug.print("{s}\n", .{arguments.program_name});
}
