const std = @import("std");
const mvzr = @import("mvzr");
const Allocator = std.mem.Allocator;
const print = std.debug.print;
const parseInt = std.fmt.parseInt;
const eql = std.mem.eql;

pub fn main() !void {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    var args = try std.process.ArgIterator.initWithAllocator(allocator);
    _ = args.next();
    const filename = args.next().?;
    const contents = try readFile(allocator, filename);
    parse(contents);
}

fn readFile(allocator: Allocator, filename: [:0]const u8) ![]u8 {
    const file = try std.fs.cwd().openFile(filename, .{});
    defer file.close();

    const file_size = try file.getEndPos();
    const data = try allocator.alloc(u8, file_size);

    _ = try file.readAll(data);
    return data;
}

fn parse(haystack: []const u8) void {
    var one: u32 = 0;
    var two: u32 = 0;
    var do = true;

    const reNum = mvzr.compile("\\d+").?;

    const patt = "mul\\(\\d{1,3},\\d{1,3}\\)|do\\(\\)|don't\\(\\)";
    const re = mvzr.compile(patt).?;
    var matches = re.iterator(haystack);
    while (matches.next()) |m| {
        const token = m.slice;
        if (eql(u8, token, "do()")) {
            do = true;
        } else if (eql(u8, token, "don't()")) {
            do = false;
        } else {
            var prod: u32 = 1;
            var nums = reNum.iterator(token);
            while (nums.next()) |n| {
                prod *= parseInt(u32, n.slice, 10) catch 1;
            }
            one += prod;
            if (do) {
                two += prod;
            }
        }
    }

    print("Part 1: {}\n", .{one});
    print("Part 2: {}\n", .{two});
}
