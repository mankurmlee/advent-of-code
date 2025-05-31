class Map()
{
    public Vec offset = new();
    public List<Vec> Asteroids = [];

    public Map(string filename): this()
    {
        string[] lines = File.ReadAllLines(filename);
        List<string> data = new(lines);

        int w = data[0].Length;
        int h = data.Count;

        for (int y = 0; y < h; y++)
        {
            string row = data[y];
            for (int x = 0; x < w; x++)
                if (row[x] == '#')
                    Asteroids.Add(new(x, y));
        }
    }

    public Map Recentre(Map m, Vec newOrigin)
    {
        Map newMap = new();

        foreach (Vec other in Asteroids)
        {
            if (newOrigin.Equals(other)) continue;
            Vec delta = newOrigin.Delta(other);
            newMap.Asteroids.Add(delta);
        }
        newMap.offset = newOrigin;

        return newMap;
    }

    public RadialMap FindMostSuitable()
    {
        RadialMap best = new();
        foreach (Vec ast in Asteroids)
        {
            RadialMap r = new(Recentre(this, ast));
            if (r.trajs.Count > best.trajs.Count)
                best = r;
        }
        return best;
    }
}

struct Vec(int x, int y)
{
    private static readonly int[] primes = [2, 3, 5, 7, 11, 13, 17, 19, 23];
    private const double radToDeg = 180.0 / Math.PI;
    public int X = x;
    public int Y = y;

    public Vec(): this(0, 0) {}

    public override readonly bool Equals(object? obj)
    {
        if (obj == null || GetType() != obj.GetType())
            return false;
        Vec other = (Vec)obj;
        return (X == other.X && Y == other.Y);
    }

    public override readonly int GetHashCode()
    {
        return HashCode.Combine(X, Y);
    }

    public readonly Vec Add(Vec other)
    {
        return new(X + other.X, Y + other.Y);
    }

    public readonly Vec Delta(Vec other)
    {
        int dx = other.X - X;
        int dy = other.Y - Y;
        return new(dx, dy);
    }

    public readonly double MagSquared()
    {
        double nx = X;
        double ny = Y;
        return nx * nx + ny * ny;
    }

    public readonly Vec Norm()
    {
        int absX = Math.Abs(X);
        int absY = Math.Abs(Y);
        int normX = X;
        int normY = Y;
        foreach (int prime in primes)
        {
            if (absX > 0 && absY > 0 && (prime > absX || prime > absY))
                break;

            while (normX % prime == 0 && normY % prime == 0)
            {
                normX /= prime;
                normY /= prime;
            }
        }
        return new(normX, normY);
    }

    public readonly double ToDeg()
    {
        double theta = Math.Atan2(Y, X) * radToDeg + 90;
        if (theta < 0) theta += 360;
        return theta;
    }
}

class RadialMap()
{
    public Vec myLoc = new();
    public Dictionary<Vec, List<Vec>> trajs = [];

    public RadialMap(Map m): this()
    {
        foreach (Vec ast in m.Asteroids)
        {
            Vec norm = ast.Norm();

            if (!trajs.ContainsKey(norm))
                trajs[norm] = [];

            trajs[norm].Add(ast);

            foreach (var kv in trajs)
                kv.Value.Sort((Vec a, Vec b) => {
                    return -a.MagSquared().CompareTo(b.MagSquared());
                });
        }
        myLoc = m.offset;
    }

    public Vec Vapourise(int n)
    {
        Vec last = new();

        List<Vec> toRemove = [];
        List<Vec> keys = [];
        int j = 0;
        for (int i = 0; i < n; i++)
        {
            if (j == keys.Count)
            {
                foreach (Vec k in toRemove)
                    trajs.Remove(k);
                toRemove = [];
                keys = [.. trajs.Keys];
                keys.Sort((Vec a, Vec b) => a.ToDeg().CompareTo(b.ToDeg()));
                j = 0;
                if (keys.Count == 0)
                {
                    Console.WriteLine("Out of asteroids!");
                    break;
                }
            }

            Vec key = keys[j++];
            int lastIdx = trajs[key].Count - 1;
            last = trajs[key][lastIdx];
            trajs[key].RemoveAt(lastIdx);

            if (trajs[key].Count == 0)
                toRemove.Add(key);
        }

        return myLoc.Add(last);
    }
}

class Program
{
    static void Main(string[] args)
    {
        string filename = args[0];
        Map m = new(filename);

        RadialMap station = m.FindMostSuitable();
        Console.WriteLine($"Part 1: {station.trajs.Count}");

        Vec last = station.Vapourise(200);
        int idx = last.X * 100 + last.Y;
        Console.WriteLine($"Part 2: {idx}");
    }
}
