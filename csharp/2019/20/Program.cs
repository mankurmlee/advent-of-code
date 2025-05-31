class Map
{
    private readonly int Width;
    private readonly int Height;
    private readonly char[][] Grid;

    private readonly Dictionary<Vec, Vec> Portals = [];
    private readonly HashSet<Vec> Seen = [];

    public Map(string filename)
    {
        string[] lines = File.ReadAllLines(filename);
        Height = lines.Length;
        for (int i = 0; i < Height; i++)
            if (Width < lines[i].Length)
                Width = lines[i].Length;

        Grid = new char[Height][];
        for (int i = 0; i < Height; i++)
        {
            Grid[i] = new char[Width];
            lines[i].CopyTo(0, Grid[i], 0, lines[i].Length);
        }
    }

    public override string ToString()
    {
        List<string> c = [];
        for (int i = 0; i < Height; i++)
            c.Add(new(Grid[i]));
        return string.Join(Environment.NewLine, c);
    }

    private Dictionary<string, List<(Vec, Vec)>> PortalsByName()
    {
        Dictionary<string, List<(Vec, Vec)>> res = [];

        for (int y = 0; y < Height; y++)
            for (int x = 0; x < Width; x++)
            {
                if (Grid[y][x] != '.') continue;
                Vec v = new(x, y);
                foreach (var (delta, rev) in new List<(Vec, bool)>() {
                    (new(0, -1), true), (new(1, 0), false),
                    (new(0, 1), false), (new(-1, 0), true),
                })
                {
                    string pName = GetPortal(v, delta, rev);
                    if (pName == "") continue;
                    if (!res.TryGetValue(pName, out List<(Vec, Vec)>? conns))
                    {
                        conns = [];
                        res[pName] = conns;
                    }
                    Vec adj = v.Add(delta);
                    conns.Add((v, adj));

                    // string s = string.Join(", ", conns);
                    // Console.WriteLine($"{pName}: {s}");
                    break;
                }
            }

        return res;
    }

    public Graph ToGraph()
    {
        var portsByName = PortalsByName();
        Portals.Clear();

        foreach (var exits in portsByName.Values)
        {
            if (exits.Count != 2) continue;
            (Vec posA, Vec vposA) = exits[0];
            (Vec posB, Vec vposB) = exits[1];
            Portals[vposA] = posB;
            Portals[vposB] = posA;

            // Console.WriteLine($"{vposA} => {posB}");
            // Console.WriteLine($"{vposB} => {posA}");
        }

        Vec start = portsByName["AA"][0].Item1;
        Vec end = portsByName["ZZ"][0].Item1;

        Console.WriteLine($"Start: {start}");
        Console.WriteLine($"End: {end}");

        var nodes = FindNodes(start, end);

        return new Graph(nodes, start, end);
    }

    private Dictionary<Vec, List<Edge>> FindNodes(Vec start, Vec end)
    {
        Dictionary<Vec, List<Edge>> nodes = [];

        Seen.Clear();
        Seen.Add(start);
        Seen.Add(end);

        Stack<Vec> q = new();
        q.Push(start);
        q.Push(end);

        while (q.Count > 0)
        {
            Vec node = q.Pop();
            List<Edge> edges = FindEdges(node);
            if (edges.Count == 0) continue;
            nodes[node] = edges;

            string mkl = string.Join(", ", edges);
            Console.WriteLine($"{node} : {mkl}");

            foreach (var e in edges)
            {
                if (Seen.Contains(e.Adj)) continue;
                Seen.Add(e.Adj);
                q.Push(e.Adj);
            }
        }
        return nodes;
    }

    private List<Edge> FindEdges(Vec start)
    {
        Dictionary<Vec, Edge> best = [];
        foreach (Vec dir in Vec.DIR)
        {
            Edge? e = FindEdge(start, dir);
            if (!e.HasValue) continue;
            if (best.TryGetValue(e.Value.Adj, out Edge b))
                if (b.Cost <= e.Value.Cost) continue;

            best[e.Value.Adj] = e.Value;
        }
        return [.. best.Values];
    }

    private Edge? FindEdge(Vec start, Vec dir)
    {
        Vec pos = start;
        Vec delta = dir;
        int floorChange = 0;
        int cost = 0;
        Vec last;
        int options;
        do
        {
            last = pos;
            pos.AddEquals(delta);
            cost++;

            // Traverse any portals first
            if (Portals.TryGetValue(pos, out Vec dest))
            {
                if (pos.x == 1 ||
                    pos.y == 1 ||
                    pos.x == Width - 2 ||
                    pos.y == Height - 2)
                    floorChange -= 1;
                else
                    floorChange += 1;

                pos = dest;
                break;
            }

            // If we're not in an open passage then bail out
            if (Grid[pos.y][pos.x] != '.')
                return null;

            // Check if we've hit a known node
            if (Seen.Contains(pos)) break;

            // Check if it's a junction type node
            options = 0;
            foreach (Vec d in Vec.DIR)
            {
                Vec n = pos.Add(d);
                if (Portals.ContainsKey(n))
                    return new(pos, floorChange, cost);
                if (n == last) continue;
                if (Grid[n.y][n.x] != '.') continue;
                options++;
                delta = d;
            }
            if (options == 0) return null;
        } while(options == 1);
        return new(pos, floorChange, cost);
    }

    private string GetPortal(Vec start, Vec delta, bool reverse)
    {
        char[] res = new char[2];
        Vec v = start;
        v.AddEquals(delta);
        char c = Grid[v.y][v.x];
        if (c < 'A' || c > 'Z') return "";
        res[0] = c;
        v.AddEquals(delta);
        c = Grid[v.y][v.x];
        if (c < 'A' || c > 'Z') return "";
        res[1] = c;
        if (reverse)
            Array.Reverse(res);
        return new string(res);
    }
}

