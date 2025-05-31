from dataclasses import dataclass
import sys
import time

class Grid:
    def __init__(self, filename):
        s = 0
        with open(filename) as f:
            s = int(next(f).strip())

        self.serial = s
        cells = [0] * 300 * 300
        for y in range(300):
            for x in range(300):
                cells[y*300+x] = self._power_level(x+1, y+1)

        self.cells = cells

    def _power_level(self, x: int, y: int):
        rack_id = x + 10
        p = rack_id * y + self.serial
        return ((p * rack_id) // 100) % 10 - 5

    def best_for_size(self, size: int):
        cells = self.cells
        bx = 0
        by = 0
        bp = 0
        n = 301 - size
        memo = [0] * 300
        for x in range(n):
            # create sq
            p = 0
            s = x
            for y in range(size):
                memo[y] = sum(cells[s:s+size])
                p += memo[y]
                s += 300

            # test sq
            if p > bp:
                bp = p
                bx, by = x, 0

            for y in range(size, 300):
                # move sq
                memo[y] = sum(cells[s:s+size])
                p += memo[y] - memo[y-size]
                s += 300

                # test sq
                if p > bp:
                    bp = p
                    bx, by = x, y+1-size

        return bx+1, by+1, bp

    def best(self):
        most = 0
        coord = 0, 0, 0
        for i in range(1, 300):
            (x, y, pow) = self.best_for_size(i)
            if pow > most:
                most = pow
                coord = x, y, i
        return coord



def main():
    g = Grid(sys.argv[1])

    (x, y, _) = g.best_for_size(3)
    print(f"Part 1: {x},{y}")

    (x, y, s) = g.best()
    print(f"Part 2: {x},{y},{s}")

if __name__ == "__main__":
    main()
