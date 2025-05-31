from collections import deque
import fileinput

dir: dict[str, tuple[int, int]] = {"N": (0, -1), "E": (1, 0), "S": (0, 1), "W": (-1, 0)}

def pipe_split(line):
    options: list[str] = []
    s = 0
    d = 0
    for i, ch in enumerate(line):
        if ch == "(":
            d += 1
        elif ch == ")":
            d -= 1
        if d == 0 and ch == "|":
            options.append(line[s:i])
            s = i + 1
    options.append(line[s:])
    return options

def block_split(line: str):
    d = 0
    for i, ch in enumerate(line):
        if ch == "(":
            d += 1
        elif ch == ")":
            d -= 1
        if d == 0:
            return [p + line[i+1:] for p in pipe_split(line[1:i])]
    return [line]

class Grid:
    def __init__(self, route: str):
        grid: dict[tuple[int, int], set[tuple[int, int]]] = {}
        visited: dict[str, set[tuple[int, int]]] = {}
        q = [((0, 0), r) for r in pipe_split(route[1:-1])]
        while len(q) > 0:
            v, r = q.pop()
            x, y = v
            if r not in visited:
                visited[r] = set()
            been = visited[r]
            if v in been:
                continue
            been.add(v)
            if v not in grid:
                grid[v] = set()
            adjs = grid[v]
            if len(r) > 0:
                if r[0] in dir:
                    dx, dy = dir[r[0]]
                    n = (x + dx, y + dy)
                    adjs.add(n)
                    if n not in grid:
                        grid[n] = set()
                    grid[n].add(v)
                    q.append((n, r[1:])) # type: ignore
                else:
                    q.extend([(v, o) for o in block_split(r)])
        self.grid = grid

    def find_furthest(self):
        been = {}
        q = deque()
        q.append(((0, 0), 0))
        while len(q) > 0:
            v, c = q.popleft()
            for a in self.grid[v]:
                if a not in been:
                    been[a] = c+1
                    q.append((a, c+1))
        return been.values()

g = Grid(next(fileinput.input()).strip())
doors = g.find_furthest()
print(f"Part 1: {max(doors)}")
print(f"Part 2: {sum(d > 1000 for d in doors)}")
