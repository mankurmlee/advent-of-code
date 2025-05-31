using System.Reflection.Emit;

class Machine
{
    private int ip = 0;
    private int offset = 0;
    private long output = 0;
    private List<long> data;

    public Machine(string filename)
    {
        string[] lines = File.ReadAllLines(filename);
        string[] nums = lines[0].Split(",");

        List<long> data = [];

        foreach (string num in nums)
            data.Add(long.Parse(num));

        this.data = data;
    }

    public Machine(List<long> data)
    {
        this.data = data;
    }

    public Machine Clone()
    {
        List<long> data = new(this.data);
        return new Machine(data);
    }

    public long Execute(int input)
    {
        while (data[ip] != 99)
        {
            string opcode = data[ip].ToString().PadLeft(5, '0');

            switch (opcode[4])
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
                case '9':
                    AddBase(opcode, ip);
                    ip += 2;
                    break;
                default:
                    Console.WriteLine($"Unexpected opcode: {opcode[4]}");
                    break;
            }
        }
        return output;
    }

    private void AddBase(string opcode, int ip)
    {
        long arg1 = GetArg(opcode, 1);
        offset += (int)arg1;
    }

    private void EqualsTo(string opcode, int ip)
    {
        long arg1 = GetArg(opcode, 1);
        long arg2 = GetArg(opcode, 2);
        Set(opcode, 3, arg1 == arg2 ? 1 : 0);
    }

    private void LessThan(string opcode, int ip)
    {
        long arg1 = GetArg(opcode, 1);
        long arg2 = GetArg(opcode, 2);
        Set(opcode, 3, arg1 < arg2 ? 1 : 0);
    }

    private int JumpIfTrue(string opcode, int ip)
    {
        long arg1 = GetArg(opcode, 1);
        long arg2 = GetArg(opcode, 2);
        if (arg1 != 0)
            return (int)arg2;
        return ip + 3;
    }

    private int JumpIfFalse(string opcode, int ip)
    {
        long arg1 = GetArg(opcode, 1);
        long arg2 = GetArg(opcode, 2);
        if (arg1 == 0)
            return (int)arg2;
        return ip + 3;
    }

    private void Add(string opcode, int ip)
    {
        long arg1 = GetArg(opcode, 1);
        long arg2 = GetArg(opcode, 2);
        Set(opcode, 3, arg1 + arg2);
    }

    private void Mul(string opcode, int ip)
    {
        long arg1 = GetArg(opcode, 1);
        long arg2 = GetArg(opcode, 2);
        Set(opcode, 3, arg1 * arg2);
    }

    private void Store(string opcode, int ip, long input)
    {
        Set(opcode, 1, input);
    }

    private long Load(string opcode, int ip)
    {
        long arg1 = GetArg(opcode, 1);
        return arg1;
    }

    private long GetArg(string opcode, int argno)
    {
        int ptr = ip + argno;
        char mode = opcode[3 - argno];
        int n = data.Count;
        switch (mode)
        {
            case '0':
                ptr = ptr < n ? (int)data[ptr] : 0;
                break;
            case '1':
                break;
            case '2':
                ptr = ptr < n ? (int)data[ptr] : 0;
                ptr += offset;
                break;
            default:
                Console.WriteLine($"Unexpected parameter mode: {mode}");
                break;
        }
        return ptr < n ? data[ptr] : 0;
    }

    private void Set(string opcode, int argno, long val)
    {
        int ptr = ip + argno;
        char mode = opcode[3 - argno];
        int n = data.Count;
        ptr = (int)data[ptr];

        if (mode == '2')
            ptr += offset;

        if (ptr >= n)
            data.AddRange(new long[ptr - n + 1]);

        data[ptr] = val;
    }
}

class Program
{
    static void Main(string[] args)
    {
        string filename = args[0];

        Machine p = new(filename);
        Machine q = p.Clone();

        long output = p.Execute(1);
        Console.WriteLine($"Part 1: {output}");

        output = q.Execute(2);
        Console.WriteLine($"Part 2: {output}");
    }
}