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
    public void AddEquals(Vec other)
    {
        x += other.x;
        y += other.y;
    }
    public override readonly bool Equals(object? obj)
    {
        if (obj == null || obj is not Vec) return false;
        Vec other = (Vec)obj;
        return x == other.x && y == other.y;
    }

    public override readonly int GetHashCode()
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
