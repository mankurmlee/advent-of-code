import fileinput

class Scoreboard:
    def __init__(self):
        self.a: int = 0
        self.b: int = 1
        self.data: list = [3, 7]

    def extend(self):
        tot = self.data[self.a] + self.data[self.b]
        if tot >= 10:
            self.data.append(1)
        self.data.append(tot % 10)
        n = len(self.data)
        self.a = (self.a + self.data[self.a] + 1) % n
        self.b = (self.b + self.data[self.b] + 1) % n

    def part_one(self, n: int):
        m = n + 10
        while len(self.data) < m:
            self.extend()
        return "".join(map(str, self.data[n:m]))

    def part_two(self, n):
        patt = list(map(int, str(n)))
        n = len(patt)
        while True:
            self.extend()
            if self.data[-n:] == patt:
                return len(self.data) - n
            if self.data[-n-1:-1] == patt:
                return len(self.data) - n - 1

n = int(next(fileinput.input()))

s = Scoreboard()
res = s.part_one(n)
print(f"Part 1: {res}")

s = Scoreboard()
res = s.part_two(n)
print(f"Part 2: {res}")