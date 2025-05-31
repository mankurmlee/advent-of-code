import sys

class Vec:
    def __init__(self, x: int, y: int):
        self.x = x
        self.y = y

    def __str__(self):
        return f"{self.x}, {self.y}"

    def __eq__(self, other):
        return self.x == other.x and self.y == other.y

    def __hash__(self):
        return hash((self.x, self.y))

class Rect:
    def __init__(self, id, lo: Vec, hi: Vec):
        self.id = id
        self.lo = lo
        self.hi = hi

    @classmethod
    def Parse(cls, line: str):
        line = line.replace("#", " ").replace(",", " ").replace(":", " ").replace("x", " ")
        data = line.split()
        id = data[0]
        x = int(data[2])
        y = int(data[3])
        w = int(data[4])
        h = int(data[5])
        lo = Vec(x, y)
        hi = Vec(x + w - 1, y + h - 1)
        return cls(id, lo, hi)

    def __str__(self):
        return f"<{self.id}: ({self.lo}, {self.hi})>"

class Puzzle:
    def __init__(self, filename):
        self.claims = []

        with open(filename) as f:
            data = f.read();
            lines = data.splitlines();
            for l in lines:
                l = l.strip()
                if l == "":
                    continue
                self.claims.append(Rect.Parse(l))

    def __str__(self):
        txt = []
        for p in self.claims:
            txt.append(str(p))
        return ", ".join(txt)

    def overlap_area(self):
        vscreen = {}
        for claim in self.claims:
            for y in range(claim.lo.y, claim.hi.y+1):
                for x in range(claim.lo.x, claim.hi.x+1):
                    v = Vec(x, y)
                    if v not in vscreen:
                        vscreen[v] = 0
                    vscreen[v] += 1

        area = 0
        for v in vscreen.values():
            if v > 1:
                area += 1

        return area

    def part_two(self):
        seen = set()
        n = len(self.claims)
        for i in range(n):
            if i in seen:
                continue
            seen.add(i)
            j = self.has_overlap(i)
            if j == -1:
                return self.claims[i].id
            seen.add(j)

    def has_overlap(self, i):
        n = len(self.claims)
        for j in range(n):
            if i == j:
                continue
            if self.overlaps(i, j):
                return j
        return -1

    def overlaps(self, i, j):
        p = self.claims[i]
        q = self.claims[j]
        if (p.lo.x > q.hi.x or q.lo.x > p.hi.x or
            p.lo.y > q.hi.y or q.lo.y > p.hi.y):
            return False
        return True

def main():
    p = Puzzle(sys.argv[1])
    n = p.overlap_area()
    print(f"Part 1: {n}")

    n = p.part_two()
    print(f"Part 2: {n}")

if __name__ == "__main__":
    main()