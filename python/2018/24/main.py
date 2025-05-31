from dataclasses import dataclass
import fileinput
import re

@dataclass
class Army:
    group_no: int
    units: int
    hit_points: int
    attack: int
    attack_type: str
    initiative: int
    def_mod: dict[str, int]

    def effective_power(self):
        return self.attack * self.units

    def get_mod(self, elem):
        return self.def_mod[elem] if elem in self.def_mod else 1

class War:
    def __init__(self, boost=0):
        lines = [l.strip() for l in fileinput.input()]
        mid = lines.index("")
        allies = {i+1: self._parse_army(i+1, l) for i, l in enumerate(lines[1:mid])}
        axis = {i+1: self._parse_army(i+1, l) for i, l in enumerate(lines[mid+2:])}
        self.all = [(lines[0][:-1], allies), (lines[mid+1][:-1], axis)]
        for a in allies.values():
            a.attack += boost

    def _parse_army(self, groupno, line):
        n = list(map(int, re.findall(r'\d+', line)))
        type = line.split()[-5]
        mods = {}
        mod_list = re.findall(r'\([^)]+\)', line)
        if len(mod_list) == 1:
            for l in mod_list[0][1:-1].split(";"):
                words = l.replace(",", "").split()
                for t in words[2:]:
                    mods[t] = 2 if words[0] == "weak" else 0
        return Army(groupno, n[0], n[1], n[2], type, n[3], mods)

    def _target_selection(self, a_idx, d_idx):
        out = []
        a_alliance, a_army = self.all[a_idx]
        attackers: list[Army] = sorted(a_army.values(), reverse=True,
            key=lambda g: (g.effective_power(), g.initiative))
        defenders: list[Army] = self.all[d_idx][1].values()
        selected = set()
        for a in attackers:
            targets = []
            for d in defenders:
                if d.units <= 0 or d.group_no in selected:
                    continue
                damage = a.effective_power() * d.get_mod(a.attack_type)
                if damage == 0:
                    continue
                targets.append((d.group_no, damage, d.effective_power(), d.initiative))
            if len(targets) == 0:
                continue
            targets.sort(key=lambda x: x[1:], reverse=True)
            target = targets[0][0]
            selected.add(target)
            out.append((a.initiative, a_idx, a.group_no, target))
        return out

    def attacking(self, t):
        _, a_idx, a_group_no, d_group_no = t
        d_idx = (a_idx + 1) % 2
        a_name, a_army = self.all[a_idx]
        d_name, d_army = self.all[d_idx]
        a: Army = a_army[a_group_no]
        d: Army = d_army[d_group_no]
        if a.units <= 0 or d.units <= 0:
            return 0
        damage = a.effective_power() * d.get_mod(a.attack_type)
        units_lost = min(damage // d.hit_points, d.units)
        d.units -= units_lost
        return units_lost

    def fight(self):
        for _, a in self.all:
            if len(a) == 0:
                return False
        attack_order = [*self._target_selection(1, 0), *self._target_selection(0, 1)]
        attack_order.sort(reverse=True)

        dead_guys = sum(self.attacking(t) for t in attack_order)
        if dead_guys == 0:
            return False

        for _, armies in self.all:
            for k in set(k for k, v in armies.items() if v.units <= 0):
                del armies[k]
        return True

    def winner(self):
        an, a = self.all[0]
        bn, b = self.all[1]
        if len(a) == 0:
            return f"{bn} wins"
        if len(b) == 0:
            return f"{an} wins"
        return "Undecided"

    def units_remaining(self):
        return sum(a.units for _, s in self.all for a in s.values())

def simulate(boost=0):
    w = War(boost)
    while w.fight():
        pass
    return w.winner(), w.units_remaining()

def binsearch():
    lo = hi = i = 1
    rem = 0
    while hi == 1:
        i <<= 3
        w, rem = simulate(i)
        print(f"Trying boost of {i}...{w} with {rem} units remaining")
        if w == "Immune System wins":
            hi = i
        else:
            lo = i

    while hi - lo > 1:
        i = (hi + lo) >> 1
        w, u = simulate(i)
        print(f"Trying boost of {i}...{w} with {u} units remaining")
        if w == "Immune System wins":
            rem = u
            hi = i
        else:
            lo = i
    return rem

_, u = simulate()
print(f"Part 1: {u}")

u = binsearch()
print(f"Part 2: {u}")
