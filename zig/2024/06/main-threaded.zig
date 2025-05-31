pub fn main() !void {
    // var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    // defer _ = gpa.deinit();
    // const alloc = gpa.allocator();

    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const alloc = arena.allocator();

    var thread_arena: std.heap.ThreadSafeAllocator = .{
        .child_allocator = alloc,
    };

    const core_count = try std.Thread.getCpuCount();
    print("Number of cores: {d}\n", .{core_count});
    var pool: std.Thread.Pool = undefined;
    try pool.init(std.Thread.Pool.Options{
        .allocator = thread_arena.allocator(),
        .n_jobs = core_count,
    });
    defer pool.deinit();

    var args = try std.process.ArgIterator.initWithAllocator(alloc);
    defer args.deinit();

    _ = args.next();
    const filename = args.next().?;

    var p = try Puzzle.init(thread_arena.allocator(), filename);
    defer p.deinit();

    try p.solve(&pool);
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
    puzzle: *const Puzzle,
    obj: Vec = .{ -1, -1 },

    const Self = @This();
    const dirs = [4]Vec{ .{ 0, -1 }, .{ 1, 0 }, .{ 0, 1 }, .{ -1, 0 } };

    pub fn init(p: *const Puzzle) Self {
        return Self{
            .guard = .{ .pos = p.start },
            .puzzle = p,
        };
    }

    pub fn next(self: *Self) Vec {
        const p = self.guard.pos;
        var out: Vec = p + dirs[self.guard.dir];
        while (self.puzzle.walls.contains(out) or out[0] == self.obj[0] and out[1] == self.obj[1]) {
            self.guard.dir +%= 1;
            out = p + dirs[self.guard.dir];
        }
        self.guard.pos = out;
        return out;
    }

    pub fn inBounds(self: Self) bool {
        const x, const y = self.guard.pos;
        const w, const h = self.puzzle.size;
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

    pub fn solve(self: *Self, pool: *std.Thread.Pool) !void {
        var visited = AutoHashMap(Vec, void).init(self.alloc);
        defer visited.deinit();

        var w = Walker.init(self);

        var pos = w.guard.pos;
        while (w.inBounds()) {
            try visited.put(pos, {});
            pos = w.next();
        }
        print("Part 1: {d}\n", .{visited.count()});

        const num = visited.count();
        const res = try self.alloc.alloc(bool, num);
        defer self.alloc.free(res);
        var wg: std.Thread.WaitGroup = .{};
        wg.reset();

        const start: [2]i16 = self.start;

        var i: usize = 0;
        var spots = visited.keyIterator();
        while (spots.next()) |p| : (i += 1) {
            const s = p.*;
            if (s[0] == start[0] and s[1] == start[1]) {
                continue;
            }
            pool.spawnWg(&wg, Self.goWillLoop, .{ self, s, &res[i] });
        }

        wg.wait();

        var loop_count: i16 = 0;
        for (res) |r| {
            if (r) loop_count += 1;
        }

        print("Part 2: {d}\n", .{loop_count});
    }

    pub fn goWillLoop(self: *Self, obj: Vec, res: *bool) void {
        res.* = self.willLoop(obj) catch false;
    }

    pub fn willLoop(self: *Self, obj: Vec) !bool {
        var w = Walker.init(self);
        w.obj = obj;

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
