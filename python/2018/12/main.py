from dataclasses import dataclass
import sys

class Observer:
    base: int
    gradient: int = 0
    memo = []

    def __init__(self, base):
        self.base = base

    def notify(self, item):
        if self.gradient > 0:
            return
        self.memo.append(item)
        n = len(self.memo)
        if n % 2 == 0:
            return
        mid = n >> 1
        if mid == 0:
            return
        for i in range(mid, n):
            if self.memo[i] != item:
                return
        self.gradient = item
        self.memo = self.memo[:mid]

    def predict(self, n):
        extra = n - len(self.memo)
        if extra < 1:
            return self.base + sum(self.memo[:n])
        return self.base + sum(self.memo) + extra * self.gradient

@dataclass
class Cavern:
    plants: set[int]
    rules: set[str]
    first: int
    last: int
    observer: Observer

    @classmethod
    def from_file(cls, filename):
        plants = set()
        rules = set()
        first = -1
        last = -1
        with open(filename) as f:
            state = next(f).strip().split()[2]
            n = len(state)
            for i in range(len(state)):
                if state[i] == "#":
                    plants.add(i)
                    if first == -1:
                        first = i
                    last = i
            next(f)
            for l in f:
                words = l.strip().split()
                if words[2] != "#":
                    continue
                rules.add(words[0])
            return Cavern(plants, rules, first, last, Observer(sum(plants)))

    def update(self):
        plants = set()
        start = self.first - 2
        finish = self.last + 2
        first = self.last
        last = self.first
        for i in range(start,finish):
            w = self.get_window(i)
            if w not in self.rules:
                continue
            plants.add(i)
            if i < first:
                first = i
            if i > last:
                last = i

        # Notify the observer of the delta
        old_sum = sum(self.plants)
        new_sum = sum(plants)
        self.observer.notify(new_sum - old_sum)

        self.plants = plants
        self.first = first
        self.last = last

    def get_window(self, n):
        ch = []
        for i in range(n-2, n+3):
            c = '#' if i in self.plants else '.'
            ch.append(c)
        return "".join(ch)


def main():
    c = Cavern.from_file(sys.argv[1])

    for _ in range(20):
        c.update()
    s = sum(c.plants)
    print(f"Part 1: {s}")

    o = c.observer
    while o.gradient == 0:
        c.update()

    n = o.predict(50000000000)
    print(f"Part 2: {n}")

if __name__ == "__main__":
    main()
