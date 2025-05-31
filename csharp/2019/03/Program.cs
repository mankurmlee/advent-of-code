class Puzzle
{
    public Wire Wire1 { get; }
    public Wire Wire2 { get; }

    public Puzzle(string filename)
    {
        string[] lines = File.ReadAllLines(filename);
        Wire1 = new(lines[0].Split(","));
        Wire2 = new(lines[1].Split(","));
    }

    public List<(int, int)> Intersections()
    {
        List<(int, int)> res = [];
        foreach (var s1 in Wire1.Sections)
        {
            foreach (var s2 in Wire2.Sections)
            {
                var (ok, pos) = s1.Intersects(s2);
                if (ok) {
                    res.Add(pos);
                }
            }
        }
        return res;
    }

    public int TaxiDist(bool partOne)
    {
        int min = 0;
        int d = 0;
        foreach (var pos in Intersections())
        {
            if (pos == (0, 0)) continue;
            if (partOne)
            {
                d = Math.Abs(pos.Item1) + Math.Abs(pos.Item2);
            }
            else
            {
                d = Wire1.Dist(pos) + Wire2.Dist(pos);
            }
            if (min == 0 || d < min)
            {
                min = d;
            }
        }
        return min;
    }
}

class Wire
{
    public Section[] Sections { get; }

    public Wire(string[] path)
    {
        int n = path.Length;
        Sections = new Section[n];

        var pos = (0, 0);
        for (int i = 0; i < n; i++)
        {
            string p = path[i];
            char direction = p[0];
            int steps = int.Parse(p[1..]);

            var last = pos;
            Section s = new((0, 0));
            switch (direction)
            {
                case 'U':
                    pos.Item2 += steps;
                    s.Lo = last;
                    s.Hi = pos;
                    break;
                case 'D':
                    pos.Item2 -= steps;
                    s.Lo = pos;
                    s.Hi = last;
                    break;
                case 'L':
                    pos.Item1 -= steps;
                    s.Lo = pos;
                    s.Hi = last;
                    break;
                case 'R':
                    pos.Item1 += steps;
                    s.Lo = last;
                    s.Hi = pos;
                    break;
            }
            Sections[i] = s;
        }
    }

    public int Dist((int, int) pos)
    {
        var start = (0, 0);
        (int, int)[] last = [start];
        int d = 0;
        Section spos = new(pos);
        foreach (var s in Sections)
        {
            var (ok, _) = s.Intersects(spos);
            if (!ok)
            {
                d += GetDist(s.Lo, s.Hi);
                last = [s.Lo, s.Hi];
                continue;
            }
            if (last.Contains(s.Lo))
            {
                start = s.Lo;
            }
            else
            {
                start = s.Hi;
            }
            d += GetDist(start, pos);
            break;
        }
        return d;
    }

    private static int GetDist((int, int) a, (int, int) b)
    {
        return Math.Abs(a.Item1 - b.Item1) + Math.Abs(a.Item2 - b.Item2);
    }
}

class Section((int, int) point)
{
    public (int, int) Lo { get; set; } = point;
    public (int, int) Hi { get; set; } = point;

    public (bool, (int, int)) Intersects(Section s2)
    {
        if (
            (Hi.Item1 < s2.Lo.Item1) ||
            (s2.Hi.Item1 < Lo.Item1) ||
            (Hi.Item2 < s2.Lo.Item2) ||
            (s2.Hi.Item2 < Lo.Item2)
        )
        {
            return (false, (0, 0));
        }

        if (Lo.Item1 == Hi.Item1) {
            return (true, (Lo.Item1, s2.Lo.Item2));
        }
        return (true, (s2.Lo.Item1, Lo.Item2));
    }
}

class Program
{
    static void Main(string[] args)
    {
        string filename = args[0];
        Puzzle p = new(filename);

        Console.WriteLine($"Part 1: {p.TaxiDist(true)}");
        Console.WriteLine($"Part 2: {p.TaxiDist(false)}");
    }
}
