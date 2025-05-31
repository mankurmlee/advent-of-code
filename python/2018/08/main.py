from dataclasses import dataclass
import sys

@dataclass
class Node:
    num_children: int
    num_metadata: int
    children: list
    metadata: list[int]

    @classmethod
    def from_file(cls, filename):
        words = []
        with open(filename) as f:
            words = f.readline().strip().split()
        return cls.parse((int(w) for w in words))

    @classmethod
    def parse(cls, data):
        nc = next(data)
        nm = next(data)
        c = [cls.parse(data) for _ in range(nc)]
        m = [next(data) for _ in range(nm)]
        return cls(nc, nm, c, m)

    def sum_metadata(self):
        tot = sum(n.sum_metadata() for n in self.children)
        return tot + sum(self.metadata)

    def value(self) -> int:
        if self.num_children == 0:
            return sum(self.metadata)

        cache = {}
        v = 0
        for i in self.metadata:
            if i < 1 or i > self.num_children:
                continue
            if i not in cache:
                cache[i] = self.children[i - 1].value()
            v += cache[i]
        return v


def main():
    root = Node.from_file(sys.argv[1])

    n = root.sum_metadata()
    print(f"Part 1: {n}")

    n = root.value()
    print(f"Part 2: {n}")

if __name__ == "__main__":
    main()
