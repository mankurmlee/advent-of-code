import fileinput
import re

def addr(r, a, b, c): r[c] = r[a] + r[b]
def addi(r, a, b, c): r[c] = r[a] + b
def mulr(r, a, b, c): r[c] = r[a] * r[b]
def muli(r, a, b, c): r[c] = r[a] * b
def banr(r, a, b, c): r[c] = r[a] & r[b]
def bani(r, a, b, c): r[c] = r[a] & b
def borr(r, a, b, c): r[c] = r[a] | r[b]
def bori(r, a, b, c): r[c] = r[a] | b
def setr(r, a, b, c): r[c] = r[a]
def seti(r, a, b, c): r[c] = a
def gtir(r, a, b, c): r[c] = int(a > r[b])
def gtri(r, a, b, c): r[c] = int(r[a] > b)
def gtrr(r, a, b, c): r[c] = int(r[a] > r[b])
def eqir(r, a, b, c): r[c] = int(a == r[b])
def eqri(r, a, b, c): r[c] = int(r[a] == b)
def eqrr(r, a, b, c): r[c] = int(r[a] == r[b])

functions = [
    addr, addi, mulr, muli, banr, bani, borr, bori,
    setr, seti, gtir, gtri, gtrr, eqir, eqri, eqrr,
]

def parse(line):
    return list(map(int, re.findall(r'\d+', line)))

def behaves_like(instruction, before, after):
    count = 0
    for f in functions:
        r = list(before)
        f(r, *instruction[1:])
        if r == after:
            count += 1
    return count

def candidates(f):
    res = [i for i in range(16) if i not in found]
    for instruction, before, after in tests:
        opcode = instruction[0]
        if opcode in res:
            r = list(before)
            f(r, *instruction[1:])
            if r != after:
                res.remove(opcode)
    return res

lines = list(fileinput.input())

tests = []
before = []
instruction = []
count = 0
for line in lines:
    if "Before" in line:
        before = parse(line)
    elif "After" in line:
        after = parse(line)
        test = (instruction, before, after)
        tests.append(test)
        if behaves_like(*test) >= 3:
            count += 1
    else:
        instruction = parse(line)

print(f"Part 1: {count}")

found = {}
while len(found) < len(functions):
    for f in functions:
        if f not in found.values():
            c = candidates(f)
            if len(c) == 1:
                found[c[0]] = f

r = [0] * 4
i = max(i for i, x in enumerate(lines) if not x.strip()) + 1
for line in lines[i:]:
    op, a, b, c = parse(line)
    found[op](r, a, b, c)

print(f"Part 2: {r[0]}")