struct Edge(Vec adj, int floorChange, int cost)
{
    public Vec Adj = adj;
    public int FloorChange = floorChange;
    public int Cost = cost;
    public override readonly string ToString()
    {
        return $"<{Adj}, {FloorChange}, {Cost}>";
    }
}

struct Path(Vec pos, int level, int cost)
{
    public Vec Pos = pos;
    public int Level = level;
    public int Cost = cost;

    public readonly Path Move(Edge e)
    {
        return new(e.Adj, Level + e.FloorChange, Cost + e.Cost);
    }
}

class Graph(Dictionary<Vec, List<Edge>> nodes, Vec start, Vec end)
{
    public readonly Dictionary<Vec, List<Edge>> Nodes = nodes;
    public readonly Vec Start = start;
    public readonly Vec End = end;

    public Path Solve(bool recurse)
    {
        Dictionary<(Vec, int), int> best = [];
        best.Add((Start, 0), 0);

        Path start = new(Start, 0, 0);
        PriorityQueue<Path, int> q = new();
        q.Enqueue(start, start.Cost);

        while (q.Count > 0)
        {
            Path p = q.Dequeue();
            if (p.Pos == End && (p.Level == 0 || !recurse))
                return p;
            foreach (Edge e in Nodes[p.Pos])
            {
                Path next = p.Move(e);
                if (recurse && next.Level < 0) continue;
                if (best.TryGetValue((next.Pos, next.Level), out int cost))
                    if (next.Cost >= cost) continue;
                best[(next.Pos, next.Level)] = next.Cost;
                q.Enqueue(next, next.Cost);
            }
        }

        return start;
    }
}

class Program
{
    static void Main(string[] args)
    {
        string filename = args[0];
        Map m = new(filename);
        Console.WriteLine(m);

        Graph g = m.ToGraph();
        Path p = g.Solve(false);
        Console.WriteLine($"Part 1: {p.Cost}");

        p = g.Solve(true);
        Console.WriteLine($"Part 2: {p.Cost}");
    }
}
