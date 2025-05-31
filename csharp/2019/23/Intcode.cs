
class Intcode
{
    private int Ip = 0;
    private int BasePtr = 0;
    private readonly Queue<long> InputBuffer = [];
    private readonly List<long> OutputBuffer = [];
    private readonly List<long> Program;

    public Intcode(string filename)
    {
        string[] lines = File.ReadAllLines(filename);
        string[] nums = lines[0].Split(",");

        List<long> data = [];

        foreach (string num in nums)
            data.Add(long.Parse(num));

        Program = data;
    }

    public Intcode(List<long> data)
    {
        Program = data;
    }

    public Intcode Clone()
    {
        List<long> data = new(Program);
        return new Intcode(data);
    }

    public void Push(long input)
    {
        InputBuffer.Enqueue(input);
    }

    public void WriteLine(string input)
    {
        foreach (char c in input)
            InputBuffer.Enqueue(c);
        InputBuffer.Enqueue(10);
    }

    public bool Execute()
    {
        OutputBuffer.Clear();
        while (Program[Ip] != 99)
        {
            string opcode = Program[Ip].ToString().PadLeft(5, '0');

            switch (opcode[4])
            {
                case '1':
                    Add(opcode);
                    Ip += 4;
                    break;
                case '2':
                    Mul(opcode);
                    Ip += 4;
                    break;
                case '3':
                    if (InputBuffer.Count == 0)
                        return false;
                    Store(opcode);
                    Ip += 2;
                    break;
                case '4':
                    OutputBuffer.Add(Load(opcode));
                    Ip += 2;
                    break;
                case '5':
                    Ip = JumpIfTrue(opcode);
                    break;
                case '6':
                    Ip = JumpIfFalse(opcode);
                    break;
                case '7':
                    LessThan(opcode);
                    Ip += 4;
                    break;
                case '8':
                    EqualsTo(opcode);
                    Ip += 4;
                    break;
                case '9':
                    AddBase(opcode);
                    Ip += 2;
                    break;
                default:
                    System.Console.WriteLine($"Unexpected opcode: {opcode[4]}");
                    break;
            }
        }
        return true;
    }

    private void AddBase(string opcode)
    {
        long arg1 = GetArg(opcode, 1);
        BasePtr += (int)arg1;
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
        return Ip + 3;
    }

    private int JumpIfFalse(string opcode)
    {
        long arg1 = GetArg(opcode, 1);
        long arg2 = GetArg(opcode, 2);
        if (arg1 == 0)
            return (int)arg2;
        return Ip + 3;
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
        long input = InputBuffer.Dequeue();
        Set(opcode, 1, input);
    }

    private long Load(string opcode)
    {
        long arg1 = GetArg(opcode, 1);
        return arg1;
    }

    private long GetArg(string opcode, int argno)
    {
        int ptr = Ip + argno;
        char mode = opcode[3 - argno];
        int n = Program.Count;
        switch (mode)
        {
            case '0':
                ptr = ptr < n ? (int)Program[ptr] : 0;
                break;
            case '1':
                break;
            case '2':
                ptr = ptr < n ? (int)Program[ptr] : 0;
                ptr += BasePtr;
                break;
            default:
                System.Console.WriteLine($"Unexpected parameter mode: {mode}");
                break;
        }
        return ptr < n ? Program[ptr] : 0;
    }

    private void Set(string opcode, int argno, long val)
    {
        int ptr = Ip + argno;
        char mode = opcode[3 - argno];
        int n = Program.Count;
        ptr = (int)Program[ptr];

        if (mode == '2')
            ptr += BasePtr;

        if (ptr >= n)
            Program.AddRange(new long[ptr - n + 1]);

        Program[ptr] = val;
    }

    public void Poke(int ptr, int val)
    {
        Program[ptr] = val;
    }

    public string ReadString()
    {
        char[] chars = new char[OutputBuffer.Count];
        for (int i = 0; i < OutputBuffer.Count; i++)
            chars[i] = (char)OutputBuffer[i];
        return new string(chars);
    }
    public List<long> Read()
    {
        return new List<long>(OutputBuffer);
    }
}
