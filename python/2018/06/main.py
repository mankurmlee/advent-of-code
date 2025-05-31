import copy
from dataclasses import dataclass
import sys

@dataclass
class Vec:
    x: int
    y: int

    @classmethod
    def parse(cls, txt):
        words = txt.replace(",", "").split()
        return cls(int(words[0]), int(words[1]))

    def __str__(self):
        return f"{self.x}, {self.y}"

    def __hash__(self):
        return hash((self.x, self.y))

    def distance_to(self, other) -> int:
        return abs(other.x - self.x) + abs(other.y - self.y)

@dataclass
class Rect:
    x: int
    y: int
    w: int
    h: int

@dataclass
class Puzzle:
    data: list[Vec]

    @classmethod
    def from_file(cls, filename):
        data = []
        with open(filename) as f:
            for l in f.read().splitlines():
                l = l.strip()
                if l == "":
                    continue
                data.append(Vec.parse(l))
        return cls(data)

    def get_bounds(self):
        lo = copy.copy(self.data[0])
        hi = copy.copy(self.data[0])
        for v in self.data[1:]:
            if v.x < lo.x:
                lo.x = v.x
            if v.y < lo.y:
                lo.y = v.y
            if v.x > hi.x:
                hi.x = v.x
            if v.y > hi.y:
                hi.y = v.y
        return Rect(lo.x, lo.y, hi.x - lo.x + 1, hi.y - lo.y + 1)

    def safest_area(self):
        area = 0
        bounds = self.get_bounds()
        limit = 32 if bounds.w < 100 else 10000
        for y in range(bounds.y, bounds.y + bounds.h):
            for x in range(bounds.x, bounds.x + bounds.w):
                sum = self.safe_sum(Vec(x, y))
                if sum < limit:
                    area += 1
        return area

    def safe_sum(self, pos: Vec):
        sum = 0
        for v in self.data:
            sum += pos.distance_to(v)
        return sum

class Canvas:
    def __init__(self, bounds):
        size = bounds.w * bounds.h
        self.bounds: Rect = bounds
        self.dist: list[int] = [size] * size
        self.owner: list[Vec] = [None] * size

    def flood_fill(self, point: Vec):
        px = point.x - self.bounds.x
        py = point.y - self.bounds.y
        for y in range(self.bounds.h):
            for x in range(self.bounds.w):
                i = y * self.bounds.w + x
                dist = abs(x - px) + abs(y - py)
                if dist < self.dist[i]:
                    self.dist[i] = dist
                    self.owner[i] = point
                elif dist == self.dist[i]:
                    self.owner[i] = Vec(-1, -1)

    def print(self):
        w = self.bounds.w
        size = w * self.bounds.h
        for i in range(0, size, w):
            print([str(v) for v in self.owner[i:i+w]])
        for i in range(0, size, w):
            print(self.dist[i:i+w])

    def infinite_set(self):
        w, h = self.bounds.w, self.bounds.h
        ignore: set[Vec] = set()
        off = (h - 1) * w
        for x in range(w):
            ignore.add(self.owner[x])
            ignore.add(self.owner[x + off])
        off = w - 1
        for y in range(h):
            i = y * w
            ignore.add(self.owner[i])
            ignore.add(self.owner[i + off])
        return ignore

    def biggest_area(self):
        ignore = self.infinite_set()
        areas = {}
        for v in self.owner:
            if v in ignore:
                continue
            if v not in areas:
                areas[v] = 0
            areas[v] += 1
        biggest = 0
        for v in areas.values():
            if v > biggest:
                biggest = v
        return biggest

def main():
    p = Puzzle.from_file(sys.argv[1])
    c = Canvas(p.get_bounds())
    for v in p.data:
        c.flood_fill(v)

    n = c.biggest_area()
    print(f"Part 1: {n}")

    n = p.safest_area()
    print(f"Part 2: {n}")

if __name__ == "__main__":
    main()
