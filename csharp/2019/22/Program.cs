abstract record Method {
    public record Reverse : Method;
    public record Cut(long N) : Method;
    public record Riffle(long N) : Method;
}

class Deck
{
    public readonly int Size;
    public int[] Cards;
    public int[] Swap;
    public Method[] Shuffles;

    public Deck(string filename)
    {
        string[] lines = File.ReadAllLines(filename);
        Size = int.Parse(lines[0]);

        Cards = new int[Size];
        Swap = new int[Size];

        for (int i = 0; i < Size; i++)
            Cards[i] = i;

        Shuffles = new Method[lines.Length - 1];
        for (int i = 1; i < lines.Length; i++)
        {
            string[] fields = lines[i].Split(' ');
            string last = fields.Last();

            if (last == "stack")
            {
                Shuffles[i - 1] = new Method.Reverse();
                continue;
            }

            int n = int.Parse(last);
            if (fields.Length == 2)
                Shuffles[i - 1] = new Method.Cut(n);
            else
                Shuffles[i - 1] = new Method.Riffle(n);
        }
    }

    private void Riffle(int n)
    {
        SwapDecks();
        int i = 0;
        foreach (int v in Swap)
        {
            Cards[i] = v;
            i = (i + n) % Size;
        }
    }

    private void Cut(int n)
    {
        SwapDecks();
        if (n < 0)
            n += Size;
        int m = Size - n;
        Array.Copy(Swap, n, Cards, 0, m);
        Array.Copy(Swap, 0, Cards, m, n);
    }

    private void Reverse()
    {
        SwapDecks();
        int max = Size - 1;
        for (int i = 0; i < Size; i++)
            Cards[max - i] = Swap[i];
    }

    private void SwapDecks()
    {
        (Cards, Swap) = (Swap, Cards);
    }

    public override string ToString()
    {
        if (Cards.Length > 25)
            return string.Join(" ", Cards[..25]);
        return string.Join(" ", Cards);
    }

    public void Shuffle()
    {
        foreach (Method s in Shuffles)
        {
            switch (s)
            {
                case Method.Reverse:
                    Reverse();
                    break;
                case Method.Cut arg:
                    Cut((int)arg.N);
                    break;
                case Method.Riffle arg:
                    Riffle((int)arg.N);
                    break;
            }
        }
    }

    public VDeck CreateVDeck()
    {
        return new VDeck(Size, Shuffles);
    }
}

class VDeck
{
    public readonly long Size;
    public readonly Method[] Steps;
    public readonly Dictionary<long, long> HCache = [];
    public readonly Dictionary<long, long> VCache = [];

    public VDeck(long size, Method[] steps)
    {
        Size = size;
        Steps = steps;
        HCache[0] = 1;
        VCache[0] = 0;
        HCache[1] = FirstCardIncrement();
        VCache[1] = BackTrace(0);
    }

    public VDeck Resize(long size)
    {
        return new VDeck(size, Steps);
    }

    public long Trace(long initPos)
    {
        long pos = initPos;
        foreach (Method s in Steps)
        {
            switch (s)
            {
                case Method.Reverse:
                    pos = Size - 1 - pos;
                    break;
                case Method.Cut arg:
                    pos -= arg.N;
                    if (pos < 0)
                        pos += Size;
                    if (pos >= Size)
                        pos %= Size;
                    break;
                case Method.Riffle arg:
                    pos = pos * arg.N % Size;
                    break;
            }
        }
        return pos;
    }
    public long BackTrace(long initPos)
    {
        long pos = initPos;
        for (int i = Steps.Length - 1; i >= 0; i--)
        {
            switch (Steps[i])
            {
                case Method.Reverse:
                    pos = Size - 1 - pos;
                    break;
                case Method.Cut arg:
                    pos += arg.N;
                    if (pos < 0)
                        pos += Size;
                    if (pos >= Size)
                        pos %= Size;
                    break;
                case Method.Riffle arg:
                    pos = BackRiffle(arg.N, pos);
                    break;
            }
        }
        return pos;
    }

