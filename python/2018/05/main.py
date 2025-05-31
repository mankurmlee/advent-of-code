import copy
import string
import sys

class Polymer:
    def __init__(self, data):
        self.data = data

    def __str__(self):
        return "".join(self.data)

    def react(self):
        n = len(self.data)
        last = self.data[n-1]
        for i in range(n - 2, 0, -1):
            if self.data[i] == last:
                continue
            if self.data[i].lower() != last.lower():
                last = self.data[i]
                continue
            del self.data[i+1]
            del self.data[i]
            n -= 2
            if i < n:
                last = self.data[i]
            else:
                last = ""

    def fully_react(self):
        self.react()
        return len(self.data)

    def remove_unit(self, lower: str):
        upper = lower.upper()
        n = len(self.data)
        for i in range(n - 1, 0, -1):
            if self.data[i] == lower or self.data[i] == upper:
                del self.data[i]

    @classmethod
    def FromFile(cls, filename):
        data = []
        with open(filename) as f:
            data = list(f.readline().strip())
        return cls(data)

    def copy(self):
        return Polymer(self.data[:])

class Solver:
    def __init__(self, polymer):
        self.polymer = copy.deepcopy(polymer)

    def part_two(self):
        best = len(self.polymer.data)
        for letter in string.ascii_lowercase:
            p = self.polymer.copy()
            p.remove_unit(letter)
            n = p.fully_react()
            if n < best:
                best = n

        return best

def main():
    p = Polymer.FromFile(sys.argv[1])
    s = Solver(p)

    n = p.fully_react()
    print(f"Part 1: {n}")

    n = s.part_two()
    print(f"Part 2: {n}")

if __name__ == "__main__":
    main()
