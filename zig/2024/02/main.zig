const std = @import("std");
const mem = std.mem;
const heap = std.heap;
const Allocator = mem.Allocator;
const ArrayList = std.ArrayList;

const String = struct {
    const Self = @This();
    allocator: Allocator,
    string: []const u8,

    pub fn init(a: Allocator, s: []const u8) !Self {
        const d = try a.alloc(u8, s.len);
        @memcpy(d, s);
        return Self{ .allocator = a, .string = d };
    }

    pub fn take(a: Allocator, s: []const u8) Self {
        return Self{ .allocator = a, .string = s };
    }

    pub fn deinit(self: *Self) void {
        self.allocator.free(self.string);
    }

    pub fn clone(self: *Self) !Self {
        return Self.init(self.allocator, self.string);
    }

    pub fn format(value: Self, comptime fmt: []const u8, options: std.fmt.FormatOptions, writer: anytype) !void {
        _ = fmt;
        _ = options;
        try std.fmt.format(writer, "{s}", .{value.string});
    }
};

const IntIterator = struct {
    const Self = @This();
    buffer: []const u8,
    index: ?usize,

    fn init(b: []const u8) Self {
        return .{ .buffer = b, .index = 0 };
    }

    fn next(self: *Self) ?i64 {
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

const LineIterator = struct {
    const Self = @This();
    buffer: String,
    index: ?usize,

    fn init(s: String) Self {
        return .{ .buffer = s, .index = 0 };
    }

    fn deinit(self: *Self) void {
        self.buffer.deinit();
    }

    pub fn format(value: Self, comptime fmt: []const u8, options: std.fmt.FormatOptions, writer: anytype) !void {
        _ = fmt;
        _ = options;
        try std.fmt.format(writer, "{s}", .{value.buffer});
    }

    fn next(self: *Self) ?[]const u8 {
        var i = self.index orelse return null;
        var started = false;
        var start = i;
        var end = i;
        for (self.buffer.string[i..]) |char| {
            i += 1;
            switch (char) {
                '\n' => break,
                '\r', '\t', ' ' => {},
                else => {
                    if (!started) {
                        started = true;
                        start = i - 1;
                    }
                    end = i;
                },
            }
        }
        self.index = if (i >= self.buffer.string.len) null else i;
        return self.buffer.string[start..end];
    }
};

fn readFile(allocator: Allocator, filename: []const u8) !LineIterator {
    const file = try std.fs.cwd().openFile(filename, .{});
    defer file.close();

    const file_size = try file.getEndPos();
    const data = try allocator.alloc(u8, file_size);

    _ = try file.readAll(data);
    const s = String.take(allocator, data);
    return LineIterator.init(s);
}

fn sgn(comptime T: type, n: T) T {
    if (n > 0) {
        return 1;
    } else if (n < 0) {
        return -1;
    }
    return 0;
}

fn safe(seq: []const i32) bool {
    if (seq.len <= 1) return true;
    var x0 = seq[0];
    var check: i32 = 0;
    for (1.., seq[1..]) |i, x| {
        const diff = x - x0;
        check += sgn(i32, diff);
        if (@abs(check) != i) {
            return false;
        }
        const adiff = @abs(diff);
        if (adiff < 1 or adiff > 3) {
            return false;
        }
        x0 = x;
    }
    return true;
}

fn almostSafe(a: Allocator, seq: []const i32) !bool {
    var child = try a.alloc(i32, seq.len - 1);
    defer a.free(child);
    for (0.., seq) |i, _| {
        var k: u16 = 0;
        for (0.., seq) |j, v| {
            if (i == j) {
                continue;
            }
            child[k] = v;
            k += 1;
        }
        if (safe(child)) {
            return true;
        }
    }
    return false;
}

const Puzzle = struct {
    allocator: Allocator,
    arena: heap.ArenaAllocator,
    data: ArrayList(ArrayList(i32)),

    const Self = @This();
    fn init(a: Allocator, filename: []const u8) !Self {
        var arena = heap.ArenaAllocator.init(a);
        var lines = try readFile(a, filename);
        defer lines.deinit();

        var data = ArrayList(ArrayList(i32)).init(arena.allocator());
        while (lines.next()) |line| {
            var row = ArrayList(i32).init(arena.allocator());
            var nums = IntIterator.init(line);
            while (nums.next()) |x| {
                try row.append(@intCast(x));
            }
            try data.append(row);
        }

        return .{
            .allocator = a,
            .arena = arena,
            .data = data,
        };
    }

    fn deinit(self: *Self) void {
        _ = self.arena.reset(.free_all);
        self.arena.deinit();
    }

    fn checkReports(self: Self) ![2]u32 {
        var n_safe: u32 = 0;
        var n_almost_safe: u32 = 0;
        for (self.data.items) |row| {
            if (safe(row.items)) {
                n_safe += 1;
                continue;
            }
            const ok = try almostSafe(self.allocator, row.items);
            if (ok) {
                n_almost_safe += 1;
            }
        }
        return [2]u32{ n_safe, n_safe + n_almost_safe };
    }
};

pub fn main() !void {
    var gpa = heap.GeneralPurposeAllocator(.{}){};
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

    var p = try Puzzle.init(allocator, filename);
    defer p.deinit();

    const res = try p.checkReports();
    std.debug.print("{}\n", .{res[0]});
    std.debug.print("{}\n", .{res[1]});
}
