const std = @import("std");

const help_text =
    \\Usage:
    \\  {s} [-build <dict.txt>] [-dict <dict.sc>] <words> ...
    \\  cat <file> | {s} [-build <dict.txt>] [-dict <dict.sc>]
    \\
;

const Error = error{
    AllocFailed,
    MissingArgument,
};

const Args = struct {
    program_name: []const u8,
    build_path: ?[]const u8,
    dict_path: ?[]const u8 = "dict.sc",
    words: std.ArrayList([]const u8),

    pub fn deinit(self: *Args) void {
        self.words.deinit();
    }
};

/// Parses the program's arguments into an `Args` struct.
///
/// The caller is responsible for calling `args.deinit()`.
pub fn parse(allocator: std.mem.Allocator) Error!Args {
    const words = std.ArrayList([]const u8).init(allocator);

    var it = std.process.args();

    var ret = Args{
        .program_name = it.next().?,
        .build_path = null,
        .dict_path = null,
        .words = words,
    };

    var has_non_flag = false;

    while (it.next()) |arg| {
        if (!has_non_flag and std.mem.eql(u8, arg, "-build")) {
            const next_arg = it.next();
            if (next_arg == null) {
                return Error.MissingArgument;
            }

            ret.build_path = next_arg.?;
        } else if (!has_non_flag and std.mem.eql(u8, arg, "-dict")) {
            const next_arg = it.next();
            if (next_arg == null) {
                const p = std.fmt.allocPrint(allocator, help_text, .{ ret.program_name, ret.program_name }) catch {
                    return Error.AllocFailed;
                };
                defer allocator.free(p);
                std.io.getStdErr().writeAll(p) catch {
                    return Error.AllocFailed;
                };
                return Error.MissingArgument;
            }

            ret.dict_path = next_arg.?;
        } else {
            has_non_flag = true;

            // an argument can possibly be multiple arguments
            var arg_it = std.mem.split(u8, arg, " ");
            while (arg_it.next()) |x| {
                ret.words.append(x) catch {
                    return Error.AllocFailed;
                };
            }
        }
    }

    return ret;
}
