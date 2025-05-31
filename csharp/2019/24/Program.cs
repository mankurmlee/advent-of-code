class Grid
{
    public bool[] Data;

    public Grid(string filename)
    {
        Data = new bool[25];
        string[] lines = File.ReadAllLines(filename);

        for (int y = 0; y < 5; y++)
            for (int x = 0; x < 5; x++)
                Data[y*5+x] = lines[y][x] == '#';
    }

    public override string ToString()
    {
        List<string> lines = [];
        for (int i = 0; i < 25; i += 5)
        {
            char[] c = new char[5];
            for (int j = 0; j < 5; j++)
                c[j] = Data[i+j] ? '#' : '.';
            lines.Add(new string(c));
        }
        return string.Join(Environment.NewLine, lines);
    }

    public int ToInt()
    {
        int n = 0;
        int x = 1;
        foreach (bool b in Data)
        {
            if (b)
                n |= x;
            x <<= 1;
        }
        return n;
    }

    public void Update()
    {
        bool[] next = new bool[25];
        for (int y = 0; y < 5; y++)
            for (int x = 0; x < 5; x++)
                next[y*5+x] = UpdateTile(x, y);
        Data = next;
    }

    private bool UpdateTile(int x, int y)
    {
        bool b = Data[y * 5 + x];
        int a = CountAdj(x, y);
        if (b && a != 1)
            return false;
        if (!b && (a == 1 || a == 2))
            return true;
        return b;
    }

    private int CountAdj(int x, int y)
    {
        int n = 0;
        foreach (Vec d in Vec.DIR)
        {
            int nx = x + d.x;
            int ny = y + d.y;
            if (nx < 0 || ny < 0 || nx >= 5 || ny >= 5) continue;
            if (Data[ny * 5 + nx])
                n++;
        }
        return n;
    }

    public int Rating()
    {
        HashSet<int> seen = [];

        int rating = ToInt();
        do
        {
            seen.Add(rating);
            Console.WriteLine(this);
            Console.WriteLine(rating);
            Update();
            rating = ToInt();
        } while (!seen.Contains(rating));

        return rating;
    }
}

class InfGrid
{
    private Dictionary<int,bool[]> Tesseract = [];
    private int Top = 0;
    private int Bot = 0;
    private readonly Dictionary<Vec,List<Vec>> Links = [];
    public InfGrid(Grid g)
    {
        Tesseract.Add(0,g.Data);
        for (int i = 0; i < 5; i++)
        {
            AddLink(new Vec(2, 1), new Vec(i, 4));
            AddLink(new Vec(3, 2), new Vec(0, i));
            AddLink(new Vec(2, 3), new Vec(i, 0));
            AddLink(new Vec(1, 2), new Vec(4, i));
        }
    }
    private void AddLink(Vec vec1, Vec vec2)
    {
        if (!Links.TryGetValue(vec1, out List<Vec>? l1))
        {
            l1 = [];
            Links[vec1] = l1;
        }
        l1.Add(vec2);
        if (!Links.TryGetValue(vec2, out List<Vec>? l2))
        {
            l2 = [];
            Links[vec2] = l2;
        }
        l2.Add(vec1);
    }

    public void Update()
    {
        Dictionary<int,bool[]> result = [];
        for (int i = Bot - 1; i <= Top + 1; i++)
        {
            bool[] g = new bool[25];
            bool bugged = false;
            for (int y = 0; y < 5; y++)
                for (int x = 0; x < 5; x++)
                {
                    if (x == 2 && y == 2) continue;
                    bool tile = UpdateTile(i, new(x, y));
                    g[y*5+x] = tile;
                    bugged |= tile;
                }

            if (!bugged) continue;

            result[i] = g;
            if (i < Bot)
                Bot = i;
            if (i > Top)
                Top = i;
        }
        Tesseract = result;
    }

    private bool UpdateTile(int i, Vec pos)
    {
        int a = CountAdj(i, pos);
        bool b = false;
        if (Tesseract.TryGetValue(i, out bool[]? g))
            b = g[pos.y * 5 + pos.x];

        if (b && a != 1)
            return false;
        if (!b && (a == 1 || a == 2))
            return true;
        return b;
    }

    private int CountAdj(int i, Vec pos)
    {
        int n = 0;
        if (Tesseract.TryGetValue(i, out bool[]? g))
        {
            foreach (Vec d in Vec.DIR)
            {
                Vec a = pos.Add(d);
                if (a.x < 0 || a.y < 0 || a.x >= 5 || a.y >= 5) continue;
                if (g[a.y * 5 + a.x])
                    n++;
            }
        }
        if (Links.TryGetValue(pos, out List<Vec>? conns))
        {
            int j = conns.Count == 5 ? i + 1 : i - 1;
            if (Tesseract.TryGetValue(j, out bool[]? grid))
                foreach (Vec a in conns)
                    if (grid[a.y * 5 + a.x])
                        n++;
        }
        return n;
    }

    public void Print()
    {
        for (int i = Bot; i <= Top; i++)
        {
            Console.WriteLine($"Depth {i}:");
            for (int y = 0; y < 5; y++)
            {
                char[] c = new char[5];
                for (int x = 0; x < 5; x++)
                    c[x] = Tesseract[i][y*5+x] ? '#' : '.';
                if (y == 2) c[2] = '?';
                string s = new(c);
                Console.WriteLine(s);
            }
        }
    }

    public int Count()
    {
        int n = 0;
        foreach (bool[] g in Tesseract.Values)
            foreach (bool b in g)
                if (b)
                    n++;
        return n;
    }
}

class Program
{
    static void Main(string[] args)
    {
        string filename = args[0];
        // string filename = @"D:\aoc\2019\24\sample.txt";
        Grid g = new(filename);
        InfGrid h = new(g);

        int rating = g.Rating();
        Console.WriteLine($"Part 1: {rating}");

        for (int i = 0; i < 200; i++)
            h.Update();

        int bugs = h.Count();
        Console.WriteLine($"Part 2: {bugs}");
    }
}
