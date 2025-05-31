class IntProgram
{
    private int ip = 0;
    private int output = 0;
    private int[] buffer = [];
    private int bufferIdx = 0;
    public List<int> Data { get; }

    public IntProgram(string filename)
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

    public IntProgram(List<int> data)
    {
        Data = data;
    }

    public IntProgram Clone()
    {
        List<int> data = new(Data);
        return new IntProgram(data);
    }

    public (bool, int) Execute(int[] input)
    {
        buffer = input;
        bufferIdx = 0;

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
                    Store(opcode, ip);
                    ip += 2;
                    break;
                case '4':
                    output = Load(opcode, ip);
                    ip += 2;
                    return (false, output);
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
        return (true, output);
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

    private void Store(string opcode, int ip)
    {
        int ptr = Data[ip+1];
        int input = buffer[bufferIdx++];
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

class Circuit
{
    private readonly IntProgram program;
    private IntPermer permer = new();
    public IntProgram[] Amps { get; }

    public Circuit(IntProgram p, int n)
    {
        program = p;
        Amps = new IntProgram[n];
    }

    public void Reset()
    {
        int n = Amps.Length;
        for (int i = 0; i < n; i++)
        {
            Amps[i] = program.Clone();
        }
    }

    public int Run(int[] phases, bool feedback)
    {
        Reset();

        int result = 0;
        int n = Amps.Length;

        bool halted = false;
        for (int i = 0; i < n; i++)
        {
            (halted, result) = Amps[i].Execute([phases[i], result]);
        }
        if (!feedback) return result;

        while (!halted)
        {
            for (int i = 0; i < n; i++)
            {
                (halted, result) = Amps[i].Execute([result]);
            }
        }

        return result;
    }

    public int FindMaxThrust(bool feedback)
    {
        List<int> phases = [0, 1, 2, 3, 4];
        if (feedback)
        {
            phases = [5, 6, 7, 8, 9];
        }
        int maxThrust = 0;
        foreach (var phase_order in permer.GenPerms(phases))
        {
            int result = Run([.. phase_order], feedback);
            if (result > maxThrust)
            {
                maxThrust = result;
            }
        }

        return maxThrust;
    }
}

class IntPermer
{
    private Dictionary<string, List<int>[]> memo = [];

    public List<int>[] GenPerms(List<int> nums)
    {
        int n = nums.Count;
        if (n == 1)
        {
            return [nums];
        }

        // Assume items are sorted....
        // nums.Sort();
        string key = string.Join(",", nums);
        if (memo.ContainsKey(key))
        {
            return memo[key];
        }

        int nRes = Factorial(n);
        List<int>[] result = new List<int>[nRes];
        int rIdx = 0;

        for (int i = 0; i < n; i++)
        {
            int v = nums[i];
            List<int> newNums = new(nums);
            newNums.RemoveAt(i);
            foreach (List<int> p in GenPerms(newNums))
            {
                List<int> newPerm = new(p) { v };
                result[rIdx++] = newPerm;
            }
        }

        memo[key] = result;
        return memo[key];
    }

    private static int Factorial(int n)
    {
        int result = 1;
        for (int i = 2; i <= n; i++)
        {
            result *= i;
        }
        return result;
    }
}

class Program
{
    static void Main(string[] args)
    {
        IntProgram p = new(args[0]);
        Circuit c = new(p, 5);

        int thrust = c.FindMaxThrust(false);
        Console.WriteLine($"Part 1: {thrust}");

        int thrust2 = c.FindMaxThrust(true);
        Console.WriteLine($"Part 2: {thrust2}");
    }
}