struct Vec(int x, int y)
{
    public static readonly Vec[] DIR = [
        new(0, -1), new (1, 0), new(0, 1), new(-1, 0)
    ];
    public int x = x;
    public int y = y;
    public Vec(): this(0, 0) {}
    public readonly Vec Add(Vec other)
    {
        return new(x + other.x, y + other.y);
    }

    public override readonly bool Equals(object? obj)
    {
        if (obj == null || obj is not Vec) return false;
        Vec other = (Vec)obj;
        return x == other.x && y == other.y;
    }

    public override int GetHashCode()
    {
        return HashCode.Combine(x, y);
    }

    public static bool operator ==(Vec a, Vec b)
    {
        return a.Equals(b);
    }

    public static bool operator !=(Vec a, Vec b)
    {
        return !a.Equals(b);
    }

    public override readonly string ToString()
    {
        return $"({x}, {y})";
    }
}

readonly struct Edge(string id, int cost)
{
    public readonly string Id = id;
    public readonly int Cost = cost;
    public override string ToString()
    {
        return $"{{'{Id}' {Cost}}}";
    }
}

class Map
{
    public readonly char[][] Grid;
    public readonly int Width;
    public readonly int Height;
    public Map(string filename)
    {
        string[] lines = File.ReadAllLines(filename);
        int n = lines.Length;
        Grid = new char[n][];
        for (int y = 0; y < n; y++)
            Grid[y] = lines[y].ToCharArray();
        Width = lines[0].Length;
        Height = n;
    }

    public Graph BuildGraph(Vec start)
    {
        Graph g = new();
        HashSet<Vec> seen = [];
        Stack<Vec> q = [];
        q.Push(start);

        while (q.Count > 0)
        {
            Vec v = q.Pop();
            seen.Add(v);
            foreach (Vec dir in Vec.DIR)
            {
                var m = FindEdge(v, dir);
                if (!m.HasValue) continue;
                var (edge, pos) = m.Value;
                if (seen.Contains(pos)) continue;
                g.Add(GetTileId(v), edge);
                if (q.Contains(pos)) continue;
                q.Push(pos);
            }
        }
        return g;
    }

    private (Edge edge, Vec pos)? FindEdge(Vec initial, Vec dir)
    {
        int cost = 0;
        Vec last = new();
        Vec pos = initial;
        while (true)
        {
            // Check if we've hit a door or key
            char tile = Grid[pos.y][pos.x];
            if (cost != 0 && tile != '.')
                return (new(GetTileId(pos), cost), pos);

            List<Vec> exits = [];
            foreach (Vec v in Vec.DIR)
            {
                // On first cycle, ignore other directions
                if (cost == 0 && v != dir) continue;
                Vec adj = pos.Add(v);

                // Check we're not backtracking
                if (cost != 0 && adj == last) continue;

                // Check coordinates are in bounds
                if (adj.x < 0 || adj.y < 0 ||
                    adj.x >= Width || adj.y >= Height) continue;

                // Check we've not hit a wall
                if (Grid[adj.y][adj.x] == '#') continue;

                exits.Add(adj);
            }

            // check if we've reached a dead end
            if (exits.Count == 0)
                return null;

            // check if we've reached an intersection
            if (exits.Count > 1)
                return (new(GetTileId(pos), cost), pos);

            last = pos;
            pos = exits[0];
            cost++;
        }
    }

    private string GetTileId(Vec pos)
    {
        char tile = Grid[pos.y][pos.x];
        if (tile != '.')
            return tile.ToString();
        int num = (pos.y << 8) + pos.x;
        return num.ToString();
    }

