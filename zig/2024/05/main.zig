pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();

    const alloc = gpa.allocator();

    var args = try std.process.ArgIterator.initWithAllocator(alloc);
    defer args.deinit();

    _ = args.next();
    const filename = args.next().?;

    var puzzle = try Puzzle.parse(alloc, filename);
    defer puzzle.deinit();

    try puzzle.run();
}

fn readFile(alloc: Allocator, filename: []const u8) ![]u8 {
    const bufsize = 4096;
    var buf = [_]u8{0} ** bufsize;

    var data = ArrayList(u8).init(alloc);
    defer data.deinit();

    const file = try std.fs.cwd().openFile(filename, .{});
    defer file.close();

    var bytesRead: usize = bufsize;
    while (bytesRead == bufsize) {
        bytesRead = try file.read(buf[0..]);
        try data.appendSlice(buf[0..bytesRead]);
    }
    return data.toOwnedSlice();
}

const Order = struct {
    lo: u8,
    hi: u8,
};

const MySlice = struct {
    start: u16,
    end: u16,
};

const Puzzle = struct {
    const Self = @This();
    alloc: Allocator,
    rules: std.AutoHashMap(Order, bool),
    manpages: []u8,
    mans: []MySlice,

    pub fn parse(alloc: Allocator, filename: []const u8) !Puzzle {
        const data = try readFile(alloc, filename);
        defer alloc.free(data);

        var rules = std.AutoHashMap(Order, bool).init(alloc);

        var manpages = ArrayList(u8).init(alloc);
        defer manpages.deinit();

        var mans = ArrayList(MySlice).init(alloc);
        defer mans.deinit();

        var section: u8 = 0;

        var lines = std.mem.splitScalar(u8, data, '\n');
        while (lines.next()) |line| {
            const trimmed = std.mem.trim(u8, line, "\r ");
            if (trimmed.len == 0) {
                section += 1;
                continue;
            }
            switch (section) {
                0 => {
                    var ordering = std.mem.tokenizeScalar(u8, trimmed, '|');
                    const lo = try parseInt(u8, ordering.next().?, 10);
                    const hi = try parseInt(u8, ordering.next().?, 10);
                    try rules.put(.{ .lo = lo, .hi = hi }, true);
                },
                1 => {
                    const offset = manpages.items.len;
                    var pages = std.mem.tokenizeScalar(u8, trimmed, ',');
                    while (pages.next()) |s| {
                        const p = try parseInt(u8, s, 10);
                        try manpages.append(p);
                    }
                    try mans.append(.{
                        .start = @intCast(offset),
                        .end = @intCast(manpages.items.len),
                    });
                },
                else => {
                    break;
                },
            }
        }

        return Puzzle{
            .alloc = alloc,
            .rules = rules,
            .manpages = try manpages.toOwnedSlice(),
            .mans = try mans.toOwnedSlice(),
        };
    }

    pub fn deinit(self: *Puzzle) void {
        self.rules.deinit();
        self.alloc.free(self.manpages);
        self.alloc.free(self.mans);
    }

    pub fn lessThan(self: *Self, a: u8, b: u8) bool {
        return self.rules.contains(.{
            .lo = a,
            .hi = b,
        });
    }

    pub fn run(self: *Self) !void {
        var same: u32 = 0;
        var diff: u32 = 0;

        var cloned = try self.alloc.dupe(u8, self.manpages);
        defer self.alloc.free(cloned);

        man_loop: for (self.mans) |m| {
            std.mem.sort(u8, cloned[m.start..m.end], self, Puzzle.lessThan);
            const iMid = (m.start + m.end - 1) >> 1;
            const mid = self.manpages[iMid];
            const mid1 = cloned[iMid];
            if (mid != mid1) {
                diff += mid1;
                continue;
            }
            for (m.start..m.end) |i| {
                const a = self.manpages[i];
                const b = cloned[i];
                if (a != b) {
                    diff += mid1;
                    continue :man_loop;
                }
            }
            same += mid1;
        }

        print("Part 1: {d}\n", .{same});
        print("Part 2: {d}\n", .{diff});
    }
};

const std = @import("std");
const Allocator = std.mem.Allocator;
const ArrayList = std.ArrayList;
const parseInt = std.fmt.parseInt;
const print = std.debug.print;
