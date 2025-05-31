from heapq import heappush, heappop
import sys

class Puzzle:
    def __init__(self, unlocks, locked_by, workers, delay):
        self.unlocks: dict[str, str] = unlocks
        self.locked_by: dict[str, str] = locked_by
        self.workers = workers
        self.delay = delay

    @classmethod
    def from_file(cls, filename):
        unlocks = {}
        locked_by = {}
        with open(filename) as f:
            workers = int(f.readline().strip())
            delay = int(f.readline().strip())
            for l in f:
                l = l.strip()
                if l == "":
                    continue
                words = l.split()
                k = words[7]
                v = words[1]

                if k not in locked_by:
                    locked_by[k] = []
                locked_by[k].append(v)

                if v not in unlocks:
                    unlocks[v] = []
                unlocks[v].append(k)

        return cls(unlocks, locked_by, workers, delay)

    def _initial_options(self):
        options = [k for k in self.unlocks.keys()]
        for k in self.locked_by.keys():
            if k in options:
                options.remove(k)

        pq = []
        for o in options:
            heappush(pq, o)
        return pq

    def single_task_order(self):
        result = []
        options = self._initial_options()
        while len(options) > 0:
            l = heappop(options)
            result.append(l)
            if l not in self.unlocks:
                continue
            for a in self.unlocks[l]:
                if self._is_open(a, result):
                    heappush(options, a)
        return "".join(result)

    def _is_open(self, letter, opened):
        for l in self.locked_by[letter]:
            if l not in opened:
                return False
        return True

    def multi_task_time(self):
        t = 0
        done = []
        processing = []
        options = self._initial_options()
        while len(options) > 0 or len(processing) > 0:
            while len(options) > 0 and len(processing) < self.workers:
                l = heappop(options)
                dt = t + ord(l) - ord('A') + 1 + self.delay
                heappush(processing, (dt, l))
            dt, l = heappop(processing)
            done.append(l)
            t = dt
            if l not in self.unlocks:
                continue
            for a in self.unlocks[l]:
                if self._is_open(a, done):
                    heappush(options, a)
        order = "".join(done)
        print(f"Multi-task order: {order}")
        return t

def main():
    p = Puzzle.from_file(sys.argv[1])

    o = p.single_task_order()
    print(f"Part 1: {o}")

    t = p.multi_task_time()
    print(f"Part 2: {t}")

if __name__ == "__main__":
    main()
