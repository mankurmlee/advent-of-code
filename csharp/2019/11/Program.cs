class Intcode
{
    private int ip = 0;
    private int basePtr = 0;
    private Queue<long> inputBuffer = [];
    private long output = 0;
    private List<long> program;

    public Intcode(string filename)
    {
        string[] lines = File.ReadAllLines(filename);
        string[] nums = lines[0].Split(",");

        List<long> data = [];

        foreach (string num in nums)
            data.Add(long.Parse(num));

        this.program = data;
    }

    public Intcode(List<long> data)
    {
        this.program = data;
    }

    public Intcode Clone()
    {
        List<long> data = new(this.program);
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
                    Console.WriteLine($"Unexpected opcode: {opcode[4]}");
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
                Console.WriteLine($"Unexpected parameter mode: {mode}");
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
}

class Robot(string filename)
{
    private readonly static Vec[] DIRVEC = [
        new (0, -1),
        new (1, 0),
        new (0, 1),
        new (-1, 0),
    ];
    private Vec pos = new();
    private int dir = 0;
    private HashSet<Vec> white = [];
    private HashSet<Vec> painted = [];
    private Intcode program = new(filename);

    public override string ToString()
    {
        int w = 0;
        int h = 0;
        foreach (Vec pos in white)
        {
            if (pos.X > w)
                w = pos.X;
            if (pos.Y > h)
                h = pos.Y;
        }

        List<string> rows = [];
        for (int y = 0; y <= h; y++)
        {
            List<char> row = [];
            for (int x = 0; x <= w; x++)
            {
                char ch = white.Contains(new Vec(x, y)) ? '#' : ' ';
                row.Add(ch);
            }
            rows.Add(string.Concat(row));
        }
        return string.Join(Environment.NewLine, rows);
    }

    public void PaintHull()
    {
        bool halted = false;
        long input = 0;
        long colour = 0;
        long turn = 0;
        while (true)
        {
            input = white.Contains(pos) ? 1 : 0;
            program.Push(input);
            (halted, colour) = program.Execute();
            if (halted) return;
            (_, turn) = program.Execute();
            Paint(colour);
            TurnWalk(turn);
        }
    }

    private void TurnWalk(long turn)
    {
        int delta = turn == 0 ? -1 : 1;
        dir = (dir + 4 + delta) % 4;
        pos.AddEquals(DIRVEC[dir]);
    }

    public void Paint(long colour)
    {
        if (colour == 0)
            white.Remove(pos);
        else
            white.Add(pos);
        painted.Add(pos);
    }

    public int CountPainted()
    {
        return painted.Count;
    }
}

struct Vec(int x, int y)
{
    public int X = x;
    public int Y = y;
    public Vec(): this(0, 0) {}

    public void AddEquals(Vec other)
    {
        X += other.X;
        Y += other.Y;
    }

    public override string ToString()
    {
        return $"{X}, {Y}";
    }
}

class Program
{
    static void Main(string[] args)
    {
        Robot r = new(args[0]);
        r.PaintHull();
        int n = r.CountPainted();
        Console.WriteLine($"Part 1: {n}");

        r = new(args[0]);
        r.Paint(1);
        r.PaintHull();
        Console.WriteLine("Part 2:");
        Console.WriteLine(r);
    }
}