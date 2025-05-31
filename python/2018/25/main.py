import fileinput
import re

def distance_to(a, b):
    return sum(abs(p - q) for p, q in zip(a, b))

def links_to(s, c):
    for a in starlinks[s]:
        if a in c:
            return True
    return False

def add_star(cons, s):
    adj = [i for i, c in enumerate(cons) if links_to(s, c)]
    if len(adj) == 0:
        cons.append({s})
        return
    c = cons[adj[0]]
    c.add(s)
    for i in reversed(adj[1:]):
        c.update(cons[i])
        del cons[i]

stars = list(tuple(map(int, re.findall(r'-?\d+', l))) for l in fileinput.input())
starlinks = {a: set(b for b in stars if a != b and distance_to(a, b) <= 3) for a in stars}
cons = []
for s in stars:
    add_star(cons, s)
print(f"Part 1: {len(cons)}")