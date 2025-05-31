class Intcode
{
    private int ip = 0;
    private int basePtr = 0;
    private readonly Queue<long> inputBuffer = [];
    private long output = 0;
    private readonly List<long> program;

    public Intcode(string filename)
    {
        string[] lines = File.ReadAllLines(filename);
        string[] nums = lines[0].Split(",");

        List<long> data = [];

        foreach (string num in nums)
            data.Add(long.Parse(num));

        program = data;
    }

    public Intcode(List<long> data)
    {
        program = data;
    }

    public Intcode Clone()
    {
        List<long> data = new(program);
        return new Intcode(data);
    }

    public void Push(long input)
    {
        inputBuffer.Enqueue(input);
    }

    public (bool, long) Execute()
    {
        while (program[ip] != 99)
        {
            string opcode = program[ip].ToString().PadLeft(5, '0');

            switch (opcode[4])
            {
                case '1':
                    Add(opcode);
                    ip += 4;
                    break;
                case '2':
                    Mul(opcode);
                    ip += 4;
                    break;
                case '3':
                    Store(opcode);
                    ip += 2;
                    break;
                case '4':
                    output = Load(opcode);
                    ip += 2;
                    return (false, output);
                case '5':
                    ip = JumpIfTrue(opcode);
                    break;
                case '6':
                    ip = JumpIfFalse(opcode);
                    break;
                case '7':
                    LessThan(opcode);
                    ip += 4;
                    break;
                case '8':
                    EqualsTo(opcode);
                    ip += 4;
                    break;
                case '9':
                    AddBase(opcode);
                    ip += 2;
                    break;
                default:
                    System.Console.WriteLine($"Unexpected opcode: {opcode[4]}");
                    break;
            }
        }
        return (true, 0);
    }

    private void AddBase(string opcode)
    {
        long arg1 = GetArg(opcode, 1);
        basePtr += (int)arg1;
    }

    private void EqualsTo(string opcode)
    {
        long arg1 = GetArg(opcode, 1);
        long arg2 = GetArg(opcode, 2);
        Set(opcode, 3, arg1 == arg2 ? 1 : 0);
    }

    private void LessThan(string opcode)
    {
        long arg1 = GetArg(opcode, 1);
        long arg2 = GetArg(opcode, 2);
        Set(opcode, 3, arg1 < arg2 ? 1 : 0);
    }

    private int JumpIfTrue(string opcode)
    {
        long arg1 = GetArg(opcode, 1);
        long arg2 = GetArg(opcode, 2);
        if (arg1 != 0)
            return (int)arg2;
        return ip + 3;
    }

    private int JumpIfFalse(string opcode)
    {
        long arg1 = GetArg(opcode, 1);
        long arg2 = GetArg(opcode, 2);
        if (arg1 == 0)
            return (int)arg2;
        return ip + 3;
    }

    private void Add(string opcode)
    {
        long arg1 = GetArg(opcode, 1);
        long arg2 = GetArg(opcode, 2);
        Set(opcode, 3, arg1 + arg2);
    }

    private void Mul(string opcode)
    {
        long arg1 = GetArg(opcode, 1);
        long arg2 = GetArg(opcode, 2);
        Set(opcode, 3, arg1 * arg2);
    }

    private void Store(string opcode)
    {
        long input = inputBuffer.Dequeue();
        Set(opcode, 1, input);
    }

    private long Load(string opcode)
    {
        long arg1 = GetArg(opcode, 1);
        return arg1;
    }

    private long GetArg(string opcode, int argno)
    {
        int ptr = ip + argno;
        char mode = opcode[3 - argno];
        int n = program.Count;
        switch (mode)
        {
            case '0':
                ptr = ptr < n ? (int)program[ptr] : 0;
                break;
            case '1':
                break;
            case '2':
                ptr = ptr < n ? (int)program[ptr] : 0;
                ptr += basePtr;
                break;
            default:
                System.Console.WriteLine($"Unexpected parameter mode: {mode}");
                break;
        }
        return ptr < n ? program[ptr] : 0;
    }

    private void Set(string opcode, int argno, long val)
    {
        int ptr = ip + argno;
        char mode = opcode[3 - argno];
        int n = program.Count;
        ptr = (int)program[ptr];

        if (mode == '2')
            ptr += basePtr;

        if (ptr >= n)
            program.AddRange(new long[ptr - n + 1]);

        program[ptr] = val;
    }

    public void Poke(int ptr, int val)
    {
        program[ptr] = val;
    }
}
