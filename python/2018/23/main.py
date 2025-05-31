from dataclasses import dataclass
import fileinput
import re

@dataclass
class Bot:
    pos: tuple[int, int, int]
    r: int

    def reaches(self, other):
        return self.r >= sum(abs(b - a) for a, b in zip(self.pos, other.pos))

    def range_overlaps(self, other):
        return self.r + other.r >= sum(abs(b - a) for a, b in zip(self.pos, other.pos))

    def dist(self):
        return sum(self.pos) - self.r

class Solver:
    def __init__(self, bots):
        n = len(bots)
        links = []
        for i, b in enumerate(bots):
            r = set(j for j in range(i+1, n) if b.range_overlaps(bots[j]))
            links.append(r)
        self.links = links
        self.cache = {frozenset({i}): frozenset({i}) for i in range(n)}

    def largest_group(self, bots):
        if bots in self.cache:
            return self.cache[bots]
        largest = frozenset()
        for b in sorted(bots):
            l = self.links[b]
            if len(l) + 1 > len(largest):
                child = self.largest_group(bots & l) | {b}
                if len(child) > len(largest):
                    largest = child
        self.cache[bots] = largest
        return largest

bot_params = [list(map(int, re.findall(r'-?\d+', l))) for l in fileinput.input()]
bots = [Bot((p[0], p[1], p[2]), p[3]) for p in bot_params]

strongest = max(bots, key=lambda b: b.r)
print(f"Part 1: {sum(1 for b in bots if strongest.reaches(b))}")

s = Solver(bots)
g = s.largest_group(frozenset(range(len(bots))))

print(f"Part 2: {max(bots[i].dist() for i in g)}")