class MoonSystem
{
    private readonly List<Moon> moons = [];
    public MoonSystem(string filename)
    {
        string[] lines = File.ReadAllLines(filename);
        foreach (string line in lines)
        {
            string[] data = line.Replace("=", " ")
                .Replace(",", " ")
                .Replace(">", " ")
                .Split(" ", StringSplitOptions.RemoveEmptyEntries);
            int x = int.Parse(data[1]);
            int y = int.Parse(data[3]);
            int z = int.Parse(data[5]);
            moons.Add(new(x, y, z));
        }
    }

    public override string ToString()
    {
        List<string> txt = [];
        foreach (Moon m in moons)
            txt.Add(m.ToString());

        return string.Join(Environment.NewLine, txt);
    }

    public void Step()
    {
        foreach (Moon m in moons)
            m.Acc = CalcAcc(m);

        foreach (Moon m in moons)
            m.Vel.AddEquals(m.Acc);

        foreach (Moon m in moons)
            m.Pos.AddEquals(m.Vel);
    }

    private Vec CalcAcc(Moon m)
    {
        Vec acc = new();
        foreach (Moon other in moons)
        {
            if (m == other) continue;
            acc.X += other.Pos.X.CompareTo(m.Pos.X);
            acc.Y += other.Pos.Y.CompareTo(m.Pos.Y);
            acc.Z += other.Pos.Z.CompareTo(m.Pos.Z);
        }
        return acc;
    }

    public int Energy()
    {
        int energy = 0;
        foreach (Moon m in moons)
            energy += m.Energy();
        return energy;
    }

    public long GetCycle(int part1)
    {
        Memo xHist = new();
        Memo yHist = new();
        Memo zHist = new();

        int n = moons.Count;
        int m = n << 1;
        int[] xs = new int[m];
        int[] ys = new int[m];
        int[] zs = new int[m];

        int count = 0;
        while (xHist.Cycle == 0 || yHist.Cycle == 0 || zHist.Cycle == 0)
        {
            for (int i = 0; i < n; i++)
            {
                int j = i << 1;
                xs[j]   = moons[i].Pos.X;
                ys[j]   = moons[i].Pos.Y;
                zs[j]   = moons[i].Pos.Z;
                xs[j+1] = moons[i].Vel.X;
                ys[j+1] = moons[i].Vel.Y;
                zs[j+1] = moons[i].Vel.Z;
            }

            xHist.Add(xs);
            yHist.Add(ys);
            zHist.Add(zs);
            Step();
            count++;
            if (count == part1)
            {
                Console.WriteLine($"Part 1: {Energy()}");
            }
        }

        return LCM(LCM(xHist.Cycle, yHist.Cycle), zHist.Cycle);
    }

    private static long LCM(long a, long b)
    {
        return a * b / GCD(a, b);
    }

    private static long GCD(long a, long b)
    {
        while (b != 0)
            (a, b) = (b, a % b);
        return a;
    }
}

class Moon(int x, int y, int z)
{
    public Vec Pos = new(x, y, z);
    public Vec Vel = new();
    public Vec Acc = new();
    public override string ToString()
    {
        return $"pos = {Pos}, vel = {Vel}, acc = {Acc}";
    }

    public int Energy()
    {
        return Pos.Energy() * Vel.Energy();
    }
}

struct Vec(int x, int y, int z)
{
    public int X = x;
    public int Y = y;
    public int Z = z;
    public Vec(): this(0, 0, 0) {}
    public override readonly string ToString()
    {
        return $"({X}, {Y}, {Z})";
    }

    public void AddEquals(Vec other)
    {
        X += other.X;
        Y += other.Y;
        Z += other.Z;
    }

    public readonly int Energy()
    {
        return Math.Abs(X) + Math.Abs(Y) + Math.Abs(Z);
    }
}

class Memo
{
    private readonly List<int[]> history = [];
    private int idx = 0;
    public int Cycle = 0;

    public void Add(int[] item)
    {
        if (Cycle > 0) return;

        history.Add([.. item]);
        if (history.Count == 1) return;

        if (!ArrayEquals(item, history[idx]))
        {
            idx = 0;
            return;
        }

        idx++;
        if ((idx << 1) < history.Count) return;

        Cycle = idx;
    }

    private static bool ArrayEquals(int[] a, int[] b)
    {
        for (int i = 0; i < 8; i++)
            if (a[i] != b[i])
                return false;

        return true;
    }
}

class Program
{
    static void Main(string[] args)
    {
        MoonSystem s = new(args[0]);
        int energyCycles = int.Parse(args[1]);
        long period = s.GetCycle(energyCycles);
        Console.WriteLine($"Part 2: {period}");
    }
}
