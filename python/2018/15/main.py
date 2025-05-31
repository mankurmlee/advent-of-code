import fileinput
from dataclasses import dataclass
from collections import deque

@dataclass
class Vec:
    x: int
    y: int

    def __add__(self, other):
        return Vec(self.x + other.x, self.y + other.y)

    def __hash__(self):
        return hash((self.x, self.y))

    def distance_to(self, other):
        return abs(other.x - self.x) + abs(other.y - self.y)

class Dirs:
    ALL = [Vec(0, -1), Vec(-1, 0), Vec(1, 0), Vec(0, 1)]

@dataclass
class Unit:
    sprite: str
    pos: Vec
    hp: int
    atk: int

    def take_turn(self, game):
        if self.hp <= 0:
            return "OK"

        enemy = "E" if self.sprite != "E" else "G"

        num_enemies = len([u for u in game.units.values() if u.sprite == enemy])
        if num_enemies == 0:
            return "END"

        attackable = game.get_enemies_in_range(self, enemy)
        if len(attackable) == 0:
            targets = game.get_move_targets(enemy)
            if len(targets) == 0:
                return "OK"
            paths = game.find_best_paths(self.pos, targets)
            if len(paths) == 0:
                return "OK"
            step = self.decide_next_step(paths)
            del game.units[self.pos]
            self.pos = step
            game.units[self.pos] = self
            attackable = game.get_enemies_in_range(self, enemy)
            if len(attackable) == 0:
                return  "OK"

        # choose weakest enemy
        attackable.sort()
        weakest = attackable[0]
        weakest.damage(game, self.atk)
        if weakest.sprite == "E" and weakest.hp <= 0:
            return "DEADELF"
        return "OK"

    def damage(self, game, dmg):
        self.hp -= dmg
        if self.hp <= 0:
            del game.units[self.pos]

    def decide_next_step(self, paths):
        if len(paths) == 1:
            return paths[0].first
        target = paths[0].pos
        filtered = [p.first for p in paths if p.pos == target]
        if len(filtered) == 1:
            return filtered[0]
        steps = sorted(filtered, key=lambda v: (v.y, v.x))
        return steps[0]

    def __str__(self):
        return f"{self.sprite}[{self.pos.x},{self.pos.y}]"

    def __lt__(self, other):
        if self.hp != other.hp:
            return self.hp < other.hp
        return reading_order(self.pos, other.pos)

class Path:
    def __init__(self, pos_, cost_, first_):
        self.pos = pos_
        self.cost = cost_
        self.first = first_

    def __lt__(self, other):
        if self.cost != other.cost:
            return self.cost < other.cost
        return reading_order(self.pos, other.pos)

    def __repr__(self):
        return f"<{self.pos}, {self.cost}, {self.first}>"

@dataclass
class Game:
    width: int
    height: int
    grid: list[str]
    units: dict[Vec, Unit]
    roundno = 0

    @classmethod
    def load(cls, elfattack):
        g = []
        u = {}

        h = 0
        for l in fileinput.input():
            l = l.strip()
            g.extend(l)
            h += 1

        w = len(g) // h
        for i, tile in enumerate(g):
            if tile != "E" and tile != "G":
                continue
            g[i] = "."
            v = Vec(i % w, i // w)
            atk = elfattack if tile == "E" else 3
            u[v] = Unit(tile, v, 200, atk)

        return cls(w, h, g, u)

    def draw(self):
        print(f"After {self.roundno} rounds:")
        g = list(self.grid)
        w = self.width
        for u in self.units.values():
            g[w*u.pos.y+u.pos.x] = u.sprite

        i = 0
        for y in range(self.height):
            print("".join(g[i:i+w]))
            i += w

    def play_round(self, part_two=False):
        units = sorted(self.units.values(), key=lambda u: (u.pos.y, u.pos.x))
        for u in units:
            res = u.take_turn(self)
            if res == "END" or part_two and res == "DEADELF":
                return res
        self.roundno += 1
        return "OK"

    def get_move_targets(self, sprite):
        in_range = []
        for u in self.units.values():
            if u.sprite != sprite:
                continue
            in_range.extend(self.get_adj(u.pos))
        return set(in_range)

    def get_adj(self, pos):
        adjs = []
        w = self.width
        used = self.units.keys()
        for d in Dirs.ALL:
            adj = pos + d
            if adj in used or self.grid[w * adj.y + adj.x] != ".":
                continue
            adjs.append(adj)
        return adjs

    def find_best_paths(self, start, dests):
        w = self.width
        q = deque([Path(start, 0, None)])

        hist = {}
        best_dist = 2<<32
        best_paths = []
        while (len(q) > 0):
            p = q.popleft()
            if p.cost > best_dist:
                break
            if p.pos in dests:
                if p.cost < best_dist:
                    best_dist = p.cost
                    best_paths = []
                best_paths.append(p)
                continue
            for d in Dirs.ALL:
                adj_pos = p.pos + d
                if self.grid[w*adj_pos.y+adj_pos.x] != ".":
                    continue
                if adj_pos in self.units:
                    continue
                adj_cost = p.cost + 1
                stepone = adj_pos if p.first == None else p.first
                if adj_pos in hist:
                    hcost, hstep = hist[adj_pos]
                    if adj_cost > hcost:
                        continue
                    if adj_cost == hcost and not reading_order(stepone, hstep):
                        continue
                hist[adj_pos] = adj_cost, stepone
                adj = Path(adj_pos, adj_cost, stepone)
                q.append(adj)
        return best_paths

    def get_enemies_in_range(self, unit, enemy):
        in_range = []
        for d in Dirs.ALL:
            adj = d + unit.pos
            if adj not in self.units:
                continue
            other = self.units[adj]
            if other.sprite != enemy:
                continue
            in_range.append(other)
        return in_range

    def total_hp(self):
        return sum(u.hp for u in self.units.values())

    def get_outcome(self):
        return self.roundno * self.total_hp()

def reading_order(a: Vec, b: Vec):
    if a.y != b.y:
        return a.y < b.y
    return a.x < b.x

def play_game(part_two=False):
    elfattack = 3 if part_two == False else 4
    res = "OK"
    while True:
        g = Game.load(elfattack)
        while True:
            res = g.play_round(part_two)
            if res != "OK":
                break
        if res == "END":
            break
        elfattack += 1
    g.draw()
    print(f"Final Elf Attack: {elfattack}")
    return g.get_outcome()

print(f"Part 1: {play_game()}")
print(f"Part 2: {play_game(True)}")

