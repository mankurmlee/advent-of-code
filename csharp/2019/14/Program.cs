class Book
{
    private readonly Dictionary<string, Recipe> recipes = [];
    private readonly Dictionary<Amount, Wallet> memo = [];

    public Book(string filename)
    {
        string[] lines = File.ReadAllLines(filename);
        foreach (string line in lines)
        {
            string[] data = line.Replace("=>", "")
                .Replace(",", "")
                .Split(' ', StringSplitOptions.RemoveEmptyEntries);
            int n = data.Length;

            Wallet input = new();
            int num = 0;
            for (int i = 0; i < n - 2; i += 2)
                input.Incr(data[i+1], int.Parse(data[i]));

            num = int.Parse(data[n-2]);
            string type = data[n-1];
            Amount output = new(num, type);

            recipes.Add(type, new Recipe(input, output));
        }
    }

    public Wallet GetCost(Amount amount)
    {
        Wallet w = new();
        if (amount.Type == "ORE")
        {
            // The cost of ORE is itself. Log this expenditure into the wallet
            w.Incr("ORE", (int)amount.Num);
            return w;
        }

        memo.TryGetValue(amount, out Wallet? value);
        if (value != null) return value;

        Recipe r = recipes[amount.Type];

        int scale = (int)(amount.Num / r.provide.Num);
        int mod   = (int)(amount.Num % r.provide.Num);
        if (mod > 0)
        {
            scale++;
            // Put unneeded items back into the wallet for later
            w.Decr(amount.Type, (int)r.provide.Num - mod);
        }

        foreach (var (type, num) in r.require.units)
        {
            // See if we already have enough items in our wallet already
            int req = (int)w.Use(type, num * scale);
            if (req < 0) continue;

            // If we don't have enough items in our wallet, we'll need to spend some ORE
            Wallet child = GetCost(new Amount(req, type));
            w.Merge(child);
        }

        // If we have too many items in the wallet we can refund them to reduce ORE required
        DoRefunds(w);

        memo[amount] = w;
        return memo[amount];
    }

    public void DoRefunds(Wallet w)
    {
        Wallet adjustment;
        do {
            adjustment = new();
            foreach (var (type, num) in w.units)
            {
                if (type == "ORE") continue;

                Recipe r = recipes[type];
                long p = r.provide.Num;
                if (p + num > 0) continue;

                long scale = (-num) / p;

                adjustment.Incr(type, p * scale);
                foreach (var (t, n) in r.require.units)
                    adjustment.Decr(t, n * scale);
            }
            w.Merge(adjustment);
        } while (adjustment.units.Count > 0);
    }

    public int Spend(long target)
    {
        long cost = GetFuelCost(1);
        int lo = (int)(target / cost);
        int hi = lo;

        while (cost <= target)
        {
            lo = hi;
            hi <<= 1;
            cost = GetFuelCost(hi);
        }

        while (hi - lo > 1)
        {
            int mid = (hi + lo) >> 1;
            cost = GetFuelCost(mid);

            if (cost > target)
                hi = mid;
            else
                lo = mid;
        }

        return lo;
    }

    private long GetFuelCost(int n)
    {
        Wallet w = GetCost(new Amount(1, "FUEL"));
        Wallet bw = new(w, n);
        DoRefunds(bw);
        return bw.Get("ORE");
    }
}

struct Amount(long n, string type)
{
    public long Num = n;
    public string Type = type;

    public override readonly string ToString()
    {
        return $"({Num}, {Type})";
    }
}

struct Recipe(Wallet i, Amount o)
{
    public Wallet require = i;
    public Amount provide = o;
}

class Wallet()
{
    public readonly Dictionary<string, long> units = [];

    public Wallet(Wallet other, int scale): this()
    {
        foreach (var (type, num) in other.units)
            units[type] = num * scale;
    }

    public override string ToString()
    {
        string txt = string.Join(" ", units);
        return $"{{{txt}}}";
    }

    public void Incr(string type, long num)
    {
        units.TryGetValue(type, out long value);
        value += num;
        if (value == 0)
            units.Remove(type);
        else
            units[type] = value;
    }

    public void Decr(string type, long num)
    {
        units.TryGetValue(type, out long value);
        value -= num;
        if (value == 0)
            units.Remove(type);
        else
            units[type] = value;
    }

    public long Get(string type)
    {
        units.TryGetValue(type, out long value);
        return value;
    }

    public long Use(string type, long num)
    {
        units.TryGetValue(type, out long value);
        long rem = value + num;
        if (rem >= 0)
            units.Remove(type);
        else
            units[type] = rem;
        return rem;
    }

    public void Merge(Wallet child)
    {
        foreach (var (type, num) in child.units)
            Incr(type, num);
    }
}

class Program
{
    static void Main(string[] args)
    {
        // string filename = "D:\\aoc\\2019\\14\\puzzle.txt";
        string filename = args[0];
        Book b = new(filename);
        Wallet w = b.GetCost(new Amount(1, "FUEL"));
        Console.WriteLine($"{w}");
        Console.WriteLine($"Part 1: {w.Get("ORE")}");

        int fuel = b.Spend(1000000000000);
        Console.WriteLine($"Part 2: {fuel}");
    }
}
