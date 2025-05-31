

class ASCII
{
    private readonly static string[] Path = [
            "A,B,A,C,A,B,C,A,B,C",
            "R,12,R,4,R,10,R,12",
            "R,6,L,8,R,10",
            "L,8,R,4,R,4,R,6",
            "n",
            ""
    ];

    private readonly Intcode prog;
    private readonly Intcode progTwo;
    public Scaffold Map;

    public ASCII(string filename)
    {
        prog = new(filename);
        progTwo = prog.Clone();
        progTwo.Poke(0, 2);
        Map = new(prog);
        Console.WriteLine(Map);
    }

    public int SumAlignmentParameters()
    {
        int sum = 0;
        foreach (Vec v in Map.GetIntersections())
            sum += v.x * v.y;
        return sum;
    }

    public int SendPath()
    {
        string input = string.Join("\n", Path);
        foreach (char c in input)
            progTwo.Push(c);

        int ret = 0;
        while (true)
        {
            var (halted, val) = progTwo.Execute();
            if (halted) break;
            ret = (int)val;
        }
        return ret;
    }
}

class Scaffold
{
    public int Width;
    public int Height;
    public char[] Data;
    public Scaffold(Intcode prog)
    {
        int w = 0;
        int h = 0;
        List<char> d = [];
        List<char> row = [];
        bool done = false;
        long v = 0;
        while (true)
        {
            (done, v) = prog.Execute();
            if (done) break;
            if (v == 10)
            {
                if (row.Count == 0) continue;
                d.AddRange(row);
                row = [];
                if (w == 0)
                    w = d.Count;
                h++;
                continue;
            }
            row.Add((char)v);
        }
        Width = w;
        Height = h;
        Data = [.. d];
    }

    public override string ToString()
    {
        List<string> s = [];
        s.Add($"Width = {Width}");
        s.Add($"Height = {Height}");
        for (int i = 0; i < Data.Length; i += Width)
        {
            int j = i + Width;
            s.Add(new string(Data[i..j]));
        }
        return string.Join(Environment.NewLine, s);
    }

    public List<Vec> GetIntersections()
    {
        List<Vec> ret = [];
        int off = 0;
        for (int y = 0; y < Height; y++)
            for (int x = 0; x < Width; x++)
            {
                if (Data[off] != '.')
                {
                    Vec v = new(x, y);
                    int adj = CountAdjScaff(v);
                    if (adj > 2) ret.Add(v);
                }
                off++;
            }
        return ret;
    }

    private int CountAdjScaff(Vec v)
    {
        int sum = 0;
        foreach (Vec dir in Vec.ADJ)
        {
            Vec a = v.Add(dir);
            if (a.x < 0 || a.y < 0 || a.x >= Width || a.y >= Height) continue;
            if (Data[a.y * Width + a.x] == '.') continue;
            sum++;
        }
        return sum;
    }
}

struct Vec(int x, int y)
{
    public static readonly Vec[] ADJ = [
        new Vec(0, -1), new Vec(0, 1), new Vec(-1, 0), new Vec(1, 0),
    ];
    public int x = x;
    public int y = y;
    public readonly Vec Add(Vec other)
    {
        return new Vec(x + other.x, y + other.y);
    }
}

class Program
{
    static void Main(string[] args)
    {
        string filename = "D:\\aoc\\2019\\17\\puzzle.txt";
        ASCII a = new(filename);
        int res = a.SumAlignmentParameters();
        Console.WriteLine($"Part 1: {res}");

        int dust = a.SendPath();
        Console.WriteLine($"Part 2: {dust}");
    }
}