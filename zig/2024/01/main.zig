const std = @import("std");
const mem = std.mem;
const Allocator = mem.Allocator;
const ArrayList = std.ArrayList;

const IntIterator = struct {
    buffer: []const u8,
    index: ?usize,

    fn init(b: []const u8) IntIterator {
        return .{ .buffer = b, .index = 0 };
    }

    fn next(self: *IntIterator) ?i64 {
        const start = self.index orelse return null;
        var i = start;
        var found = false;
        var num: i64 = 0;
        var neg: i8 = 1;
        for (self.buffer[i..]) |char| {
            i += 1;
            if (char >= '0' and char <= '9') {
                if (!found and i >= start + 2 and self.buffer[i - 2] == '-') {
                    neg = -1;
                }
                found = true;
                num = num * 10 + char - '0';
            } else if (found) {
                break;
            }
        }
        self.index = if (i >= self.buffer.len) null else i;
        return if (found) neg * num else null;
    }
};

const String = struct {
    allocator: Allocator,
    string: []const u8,

    pub fn init(a: Allocator, s: []const u8) !String {
        const d = try a.alloc(u8, s.len);
        @memcpy(d, s);
        return String{ .allocator = a, .string = d };
    }

    pub fn deinit(self: *String) void {
        self.allocator.free(self.string);
    }

    pub fn clone(self: *String) !String {
        return String.init(self.allocator, self.string);
    }

    pub fn format(value: String, comptime fmt: []const u8, options: std.fmt.FormatOptions, writer: anytype) !void {
        _ = fmt;
        _ = options;
        try std.fmt.format(writer, "{s}", .{value.string});
    }
};

const Lines = struct {
    data: ArrayList(String),

    fn deinit(self: *Lines) void {
        var data = self.data;
        while (data.items.len > 0) {
            var str = data.pop();
            str.deinit();
        }
        data.deinit();
    }
};

fn readFile(allocator: Allocator, filename: []const u8) !Lines {
    const file = try std.fs.cwd().openFile(filename, .{});
    defer file.close();

    const file_size = try file.getEndPos();
    const data = try allocator.alloc(u8, file_size);
    defer allocator.free(data);

    _ = try file.readAll(data);

    var out = ArrayList(String).init(allocator);
    var it = mem.splitScalar(u8, data, '\n');
    while (it.next()) |part| {
        const trimmed = mem.trim(u8, part, "\t\r ");
        const str = try String.init(allocator, trimmed);
        try out.append(str);
    }
    return Lines{ .data = out };
}

const Puzzle = struct {
    allocator: Allocator,
    lists: [2]ArrayList(i32),

    fn init(a: Allocator) Puzzle {
        const l = [2]ArrayList(i32){
            ArrayList(i32).init(a),
            ArrayList(i32).init(a),
        };
        return .{
            .allocator = a,
            .lists = l,
        };
    }

    fn deinit(self: *Puzzle) void {
        self.lists[0].deinit();
        self.lists[1].deinit();
    }

    fn load(self: *Puzzle, filename: []const u8) !void {
        var lines = try readFile(self.allocator, filename);
        defer lines.deinit();

        for (lines.data.items) |item| {
            var it = IntIterator.init(item.string);
            const item0: i32 = @intCast(it.next().?);
            const item1: i32 = @intCast(it.next().?);
            try self.lists[0].append(item0);
            try self.lists[1].append(item1);
        }

        mem.sort(i32, self.lists[0].items, {}, comptime std.sort.asc(i32));
        mem.sort(i32, self.lists[1].items, {}, comptime std.sort.asc(i32));
    }

    fn partOne(self: *Puzzle) i32 {
        var sum: i32 = 0;
        for (self.lists[0].items, self.lists[1].items) |i, j| {
            sum += @intCast(@abs(i - j));
        }
        return sum;
    }

    fn partTwo(self: *Puzzle, allocator: Allocator) !i32 {
        var map = std.AutoHashMap(i32, i32).init(allocator);
        defer map.deinit();
        for (self.lists[1].items) |k| {
            try map.put(k, 1 + (map.get(k) orelse 0));
        }
        var sum: i32 = 0;
        for (self.lists[0].items) |k| {
            sum += k * (map.get(k) orelse 0);
        }
        return sum;
    }
};

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer {
        const deinit_status = gpa.deinit();
        if (deinit_status == .leak) std.testing.expect(false) catch @panic("TEST FAIL");
    }
    const allocator = gpa.allocator();

    const args = try std.process.argsAlloc(allocator);
    defer std.process.argsFree(allocator, args);

    if (args.len < 2) {
        std.debug.print("Usage: main <input_file>\n", .{});
        return;
    }
    const filename = args[1];

    var p = Puzzle.init(allocator);
    defer p.deinit();

    try p.load(filename);

    std.debug.print("{}\n", .{p.partOne()});
    const p2 = try p.partTwo(allocator);
    std.debug.print("{}\n", .{p2});
}
