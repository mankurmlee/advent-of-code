import fileinput

class Grid:
    def __init__(self, memo):
        self.data = []
        self.width = 0
        for line in fileinput.input():
            self.width = len(line.strip())
            self.data.extend(line.strip())
        self.height = len(self.data) // self.width
        self.buffer = [' '] * len(self.data)
        self.memo = memo
        memo.record(self.resources())

    def step(self):
        self.buffer, self.data = self.data, self.buffer
        w = self.width
        for y in range(self.height):
            for x in range(self.width):
                i = w * y + x
                ch = self.buffer[i]
                t, l = self.count_adj(x, y)
                if ch == '.' and t >= 3:
                    self.data[i] = '|'
                elif ch == '|' and l >= 3:
                    self.data[i] = '#'
                elif ch == '#' and (t == 0 or l == 0):
                    self.data[i] = '.'
                else:
                    self.data[i] = ch
        self.memo.record(self.resources())

    def count_adj(self, x_pos, y_pos):
        w = self.width
        h = self.height
        trees = 0
        lumberyard = 0
        for dx, dy in [
            (-1, -1), (0, -1), (1, -1), (-1, 0),
            (-1,  1), (0,  1), (1,  1), ( 1, 0),
        ]:
            x, y = x_pos + dx, y_pos + dy
            if x < 0 or y < 0 or x >= w or y >= h:
                continue
            ch = self.buffer[w * y + x]
            if ch == "|":
                trees += 1
            elif ch == "#":
                lumberyard += 1
        return trees, lumberyard

    def draw(self):
        w = self.width
        i = 0
        for _ in range(self.height):
            print("".join(self.data[i:i+w]))
            i += w

    def resources(self):
        t = len(list(ch for ch in self.data if ch == '|'))
        l = len(list(ch for ch in self.data if ch == '#'))
        return t * l

class Memo:
    def __init__(self):
        self.data: list[int] = []
        self.index: dict[int, list[int]] = {}
        self.base = 0
        self.period = 0

    def record(self, x):
        if self.period > 0:
            return
        if x not in self.index:
            self.index[x] = []
        indicies = self.index[x]
        indicies.append(len(self.data))
        self.data.append(x)
        l = len(self.data)
        for j in indicies[:-1]:
            k = j + 1
            w = l - k
            i = k - w
            if i < 0 or w < 3:
                continue
            if self.data[i:k] == self.data[k:l]:
                self.period = w
                self.base = i

    def predict(self, x):
        if x >= len(self.data):
            x -= self.base
            x %= self.period
            x += self.base
        return self.data[x]

m = Memo()
g = Grid(m)
for _ in range(10):
    g.step()

r = m.data[10]
print(f"Part 1: {r}")

while m.period == 0:
    g.step()

r = m.predict(1000000000)
print(f"Part 2: {r}")
