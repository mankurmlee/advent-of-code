import sys
import multiprocessing
from multiprocessing import Pool, Array

class Grid:
    def __init__(self, filename):
        s = 0
        with open(filename) as f:
            s = int(next(f).strip())

        self.serial = s
        cells = [0] * 300 * 300
        for y in range(300):
            for x in range(300):
                cells[y*300+x] = self._power_level(x+1, y+1)

        self.cells = cells

    def _power_level(self, x: int, y: int):
        rack_id = x + 10
        p = rack_id * y + self.serial
        return ((p * rack_id) // 100) % 10 - 5

def best_for_size(cells, size: int):
    bx = 0
    by = 0
    bp = 0
    n = 301 - size
    memo = [0] * 300
    for x in range(n):
        # create sq
        p = 0
        s = x
        for y in range(size):
            memo[y] = sum(cells[s:s+size])
            p += memo[y]
            s += 300

        # test sq
        if p > bp:
            bp = p
            bx, by = x, 0

        for y in range(size, 300):
            # move sq
            memo[y] = sum(cells[s:s+size])
            p += memo[y] - memo[y-size]
            s += 300

            # test sq
            if p > bp:
                bp = p
                bx, by = x, y+1-size

    return bx+1, by+1, bp

def mod_best(cells, mod, rem):
    most = 0
    coord = 0, 0, 0, 0
    for i in range(1, 300):
        if i % mod != rem:
            continue
        (x, y, pow) = best_for_size(cells, i)
        if pow > most:
            most = pow
            coord = x, y, i, pow
    return coord

def init_worker(shared_array_, num_workers_):
    global shared_array, num_workers
    shared_array = shared_array_
    num_workers = num_workers_

def worker(index):
    return mod_best(shared_array, num_workers, index)

def best(cells):
    shared_array = Array('i', cells, lock=False)
    num = multiprocessing.cpu_count()
    with Pool(num, init_worker, (shared_array, num)) as pool:
        results = pool.map(worker, list(range(num)))

    best = 0
    res = 0, 0, 0
    for r in results:
        x, y, i, p = r
        if p > best:
            best = p
            res = x, y, i
    return res

def main():
    g = Grid(sys.argv[1])

    (x, y, _) = best_for_size(g.cells, 3)
    print(f"Part 1: {x},{y}")

    (x, y, s) = best(g.cells)
    print(f"Part 2: {x},{y},{s}")

if __name__ == "__main__":
    main()
