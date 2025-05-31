from dataclasses import dataclass
import sys

@dataclass
class Vec:
    x: int
    y: int

    def __iadd__(self, other):
        return Vec(self.x + other.x, self.y + other.y)

    def __str__(self):
        return f"{self.x},{self.y}"

    def __hash__(self):
        return hash((self.x, self.y))

    def tuple(self):
        return (self.x, self.y)

@dataclass
class Cart:
    id: tuple[int, int]
    pos: Vec
    vel: Vec
    turn: int = 0

    dirs = [Vec(0, -1), Vec(1, 0), Vec(0, 1), Vec(-1, 0)]

    def steer(self, tile):
        if tile == "|" or tile == "-":
            return

        turn_left = (
            tile == "\\" and self.vel.x == 0 or
            tile == "/" and self.vel.x != 0 or
            tile == "+" and self.turn == 0
        )

        if tile == '+':
            self.turn = (self.turn + 1) % 3
            if self.turn == 2:
                return

        d = Cart.dirs.index(self.vel)
        if turn_left:
            d -= 1
            if d < 0:
                d += 4
        else:
            d = (d + 1) % 4

        self.vel = Cart.dirs[d]

class MiningRail:

    def __init__(self, filename, collisions):
        lines = []
        with open(filename) as f:
            for l in f:
                l = l.rstrip()
                lines.append(l)

        w = 0
        h = len(lines)
        for l in lines:
            if len(l) > w:
                w = len(l)

        grid = [' '] * w * h

        cart_dir = dict({
            '^': Vec(0, -1), 'v': Vec(0, 1), '<': Vec(-1, 0), '>': Vec(1, 0),
        })
        carts: list[Cart] = []
        i = 0
        for y in range(h):
            l = lines[y]
            n = len(l)
            grid[i:i+n] = l
            for x in range(n):
                if l[x] in cart_dir.keys():
                    d = cart_dir[l[x]]
                    v = Vec(x, y)
                    c = Cart(v.tuple(), v, d)
                    carts.append(c)
                    grid[i+x] = '|' if d.x == 0 else '-'
            i += w

        self.width = w
        self.height = h
        self.grid = grid
        self.carts = carts
        self.collisions = collisions

    def update(self):
        dead = set()
        w = self.width
        n = len(self.carts)
        alive = n
        grid = self.grid
        self.carts.sort(key=lambda c: (c.pos.y, c.pos.x))
        for c in self.carts:
            c.pos += c.vel
            c.steer(grid[w * c.pos.y + c.pos.x])
            if len(set(o.pos for o in self.carts if o.id not in dead)) == alive:
                continue
            for other in self.carts:
                if other.pos == c.pos:
                    dead.add(other.id)
            alive = n - len(dead)
            self.collisions.append(c.pos)

        self.carts = [c for c in self.carts if c.id not in dead]

def main():
    collisions = []
    m = MiningRail(sys.argv[1], collisions)

    while len(m.carts) > 1:
        m.update()

    print(f"Part 1: {collisions[0]}")
    if len(m.carts) == 1:
        pos = m.carts[0].pos
        print(f"Part 2: {pos}")

if __name__ == "__main__":
    main()
