import re
import sys

class Puzzle:
    def __init__(self, filename):
        self.deltas = []
        patt = r'[,\n]'
        with open(filename) as f:
            data = f.read();
            deltas = re.split(patt, data)
            for d in deltas:
                d = d.strip()
                if d == "":
                    continue
                self.deltas.append(int(d))

    def __str__(self):
        txt = ", ".join(str(v) for v in self.deltas)
        return f"{txt}"

    def part_one(self):
        r = 0
        for d in self.deltas:
            r += d
        return r

    def part_two(self):
        seen = set()
        r = 0
        seen.add(0)
        while 1:
            for d in self.deltas:
                r += d
                if r not in seen:
                    seen.add(r)
                else:
                    return r




def main():
    p = Puzzle(sys.argv[1])
    print(f"Part 1: {p.part_one()}")
    print(f"Part 2: {p.part_two()}")

if __name__ == "__main__":
    main()