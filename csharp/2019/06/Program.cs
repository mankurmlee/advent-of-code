class Puzzle
{
    public Dictionary<string, string> Orbits { get; set; }

    public Puzzle(string filename)
    {
        Dictionary<string, string> orbits = [];

        string[] lines = File.ReadAllLines(filename);

        foreach (string line in lines)
        {
            string[] data = line.Split(")");
            orbits[data[1]] = data[0];
        }

        Orbits = orbits;
    }

    public override string ToString()
    {
        var txt = string.Join(" ", Orbits);
        return $"[{txt}]";
    }

    public int CountAll()
    {
        int cnt = 0;
        foreach (var k in Orbits.Keys)
        {
            cnt += Count(k);
        }
        return cnt;
    }

    private int Count(string key)
    {
        int cnt = 0;
        while (Orbits.ContainsKey(key))
        {
            cnt++;
            key = Orbits[key];
        }
        return cnt;
    }

    public int TransferCost(string n1, string n2)
    {
        Dictionary<string, int> d = [];

        foreach (string s in GetAncestry(n1))
        {
            d.TryGetValue(s, out int value);
            d[s] = value + 1;
        }
        foreach (string s in GetAncestry(n2))
        {
            d.TryGetValue(s, out int value);
            d[s] = value + 1;
        }

        int cnt = 0;
        foreach (int v in d.Values)
        {
            if (v == 1)
            {
                cnt++;
            }
        }
        return cnt;
    }

    private List<string> GetAncestry(string node)
    {
        List<string> anc = [];
        while (Orbits.ContainsKey(node))
        {
            node = Orbits[node];
            anc.Add(node);
        }
        return anc;
    }
}

class Program
{
    static void Main(string[] args)
    {
        Puzzle p = new(args[0]);
        // Console.WriteLine(p);

        int cnt = p.CountAll();
        Console.WriteLine($"Part 1: {cnt}");

        int xfer = p.TransferCost("YOU", "SAN");
        Console.WriteLine($"Part 2: {xfer}");
    }
}
