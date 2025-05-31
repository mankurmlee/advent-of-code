import sys

class Puzzle:
    def __init__(self, filename: str):
        self.ids = []
        with open(filename) as f:
            data = f.read();
            lines = data.splitlines();
            for l in lines:
                l = l.strip()
                if l == "":
                    continue
                self.ids.append(l)

    def __str__(self):
        txt = "\n".join(v for v in self.ids)
        return f"{txt}"

    def part_one(self) -> int:
        twos = 0
        threes = 0

        for id in self.ids:
            two = False
            three = False

            for v in self.count_letters(id).values():
                if v == 2:
                    two = True
                elif v == 3:
                    three = True

            if two:
                twos += 1
            if three:
                threes += 1

        return twos * threes

    @staticmethod
    def count_letters(id: str) -> dict[str, int]:
        cnts = {}
        for c in id:
            if c not in cnts:
                cnts[c] = 0
            cnts[c] += 1
        return cnts

    def part_two(self) -> str:
        n = len(self.ids)
        for i in range(n):
            for j in range(i+1, n):
                if self.almost_matches(self.ids[i], self.ids[j]):
                    return self.common_letters(self.ids[i], self.ids[j])
        return ""

    @staticmethod
    def almost_matches(a: str, b: str) -> bool:
        diff = 0
        for i in range(len(a)):
            if a[i] == b[i]:
                continue
            diff += 1
            if diff > 1:
                return False
        return True

    @staticmethod
    def common_letters(a: str, b: str) -> str:
        common = ""
        for p, q in zip(a, b):
            if p == q:
                common += p

        return common

def main():
    p = Puzzle(sys.argv[1])
    print(f"Part 1: {p.part_one()}")
    print(f"Part 2: {p.part_two()}")

if __name__ == "__main__":
    main()