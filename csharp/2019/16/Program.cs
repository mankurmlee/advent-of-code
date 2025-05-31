
class Signal(int[] data)
{
    public readonly int[] Data = data;

    public Signal(string filename): this(Array.Empty<int>())
    {
        string[] lines = File.ReadAllLines(filename);
        List<string> data = new(lines);
        int n = data[0].Length;
        Data = new int[n];
        char[] ch = data[0].ToCharArray();
        for (int i = 0; i < n; i++)
            Data[i] = int.Parse(ch[i].ToString());
    }

    public override string ToString()
    {
        return string.Join("", Data)[..8];
    }

    public Signal FFT()
    {
        int n = Data.Length;
        int[] output = new int[n];
        for (int j = 0; j < n; j++)
        {
            Pattern patt = new(j);
            for (int i = 0; i < n; i++)
                output[j] += Data[i] * patt.Next();

            if (output[j] < 0)
                output[j] = -output[j];

            output[j] %= 10;
        }
        return new Signal(output);
    }

    public int PartTwo()
    {
        int len = Data.Length;
        int vidx = GetIndex();
        int vlen = len * 10000;

        int size = vlen - vidx;
        if (size <= 10) return 0;

        // Create the swap buffers
        int[] inBuf = new int[size];
        int[] outBuf = new int[size];

        // Fill in the buffer
        for (int i = vidx; i < vlen; i++)
            outBuf[i - vidx] = Data[i % len];

        // Do processing in reverse as it's easier
        Array.Reverse(outBuf);
        for (int i = 0; i < 100; i++)
        {
            (inBuf, outBuf) = (outBuf, inBuf);
            SumMod(inBuf, outBuf);
        }

        // Revert order
        Array.Reverse(outBuf);
        string txt = string.Join("", outBuf[..8]);
        return int.Parse(txt);
    }

    private static void SumMod(int[] inBuf, int[] outBuf)
    {
        int len = outBuf.Length;
        int sum = 0;
        for (int i = 0; i < len; i++)
        {
            sum = (sum + inBuf[i]) % 10;
            outBuf[i] = sum;
        }
    }

    private int GetIndex()
    {
        string txt = string.Join("", Data[..7]);
        return int.Parse(txt);
    }
}

class Pattern(int repeat)
{
    private int counter = 0;
    private readonly static int[] patt = [0, 1, 0, -1];
    private readonly int times = (repeat + 1);
    private readonly int mod = 4 * (repeat + 1);
    public int Next()
    {
        counter++;
        counter %= mod;
        return patt[counter / times];
    }
}

class Program
{
    static void Main(string[] args)
    {
        string filename = args[0];
        Signal s = new(filename);
        Console.WriteLine(s);

        // for (int i = 0; i < 100; i++)
        //     s = s.FFT();
        // Console.WriteLine($"Part 1: {s}");

        s = new(filename);
        int res = s.PartTwo();
        Console.WriteLine($"Part 2: {res}");
    }
}
