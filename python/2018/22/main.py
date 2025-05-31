import fileinput
from heapq import heappop, heappush
import re

class Scanner:
    dirs = [(0, -1), (1, 0), (0, 1), (-1, 0)]
    def __init__(self):
        it = fileinput.input()
        self.depth = list(map(int, re.findall(r'\d+', next(it))))[0]
        self.target = tuple(map(int, re.findall(r'\d+', next(it))))
        self.glcache = {
            (0, 0): 0,
            self.target: 0,
        }
        self.rcache = {}

    def geologic_index(self, v):
        if v in self.glcache:
            return self.glcache[v]
        x, y = v
        if y == 0:
            idx = x * 16807
        elif x == 0:
            idx = y * 48271
        else:
            idx = self.erosion_level((x-1, y)) * self.erosion_level((x, y-1))
        self.glcache[v] = idx
        return idx

    def erosion_level(self, v):
        return (self.geologic_index(v) + self.depth) % 20183

    def region_type(self, k):
        if k in self.rcache:
            return self.rcache[k]
        v = self.erosion_level(k) % 3
        self.rcache[k] = v
        return v

    def risk_level(self):
        w, h = self.target
        return sum(self.region_type((x, y)) for x in range(w+1) for y in range(h+1))

    def best_path(self):
        target = self.target[0], self.target[1], 1
        been = {(0, 0, 1): 0}
        pq = [(0, 0, 0, 1)]
        while True:
            c, x, y, t = heappop(pq)
            if (x, y, t) == target:
                return c
            for dx, dy in self.dirs:
                nc = c + 1
                nx = x + dx
                ny = y + dy
                if nx < 0 or ny < 0:
                    continue
                if (nx, ny, t) in been and nc >= been[(nx, ny, t)]:
                    continue
                if self.region_type((nx, ny)) == t:
                    continue
                been[(nx, ny, t)] = nc
                heappush(pq, (nc, nx, ny, t)) # type: ignore
            r = self.region_type((x, y))
            for dt in range(1, 3):
                nt = (t + dt) % 3
                if nt == r:
                    continue
                nc = c + 7
                if (x, y, nt) in been and nc >= been[(x, y, nt)]: # type: ignore
                    continue
                been[(x, y, nt)] = nc # type: ignore
                heappush(pq, (nc, x, y, nt)) # type: ignore

s = Scanner()
print(f"Part 1: {s.risk_level()}")
print(f"Part 2: {s.best_path()}")