    private long BackRiffle(long n, long pos)
    {
        for (int i = 0; i < n; i++)
            if ((Size * i + pos) % n == 0)
                return (Size * i + pos) / n;

        // Shouldn't get here
        throw new NotImplementedException();
    }

    private long FirstCardIncrement()
    {
        long pos = Trace(0);
        pos = (pos + 1) % Size;
        return BackTrace(pos);
    }

    // HSI: Horizontal Skip Increment
    // To predict the next card, add the HSI and MOD by pack size
    // HSI differs from shuffle to shuffle
    // But squaring the hsi (and mod pack size) will
    // predict the hsi for current shuffle x 2
    public long CardIncrement(long i)
    {
        if (HCache.TryGetValue(i, out long cached))
            return cached;

        long hsi = 0;
        if ((i & (i - 1)) == 0)
        {
            // Powers of 2 can be recursively generated from powers of 2
            long half = CardIncrement(i >> 1);
            hsi = ModMul(half, half, Size);
        }
        else
        {
            long powtwo = 1;
            while ((powtwo & i) == 0)
                powtwo <<= 1;
            hsi = ModMul(CardIncrement(powtwo), CardIncrement(i - powtwo), Size);
        }

        Console.WriteLine($"HSI {i}: {hsi}");
        HCache[i] = hsi;
        return hsi;
    }

    // VSM: Vertical skip multiplier
    // Expects power of 2 param
    public long TopCard(long i)
    {
        if (VCache.TryGetValue(i, out long cached))
            return cached;

        long vsm = 0;
        if ((i & (i - 1)) == 0)
        {
            long half = i >> 1;
            vsm = ModMul(TopCard(half), CardIncrement(half) + 1, Size);
        }
        else
        {
            long powtwo = 1;
            while ((powtwo & i) == 0)
                powtwo <<= 1;
            long rem = i - powtwo;
            vsm = ModMul(TopCard(powtwo), CardIncrement(rem), Size);
            vsm = (vsm + TopCard(rem)) % Size;
        }

        Console.WriteLine($"VSM {i}: {vsm}");
        VCache[i] = vsm;
        return vsm;
    }

    public static long ModMul(long a, long b, long mod)
    {
        a %= mod;
        b %= mod;
        long result = 0;
        while (b > 0)
        {
            if ((b & 1) == 1)
                result = (result + a) % mod;
            a = (a << 1) % mod;
            b >>= 1;
        }
        return result;
    }

    public long CardXAfterYShuffles(long x, long y)
    {
        long topcard = TopCard(y);
        long hsi = ModMul(CardIncrement(y), x, Size);
        return (topcard + hsi) % Size;
    }
}


class Program
{
    static void Main(string[] args)
    {
        string filename = args[0];
        Deck d = new(filename);
        VDeck v = d.CreateVDeck();

        PartOne(d);
        // DebugShuffles(d);
        PartTwo(v);
    }

    private static void DebugShuffles(Deck d)
    {
        long n = d.Size;
        if (n > 20000) n = 20000;

        Console.WriteLine($"0: {d}");
        for (int i = 1; i <= n; i++)
        {
            d.Shuffle();
            Console.WriteLine($"{i}: {d}");
        }
    }

    private static void PartOne(Deck d)
    {
        if (d.Size < 2020) return;
        d.Shuffle();
        long n = Array.IndexOf(d.Cards, 2019);
        Console.WriteLine($"Part 1: {n}");
    }

    private static void PartTwo(VDeck v)
    {
        long cards = 119315717514047;
        long shufs = 101741582076661;
        v = v.Resize(cards);

        long card = v.CardXAfterYShuffles(2020, shufs);
        Console.WriteLine($"Part 2: {card}");
    }
}
