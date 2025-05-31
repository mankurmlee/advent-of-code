class Droid(string filename)
{
    private readonly Intcode prog = new(filename);
    private readonly Dictionary<Vec, int> map = new()
    {
        { new (0, 0), 1 },
    };
    private Vec pos = new();
    private Vec tankPos = new();
    public void Move(int dir)
    {
        prog.Push(dir+1);
        Vec newPos = pos.Add(Vec.DIRVEC[dir]);
        int tile = (int)prog.Execute().Item2;
        if (!map.TryGetValue(newPos, out int value))
            map[newPos] = tile;
        if (tile != 0)
            pos = newPos;
    }

    public int FindOxygenSystem()
    {
        HashSet<Vec> options = [.. Vec.DIRVEC];
        while (options.Count > 0)
        {
            Node? path = NextMove(options);
            if (path == null) {
                Console.WriteLine("Can't FIND IT");
                break;
            }
            options.Remove(path.Pos);
            Console.WriteLine($"Investigating {path.Pos}");

            // Follow the path
            foreach (int dir in path.ToDir())
                Move(dir);

            if (map[pos] == 0) continue;
            if (map[pos] == 2)
                tankPos = pos;

            foreach (Vec dir in Vec.DIRVEC)
            {
                Vec o = pos.Add(dir);
                if (map.ContainsKey(o)) continue;
                options.Add(o);
            }

            string txt = string.Join(" ", options);
        }

        for (int y = 0; y < 50; y++)
        {
            char[] row = new char[50];
            for (int x = 0; x < 50; x++)
            {
                if (map.TryGetValue(new Vec(x-25, y-25), out int value))
                {
                    if (value == 0)
                        row[x] = '#';
                    else if (value == 1)
                        row[x] = ' ';
                    else
                        row[x] = 'O';
                }
                else
                    row[x] = '.';
            }
            string txt = new(row);
            Console.WriteLine(txt);
        }


        Node? p = FindPath(new(0, 0), tankPos);
        if (p == null) return 0;
        return p.Cost;
    }

    private Node? NextMove(HashSet<Vec> options)
    {
        // Sort options
        var optList = options.ToList();
        optList.Sort((Vec a, Vec b) => pos.Dist(a).CompareTo(pos.Dist(b)));

        Node? best = FindPath(pos, optList[0]);
        if (best == null) return null;

        foreach (Vec o in optList)
        {
            if (pos.Dist(o) > best.Cost) continue;

            Node? p = FindPath(pos, o);
            if (p == null) continue;
            if (p.Cost > best.Cost) continue;

            best = p;
        }
        return best;
    }

    private Node? FindPath(Vec start, Vec dest)
    {
        PriorityQueue<Node, int> queue = new();
        Node n = new(start, dest, 0, null);
        queue.Enqueue(n, -n.Estimate);

        while (queue.Count > 0)
        {
            Node c = queue.Dequeue();
            if (c.Pos.Equals(dest)) return c;
            HashSet<Vec> visited = c.ToSet();
            foreach (Vec dir in Vec.DIRVEC)
            {
                Vec adj = c.Pos.Add(dir);

                if (visited.Contains(adj)) continue;

                if (!adj.Equals(dest))
                {
                    if (!map.TryGetValue(adj, out int tile)) continue;
                    if (tile == 0) continue;
                }

                Node o = new(adj, dest, c.Cost + 1, c);
                queue.Enqueue(o, -o.Estimate);
            }
        }

        return null;
    }

    public int OxygenSaturationTime()
    {
        int timeMax = 0;
        HashSet<Vec> visited = [];
        Node n = new(tankPos, tankPos, 0, null);
        Queue<Node> queue = [];
        queue.Enqueue(n);
        while (queue.Count > 0)
        {
            Node c = queue.Dequeue();
            foreach (Vec dir in Vec.DIRVEC)
            {
                Vec adj = c.Pos.Add(dir);
                if (visited.Contains(adj)) continue;
                if (!map.TryGetValue(adj, out int tile)) continue;
                if (tile == 0) continue;

                int newCost = c.Cost + 1;
                if (newCost > timeMax)
                    timeMax = newCost;

                Node o = new(adj, adj, newCost, null);
                queue.Enqueue(o);
                visited.Add(adj);
            }
        }
        return timeMax;
    }
}

struct Vec(int x, int y)
{
    public static readonly Vec[] DIRVEC = [
        new(0, -1), new(0, 1), new(-1, 0), new(1, 0),
    ];
    public int X = x;
    public int Y = y;
    public Vec(): this(0, 0) {}
    public override readonly string ToString()
    {
        return $"({X}, {Y})";
    }

    public override readonly bool Equals(object? obj)
    {
        if (obj == null || obj is not Vec) return false;
        Vec other = (Vec)obj;
        return X == other.X && Y == other.Y;
    }

    public static bool operator ==(Vec a, Vec b)
    {
        return a.Equals(b);
    }

    public static bool operator !=(Vec a, Vec b)
    {
        return !a.Equals(b);
    }

    public readonly Vec Add(Vec other)
    {
        return new(X + other.X, Y + other.Y);
    }

    public readonly int Dist(Vec other)
    {
        return Math.Abs(other.X - X) + Math.Abs(other.Y - Y);
    }

    public override readonly int GetHashCode()
    {
        return HashCode.Combine(X, Y);
    }

    public readonly Vec Diff(Vec other)
    {
        return new(other.X - X, other.Y - Y);
    }
}

class Node(Vec p, Vec target, int c, Node? l)
{
    public Vec Pos = p;
    public int Cost = c;
    public int Estimate = c + p.Dist(target);
    public Node? Last = l;

    public List<int> ToDir()
    {
        List<int> directions = [];
        Node now = this;
        while (now.Last != null)
        {
            Node next = now.Last;
            Vec v = next.Pos.Diff(now.Pos);
            int dir = Array.IndexOf(Vec.DIRVEC, v);
            directions.Add(dir);
            now = next;
        }
        directions.Reverse();
        return directions;
    }

    public HashSet<Vec> ToSet()
    {
        HashSet<Vec> visited = [ Pos ];
        Node now = this;
        while (now.Last != null)
        {
            now = now.Last;
            visited.Add(now.Pos);
        }
        return visited;
    }
}

class Program
{
    static void Main(string[] args)
    {
        string filename = "D:\\aoc\\2019\\15\\puzzle.txt";
        Droid droid = new(filename);
        int dist = droid.FindOxygenSystem();
        Console.WriteLine($"Part 1: {dist}");

        int t = droid.OxygenSaturationTime();
        Console.WriteLine($"Part 2: {t}");
    }
}