import fileinput
import re

class Reservoir:
    def __init__(self):
        walls = set()
        for line in fileinput.input():
            words = re.findall(r'[0-9xy]+', line)
            j = int(words[1])
            for i in range(int(words[3]), int(words[4])+1):
                w = (j, i) if words[0] == 'x' else (i, j)
                walls.add(w)
        self.x_min = min(x for x, _ in walls)
        self.x_max = max(x for x, _ in walls)
        self.y_min = min(y for _, y in walls)
        self.y_max = max(y for _, y in walls)
        self.walls = walls
        self.unsettled = set()
        self.settled = set()

    def draw(self):
        for y in range(self.y_min, self.y_max+1):
            line = []
            for x in range(self.x_min-2, self.x_max+3):
                v = (x, y)
                ch = '.'
                if v in self.walls:
                    ch = '#'
                elif v in self.unsettled:
                    ch = '|'
                elif v in self.settled:
                    ch = '~'
                line.append(ch)
            print("".join(line))

    def flow_down(self, x_start, y_start):
        y = y_start
        while True:
            y += 1
            v = (x_start, y)
            if y > self.y_max or v in self.unsettled:
                return True
            if v in self.walls or v in self.settled:
                break
            if y >= self.y_min:
                self.unsettled.add(v)

        while y > y_start + 1:
            y -= 1
            if (self.flow_sides(x_start, y)):
                return True

        return False

    def flow_sides(self, x_start, y_start):
        unsettled = False
        x_min = x_start
        x_max = x_start
        for d in [-1, 1]:
            x = x_start
            while True:
                x += d
                if (x, y_start) in self.walls:
                    x_min = min(x_min, x + 1)
                    x_max = max(x_max, x - 1)
                    break
                self.unsettled.add((x, y_start))
                below = (x, y_start+1)
                if below not in self.walls and below not in self.settled:
                    if self.flow_down(x, y_start):
                        unsettled = True
                        break

        if not unsettled:
            for x in range(x_min, x_max+1):
                self.unsettled.remove((x, y_start))
                self.settled.add((x, y_start))

        return unsettled

r = Reservoir()
r.flow_down(500, 0)
# r.draw()
print(f"Part 1: {len(r.settled) + len(r.unsettled)}")
print(f"Part 2: {len(r.settled)}")
