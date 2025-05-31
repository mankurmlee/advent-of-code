pub fn main() !void {
    // var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    // defer _ = gpa.deinit();
    // const alloc = gpa.allocator();

    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const alloc = arena.allocator();

    var args = try std.process.ArgIterator.initWithAllocator(alloc);
    defer args.deinit();

    _ = args.next();
    const filename = args.next().?;

    var p = try Puzzle.init(alloc, filename);
    defer p.deinit();

    try p.solve();
}

fn readFile(alloc: Allocator, filename: []const u8) ![]u8 {
    const file = try std.fs.cwd().openFile(filename, .{});
    defer file.close();

    const file_size = try file.getEndPos();
    const data = try alloc.alloc(u8, file_size);

    _ = try file.readAll(data);
    return data;
}

const Vec = @Vector(2, i16);
const Guard = struct { pos: Vec, dir: u2 = 0 };
const Walker = struct {
    guard: Guard,
    walls: AutoHashMap(Vec, void),
    size: Vec,

    const Self = @This();
    const dirs = [4]Vec{ .{ 0, -1 }, .{ 1, 0 }, .{ 0, 1 }, .{ -1, 0 } };

    pub fn init(p: *const Puzzle) !Self {
        return Self{
            .guard = .{ .pos = p.start },
            .walls = try p.walls.clone(),
            .size = p.size,
        };
    }

    pub fn deinit(self: *Self) void {
        self.walls.deinit();
    }

    pub fn next(self: *Self) void {
        const p = self.guard.pos;
        var out: Vec = p + dirs[self.guard.dir];
        while (self.walls.contains(out)) {
            self.guard.dir +%= 1;
            out = p + dirs[self.guard.dir];
        }
        self.guard.pos = out;
    }

    pub fn inBounds(self: Self) bool {
        const x, const y = self.guard.pos;
        const w, const h = self.size;
        return (x >= 0 and y >= 0 and x < w and y < h);
    }
};

const Puzzle = struct {
    alloc: Allocator,
    walls: AutoHashMap(Vec, void),
    start: Vec,
    size: Vec,

    const Self = @This();

    pub fn init(alloc: Allocator, filename: []const u8) !Self {
        var walls = AutoHashMap(Vec, void).init(alloc);
        var s: ?Vec = null;
        var x: i16 = 0;
        var y: i16 = 0;

        const data = try readFile(alloc, filename);
        defer alloc.free(data);

        var lines = std.mem.tokenizeAny(u8, data, " \r\n");
        while (lines.next()) |l| {
            x = 0;
            for (l) |ch| {
                const v = .{ x, y };
                switch (ch) {
                    '#' => try walls.put(v, {}),
                    '^' => s = v,
                    else => {},
                }
                x += 1;
            }
            y += 1;
        }

        return .{
            .alloc = alloc,
            .walls = walls,
            .start = s.?,
            .size = .{ x, y },
        };
    }

    pub fn deinit(self: *Self) void {
        self.walls.deinit();
    }

    pub fn solve(self: *Self) !void {
        var visited = AutoHashMap(Vec, void).init(self.alloc);
        defer visited.deinit();

        var w = try Walker.init(self);
        defer w.deinit();

        while (w.inBounds()) {
            try visited.put(w.guard.pos, {});
            w.next();
        }
        print("Part 1: {d}\n", .{visited.count()});

        const start: [2]i16 = self.start;

        var loop_count: i16 = 0;
        var spots = visited.keyIterator();
        while (spots.next()) |p| {
            const s = p.*;
            if ((s[0] != start[0] or s[1] != start[1]) and try self.willLoop(s)) {
                loop_count += 1;
            }
        }

        print("Part 2: {d}\n", .{loop_count});
    }

    pub fn willLoop(self: *Self, obj: Vec) !bool {
        var w = try Walker.init(self);
        defer w.deinit();

        try w.walls.put(obj, {});

        var past = AutoHashMap(Guard, void).init(self.alloc);
        defer past.deinit();

        while (w.inBounds()) {
            const res = try past.getOrPut(w.guard);
            if (res.found_existing) return true;
            _ = w.next();
        }
        return false;
    }
};

const std = @import("std");
const Allocator = std.mem.Allocator;
const ArrayList = std.ArrayList;
const AutoHashMap = std.AutoHashMap;
const print = std.debug.print;
