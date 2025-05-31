from dataclasses import dataclass
import sys

@dataclass
class Vec:
    x: int
    y: int

    def copy(self):
        return Vec(self.x, self.y)

@dataclass
class Proj:
    pos: Vec
    vel: Vec

    def update(self):
        self.pos.x += self.vel.x
        self.pos.y += self.vel.y

@dataclass
class Rect:
    lo: Vec
    hi: Vec

@dataclass
class SkyField:
    stars: list[Proj]

    @classmethod
    def from_file(cls, filename):
        data = []
        with open(filename) as f:
            for l in f:
                data.append(cls.parse(l))
        return cls(data)

    @staticmethod
    def parse(line):
        l = line.replace("<", " ").replace(">", " ").replace(",", " ")
        words = l.split()
        pos = Vec(int(words[1]), int(words[2]))
        vel = Vec(int(words[4]), int(words[5]))
        return Proj(pos, vel)

    def update(self):
        for i in range(len(self.stars)):
            self.stars[i].update()

    def height(self):
        lo = self.stars[0].pos.y
        hi = lo
        for p in self.stars:
            if p.pos.y < lo:
                lo = p.pos.y
            if p.pos.y > hi:
                hi = p.pos.y
        return hi - lo + 1

    def bounds(self):
        lo = self.stars[0].pos.copy()
        hi = lo.copy()
        for p in self.stars:
            if p.pos.x < lo.x:
                lo.x = p.pos.x
            if p.pos.x > hi.x:
                hi.x = p.pos.x
            if p.pos.y < lo.y:
                lo.y = p.pos.y
            if p.pos.y > hi.y:
                hi.y = p.pos.y

        return Rect(lo, hi)

    def print(self):
        r = self.bounds()
        ox = r.lo.x
        oy = r.lo.y
        stars = set((p.pos.x - ox, p.pos.y - oy) for p in self.stars)
        w = r.hi.x - ox + 1
        h = r.hi.y - oy + 1
        for y in range(h):
            line = []
            for x in range(w):
                if (x, y) in stars:
                    line.append("*")
                else:
                    line.append(" ")
            print("".join(line))

def main():
    f = SkyField.from_file(sys.argv[1])
    i = 0
    while(f.height() > 10):
        f.update()
        i += 1
    print("Part 1:")
    f.print()
    print(f"Part 2: {i}")

if __name__ == "__main__":
    main()