    public void PartTwo()
    {
        Vec start = new();
        for (int y = 0; y < Height; y++)
            for (int x = 0; x < Width; x++)
                if (Grid[y][x] == '@')
                    start = new(x, y);

        if (start.Equals(new())) return;

        Dictionary<Vec, char> replacements = new()
        {
            { new(-1, -1), '@' }, { new(0, -1), '#' }, { new(1, -1), '@' },
            { new(-1,  0), '#' }, { new(0,  0), '#' }, { new(1,  0), '#' },
            { new(-1,  1), '@' }, { new(0,  1), '#' }, { new(1,  1), '@' },
        };
        foreach (var (o, c) in replacements)
        {
            Vec v = start.Add(o);
            if (v.x < 0 || v.y < 0 || v.x >= Width || v.y > Height) continue;
            Grid[v.y][v.x] = c;
        }
    }
}

class Path(string id, int cost, Path? last)
{
    public string Id = id;
    public int Cost = cost;
    public Path? Last = last;

    public HashSet<string> ToHashSet()
    {
        HashSet<string> ret = [];
        Path p = this;
        ret.Add(p.Id);
        while (p.Last != null)
        {
            p = p.Last;
            ret.Add(p.Id);
        }
        return ret;
    }
    public bool DoorLocked(IEnumerable<char> missing)
    {
        Path p = new("", 0, this);
        while (p.Last != null)
        {
            p = p.Last;
            char c = p.Id[0];
            // We only want to examine doors i.e. upper case letter Ids...
            if (c < 'A' || c > 'Z') continue;
            if (missing.Contains(char.ToLower(c))) return true;
        }
        return false;
    }
}

class Graph
{
    private readonly Dictionary<string, List<Edge>> Nodes = [];
    private readonly Dictionary<string, Path> PathMemo = [];

    public HashSet<char> Keys = [];
    public void Add(string id, Edge e)
    {
        if (!Nodes.TryGetValue(id, out List<Edge>? srcList))
        {
            srcList = [];
            Nodes[id] = srcList;
            char c = id[0];
            if (c >= 'a' && c <= 'z')
                Keys.Add(c);
        }
        srcList.Add(e);

        if (!Nodes.TryGetValue(e.Id, out List<Edge>? destList))
        {
            destList = [];
            Nodes[e.Id] = destList;
            char c = e.Id[0];
            if (c >= 'a' && c <= 'z')
                Keys.Add(c);
        }
        destList.Add(new(id, e.Cost));
    }
    public override string ToString()
    {
        string[] elems = new string[Nodes.Count];
        int i = 0;
        foreach (var kv in Nodes)
        {
            string list = string.Join(" ", kv.Value);
            elems[i++] = $"'{kv.Key}': [{list}]";
        }
        string bunch = string.Join(",\n\t", elems);
        return $"{{\n\t{bunch}\n}}";
    }
    public Path FindPath(string depart, string dest)
    {
        string cacheKey = depart + dest;
        if (PathMemo.TryGetValue(cacheKey, out Path? cached))
            return cached;

        Path best = new(depart, 0, null);
        PriorityQueue<Path, int> queue = new();
        queue.Enqueue(best, best.Cost);

        while (queue.Count > 0)
        {
            Path p = queue.Dequeue();
            if (p.Id == dest)
            {
                best = p;
                break;
            }
            HashSet<string> visited = p.ToHashSet();
            foreach (Edge e in Nodes[p.Id])
            {
                if (visited.Contains(e.Id)) continue;
                Path next = new(e.Id, p.Cost + e.Cost, p);
                queue.Enqueue(next, next.Cost);
            }
        }

        PathMemo[cacheKey] = best;
        return best;
    }
}

class Robot
{
    public List<char> LostKeys;
    public char Pos;
    public readonly Graph MyGraph;

