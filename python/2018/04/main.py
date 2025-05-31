from functools import total_ordering
import sys
from datetime import datetime

@total_ordering
class Record:
    def __init__(self, dt: datetime, event: str):
        self.datetime: datetime = dt
        self.event: str = event

    @classmethod
    def Parse(cls, txt):
        date_string = "20" + txt[3:17]
        event = txt[19:]

        dt = datetime.strptime(date_string, "%Y-%m-%d %H:%M")
        return cls(dt, event)

    def __eq__(self, other):
        return self.datetime == other

    def __lt__(self, other):
        return self.datetime < other

    def __str__(self):
        return str(self.datetime) + " " + self.event


class Journal:
    def __init__(self, filename):
        self.records: list[Record] = []
        with open(filename) as f:
            data = f.read();
            lines = data.splitlines();
            for l in lines:
                l = l.strip()
                if l == "":
                    continue
                self.records.append(Record.Parse(l))
        self.records.sort()

    def __str__(self):
        return "\n".join(str(v) for v in self.records)

class NapLog:
    def __init__(self, filename):
        self.all_naps: dict[int,list[tuple[datetime,datetime]]] = {}

        j = Journal(filename)
        guard = 0
        asleep: datetime = datetime.now();
        for r in j.records:
            words = r.event.split();
            if words[0] == "Guard":
                guard = int(words[1][1:])
                if guard not in self.all_naps:
                    self.all_naps[guard] = []
            elif words[1] == "asleep":
                asleep = r.datetime
            elif words[1] == "up":
                self.all_naps[guard].append((asleep, r.datetime))

    def part_one(self):
        g = self.sleepiest()
        m = self.deepest(g)[0]
        return g * m

    def sleepiest(self):
        most = 0
        guard = 0

        for g, naps in self.all_naps.items():
            slept = 0
            for s, w in naps:
                slept += w.minute - s.minute
            if slept > most:
                most = slept
                guard = g

        return guard

    def deepest(self, guard) -> tuple[int, int]:
        minutes = {}
        for stime, etime in self.all_naps[guard]:
            for m in range(stime.minute, etime.minute):
                if m not in minutes:
                    minutes[m] = 0
                minutes[m] += 1

        freq = 0
        minute = 0
        for m, f in minutes.items():
            if f > freq:
                freq = f
                minute = m
        return (minute, freq)

    def part_two(self) -> int:
        res = 0
        most = 0
        for g in self.all_naps.keys():
            m, f = self.deepest(g)
            if f > most:
                most = f
                res = g * m
        return res

def main():
    l = NapLog(sys.argv[1])
    r = l.part_one();
    print(f"Part 1: {r}")

    r = l.part_two();
    print(f"Part 2: {r}")


if __name__ == "__main__":
    main()