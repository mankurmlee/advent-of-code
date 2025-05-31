from dataclasses import dataclass
import sys

@dataclass
class Node:
    last: int
    next: int

class Circle:
    def __init__(self, size=1):
        self.current = 0
        self.marble = 1
        self.data = [Node(0, 0)] * (size + 1)

    def __str__(self):
        return str(self.data)

    def take_turn(self):
        if self.marble % 23 == 0:
            return self.remove()
        self.insert()
        return 0

    def insert(self):
        last = self.data[self.current].next
        next = self.data[last].next
        self.data[last].next = self.marble
        self.data[next].last = self.marble
        self.data[self.marble] = Node(last, next)
        self.current = self.marble
        self.marble += 1

    def remove(self):
        s = self.marble
        c = self.current
        for _ in range(7):
            c = self.data[c].last
        n = self.data[c]
        self.data[n.last].next = n.next
        self.data[n.next].last = n.last
        self.current = n.next
        self.marble += 1
        # del self.data[c]
        return s + c

class Game:
    def __init__(self, filename):
        p = 0
        m = 0
        with open(filename) as f:
            words = f.readline().strip().split()
            p = int(words[0])
            m = int(words[6])

        self.players = p
        self.marbles = m

    def play(self, scale=1):
        scores = [0] * self.players
        size = self.marbles * scale
        circle = Circle(size)
        for i in range(size):
            scores[i % self.players] += circle.take_turn()
        return scores

def main():
    g = Game(sys.argv[1])

    s = max(g.play())
    print(f"Part 1: {s}")

    s = max(g.play(100))
    print(f"Part 2: {s}")

if __name__ == "__main__":
    main()