    public Robot(IEnumerable<char> lost, char pos, Graph graph)
    {
        LostKeys = lost.ToList();
        LostKeys.Sort();
        Pos = pos;
        MyGraph = graph;
    }
    public Robot Clone()
    {
        List<char> l = new(LostKeys);
        return new(l, Pos, MyGraph);
    }
    public override bool Equals(object? obj)
    {
        if (obj == null || obj is not Robot) return false;
        Robot other = (Robot)obj;
        if (LostKeys.Count != other.LostKeys.Count) return false;
        if (LostKeys.Count == 0) return true;
        if (Pos != other.Pos) return false;
        for (int i = 0; i < LostKeys.Count; i++)
            if (LostKeys[i] != other.LostKeys[i])
                return false;
        return true;
    }
    public override int GetHashCode()
    {
        char p = Pos;
        if (LostKeys.Count == 0) p = '_';
        string s = new([p, .. LostKeys]);
        return s.GetHashCode();
    }
    public override string ToString()
    {
        if (LostKeys.Count == 0)
            return "_";
        string l = string.Join(" ", LostKeys);
        string p = Pos.ToString();
        return $"{{{p} [{l}]}}";
    }
}

class Save(List<Robot> robots)
{
    public List<Robot> Robots = robots;

    public Save Clone()
    {
        return new(new(Robots));
    }
    public override bool Equals(object? obj)
    {
        if (obj == null || obj is not Save) return false;
        Save other = (Save)obj;
        for (int i = 0; i < Robots.Count; i++)
            if (!Robots[i].Equals(other.Robots[i]))
                return false;
        return true;
    }
    public override int GetHashCode()
    {
        unchecked {
            int i = 17;
            foreach (Robot r in Robots)
                i += 23 * r.GetHashCode();
            return i;
        }
    }
    public override string ToString()
    {
        string r = string.Join(" ", Robots);
        return $"[{r}]";
    }
    public int KeyCount()
    {
        int c = 0;
        foreach (Robot r in Robots)
            c += r.LostKeys.Count;
        return c;
    }
    public IEnumerable<char> LostKeys()
    {
        List<char> l = [];
        foreach (Robot r in Robots)
            l.AddRange(r.LostKeys);
        return l;
    }
}

class Solver
{
    public List<Graph> Graphs = [];
    private Dictionary<Save,int> FewestMemo = [];

    public int FewestSteps(Map m)
    {
        FewestMemo = [];
        BuildGraphs(m);
        return GetFewest(StartState());
    }

    private Save StartState()
    {
        List<Robot> robots = [];
        foreach (Graph g in Graphs)
            robots.Add(new(g.Keys, '@', g));

        return new(robots);
    }

    private int GetFewest(Save s)
    {
        if (FewestMemo.TryGetValue(s, out int cached))
            return cached;

        int best = int.MaxValue;
        for (int i = 0; i < s.Robots.Count; i++)
        {
            Robot r = s.Robots[i];
            foreach (char c in r.LostKeys)
            {
                string depart = r.Pos.ToString();
                string dest = c.ToString();

                Path p = r.MyGraph.FindPath(depart, dest);
                if (p.DoorLocked(s.LostKeys())) continue;
                int cost = p.Cost;
                if (cost == 0) continue;

                // Recurse
                if (s.KeyCount() > 1)
                {
                    Robot childRobot = r.Clone();
                    childRobot.LostKeys.Remove(c);
                    childRobot.Pos = c;
                    Save childSave = s.Clone();
                    childSave.Robots[i] = childRobot;
                    int childCost = GetFewest(childSave);
                    cost += childCost;
                }

                // Keep cheapest
                if (cost < best)
                    best = cost;
            }
        }

        FewestMemo[s] = best;
        return best;
    }

    private void BuildGraphs(Map m)
    {
        Graphs = [];
        for (int y = 0; y < m.Height; y++)
            for (int x = 0; x < m.Width; x++)
                if (m.Grid[y][x] == '@')
                    Graphs.Add(m.BuildGraph(new(x, y)));
    }

    public override string ToString()
    {
        return string.Join(", ", Graphs);
    }
}

class Program
{
    static void Main(string[] args)
    {
        string filename = args[0];
        // string filename = "D:\\aoc\\2019\\18\\sample5.txt";
        Map m = new(filename);
        Solver s = new();

        int cost = s.FewestSteps(m);
        Console.WriteLine($"Part 1: {cost}");

        m.PartTwo();
        cost = s.FewestSteps(m);
        Console.WriteLine($"Part 2: {cost}");
    }
}
