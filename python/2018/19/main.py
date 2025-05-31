import fileinput

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

functions = {
	"addr": addr, "addi": addi, "mulr": mulr, "muli": muli, "banr": banr, "bani": bani, "borr": borr, "bori": bori,
    "setr": setr, "seti": seti, "gtir": gtir, "gtri": gtri, "gtrr": gtrr, "eqir": eqir, "eqri": eqri, "eqrr": eqrr,
}

def run(part_two = False):
	r = [0] * 6
	if part_two:
		r[0] = 1
	ip = 0
	while ip >= 0 and ip < len(lines):
		if part_two and ip == 11:
			return r[2]
		f = functions[lines[ip].split()[0]]
		a, b, c = map(int, lines[ip].split()[1:])
		r[i] = ip
		f(r, a, b, c)
		ip = r[i] + 1
	return r[0]

lines = [l.strip() for l in fileinput.input()]
i = int(lines[0].split()[1])
lines = lines[1:]

print(f"Part 1: {run()}")

num = run(True)
a = 0
for i in range(1, num+1):
	if num % i == 0:
		a += i
print(f"Part 2: {a}")