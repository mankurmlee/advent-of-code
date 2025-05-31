const std = @import("std");

pub fn main() !void {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    var args = try std.process.ArgIterator.initWithAllocator(allocator);
    _ = args.next();
    const filename = args.next().?;

    const ws = try Puzzle.load(allocator, filename);
    ws.partOne();
    ws.partTwo();
}

fn readFile(allocator: std.mem.Allocator, filename: []const u8) ![]u8 {
    const file = try std.fs.cwd().openFile(filename, .{});
    defer file.close();

    const file_size = try file.getEndPos();
    const data = try allocator.alloc(u8, file_size);

    _ = try file.readAll(data);
    return data;
}

const Puzzle = struct {
    width: u8 = undefined,
    height: u8 = undefined,
    cells: []u8 = undefined,

    pub fn load(allocator: std.mem.Allocator, filename: []const u8) !Puzzle {
        const buffer = try readFile(allocator, filename);
        var width: u8 = 0;
        var height: u8 = 0;
        var cells = std.ArrayList(u8).init(allocator);
        defer cells.deinit();

        var lines = std.mem.splitScalar(u8, buffer, '\n');
        if (lines.peek()) |line| {
            const trimmed = std.mem.trim(u8, line, "\r ");
            width = @intCast(trimmed.len);
        }

        while (lines.next()) |line| {
            const trimmed = std.mem.trim(u8, line, "\r ");
            if (trimmed.len == 0) {
                continue;
            }
            try cells.appendSlice(trimmed);
            height += 1;
        }

        return Puzzle{
            .width = width,
            .height = height,
            .cells = try cells.toOwnedSlice(),
        };
    }

    pub fn partOne(self: *const Puzzle) void {
        const dirs = [_][2]i16{
            .{ 0, -1 }, .{ 1, -1 }, .{ 1, 0 },  .{ 1, 1 },
            .{ 0, 1 },  .{ -1, 1 }, .{ -1, 0 }, .{ -1, -1 },
        };
        const w = self.width;
        const h = self.height;
        const n = self.cells.len;

        var sum: u16 = 0;

        for (dirs) |d| {
            const dx, const dy = d;
            var i: u32 = 0;
            while (i < n) : (i += 1) {
                var x: i16 = @intCast(i % w);
                var y: i16 = @intCast(i / w);
                for ("XMAS") |expected| {
                    if (x < 0 or y < 0 or x >= w or y >= h or self.cells[@intCast(y * w + x)] != expected) {
                        break;
                    }
                    x += dx;
                    y += dy;
                } else {
                    sum += 1;
                }
            }
        }
        std.debug.print("Part 1: {d}\n", .{sum});
    }

    pub fn partTwo(self: *const Puzzle) void {
        var sum: u16 = 0;
        var y: u8 = 1;
        while (y < self.height - 1) : (y += 1) {
            var x: u8 = 1;
            while (x < self.width - 1) : (x += 1) {
                if (self.xmas(x, y)) {
                    sum += 1;
                }
            }
        }
        std.debug.print("Part 2: {d}\n", .{sum});
    }

    fn xmas(self: *const Puzzle, x: u8, y: u8) bool {
        const w: i16 = self.width;
        if (self.cells[@intCast(y * w + x)] != 'A') {
            return false;
        }
        const diags = [_][2]i16{
            .{ -1, -1 }, .{ 1, 1 },
            .{ -1, 1 },  .{ 1, -1 },
        };
        var letters: [4]u8 = undefined;
        for (diags, 0..) |d, i| {
            const x1: i16 = x + d[0];
            const y1: i16 = y + d[1];
            const o: i16 = @intCast(y1 * w + x1);
            const l = self.cells[@intCast(o)];
            if (l != 'M' and l != 'S') return false;
            letters[i] = l;
        }
        if ((letters[0] != 'M' or letters[1] != 'S') and (letters[0] != 'S' or letters[1] != 'M')) {
            return false;
        }
        if ((letters[2] != 'M' or letters[3] != 'S') and (letters[2] != 'S' or letters[3] != 'M')) {
            return false;
        }
        return true;
    }
};
