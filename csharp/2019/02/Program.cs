class Puzzle
{
    public List<int> Data { get; }

    public Puzzle(string filename)
    {
        string[] lines = File.ReadAllLines(filename);
        string[] nums = lines[0].Split(",");

        List<int> data = [];

        foreach (string num in nums)
        {
            data.Add(int.Parse(num));
        }
        Data = data;
    }

    public Puzzle(List<int> data)
    {
        Data = data;
    }

    public Puzzle Clone()
    {
        List<int> data = new(Data);
        return new Puzzle(data);
    }

    public override string ToString()
    {
        return "[" + string.Join(" ", Data) + "]";
    }

    public int Execute(int noun, int verb)
    {
        Data[1] = noun;
        Data[2] = verb;
        for (int ip = 0; Data[ip] != 99; ip += 4)
        {
            int ptr1 = Data[ip+1];
            int ptr2 = Data[ip+2];
            int ptr3 = Data[ip+3];
            switch (Data[ip])
            {
                case 1:
                    Data[ptr3] = Data[ptr1] + Data[ptr2];
                    break;
                case 2:
                    Data[ptr3] = Data[ptr1] * Data[ptr2];
                    break;
            }
        }
        return Data[0];
    }

    public int FindInputs(int result)
    {
        for (int v = 0; v < 100; v++)
        {
            for (int n = 0; n < 100; n++)
            {
                Puzzle c = Clone();
                if (c.Execute(n, v) == result)
                {
                    return 100 * n + v;
                }
            }
        }
        return 0;
    }
}

class Program
{
    static void Main(string[] args)
    {
        string filename = args[0];
        Puzzle p = new(filename);
        Puzzle c = p.Clone();
        Console.WriteLine(p);

        c.Execute(12, 2);
        Console.WriteLine("Part 1: " + c.Data[0]);

        int r = p.FindInputs(19690720);
        Console.WriteLine($"Part 2: {r}");
    }
}