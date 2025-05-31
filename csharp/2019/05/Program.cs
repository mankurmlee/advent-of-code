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

    public int Execute(int input)
    {
        int output = 0;
        int ip = 0;
        while (Data[ip] != 99)
        {
            string opcode = Data[ip].ToString().PadLeft(4, '0');
            switch (opcode[3])
            {
                case '1':
                    Add(opcode, ip);
                    ip += 4;
                    break;
                case '2':
                    Mul(opcode, ip);
                    ip += 4;
                    break;
                case '3':
                    Store(opcode, ip, input);
                    ip += 2;
                    break;
                case '4':
                    output = Load(opcode, ip);
                    ip += 2;
                    break;
                case '5':
                    ip = JumpIfTrue(opcode, ip);
                    break;
                case '6':
                    ip = JumpIfFalse(opcode, ip);
                    break;
                case '7':
                    LessThan(opcode, ip);
                    ip += 4;
                    break;
                case '8':
                    EqualsTo(opcode, ip);
                    ip += 4;
                    break;
            }
        }
        return output;
    }

    private void EqualsTo(string opcode, int ip)
    {
        int arg1 = GetArg(ip+1, opcode[1]);
        int arg2 = GetArg(ip+2, opcode[0]);
        int ptr = Data[ip+3];
        Data[ptr] = arg1 == arg2 ? 1 : 0;
    }

    private void LessThan(string opcode, int ip)
    {
        int arg1 = GetArg(ip+1, opcode[1]);
        int arg2 = GetArg(ip+2, opcode[0]);
        int ptr = Data[ip+3];
        Data[ptr] = arg1 < arg2 ? 1 : 0;
    }

    private int JumpIfTrue(string opcode, int ip)
    {
        int arg1 = GetArg(ip+1, opcode[1]);
        int arg2 = GetArg(ip+2, opcode[0]);
        if (arg1 != 0)
        {
            return arg2;
        }
        return ip + 3;
    }

    private int JumpIfFalse(string opcode, int ip)
    {
        int arg1 = GetArg(ip+1, opcode[1]);
        int arg2 = GetArg(ip+2, opcode[0]);
        if (arg1 == 0)
        {
            return arg2;
        }
        return ip + 3;
    }

    private void Add(string opcode, int ip)
    {
        int arg1 = GetArg(ip+1, opcode[1]);
        int arg2 = GetArg(ip+2, opcode[0]);
        int ptr = Data[ip+3];
        Data[ptr] = arg1 + arg2;
    }

    private void Mul(string opcode, int ip)
    {
        int arg1 = GetArg(ip+1, opcode[1]);
        int arg2 = GetArg(ip+2, opcode[0]);
        int ptr = Data[ip+3];
        Data[ptr] = arg1 * arg2;
    }

    private void Store(string opcode, int ip, int input)
    {
        int ptr = Data[ip+1];
        Data[ptr] = input;
    }

    private int Load(string opcode, int ip)
    {
        int arg1 = GetArg(ip+1, opcode[1]);
        return arg1;
    }

    private int GetArg(int v, char mode)
    {
        if (mode == '1')
        {
            return Data[v];
        }
        int ptr = Data[v];
        return Data[ptr];
    }
}

class Program
{
    static void Main(string[] args)
    {
        Puzzle c = new(args[0]);

        Puzzle p = c.Clone();
        int output = p.Execute(1);
        Console.WriteLine($"Part 1: {output}");

        p = c.Clone();
        output = p.Execute(5);
        Console.WriteLine($"Part 2: {output}");
    }
}