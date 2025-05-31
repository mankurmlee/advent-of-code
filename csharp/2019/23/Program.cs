struct Packet(long dest, long x, long y)
{
    public long Dest = dest;
    public long X = x;
    public long Y = y;

    public override readonly string ToString()
    {
        return $"<{Dest}, {X}, {Y}>";
    }
}

class NIC
{
    private readonly Intcode Prog;
    public int Address;
    private readonly Queue<Packet> Incoming = [];

    public NIC(string filename)
    {
        Prog = new(filename);
        Address = 0;
    }

    public NIC(Intcode p, int addr)
    {
        Prog = p;
        Address = addr;
    }

    public NIC Clone(int addr)
    {
        Intcode p = Prog.Clone();
        return new(p, addr);
    }

    public void Init()
    {
        Prog.Push(Address);
    }

    public List<Packet> Run()
    {
        if (Incoming.Count > 0)
            while (Incoming.Count > 0)
            {
                Packet p = Incoming.Dequeue();
                Prog.Push(p.X);
                Prog.Push(p.Y);
            }
        else
            Prog.Push(-1);

        Prog.Execute();

        List<long> runOutput = Prog.Read();
        int n = runOutput.Count;

        List<Packet> res = [];
        for (int i = 0; i < n; i += 3)
        {
            long addr = runOutput[i];
            long x = runOutput[i+1];
            long y = runOutput[i+2];
            res.Add(new(addr, x, y));
        }

        return res;
    }

    public void Message(Packet p)
    {
        Incoming.Enqueue(p);
    }
}

class Network
{
    private readonly NIC[] Machines;
    private Packet NAT;
    private Packet LastNAT;
    public Network(string filename, int n)
    {
        NIC m = new(filename);

        Machines = new NIC[n];
        for (int i = 0; i < n; i++)
        {
            NIC clone = m.Clone(i);
            clone.Init();
            Machines[i] = clone;
        }
    }

    public int OneCycle()
    {
        List<Packet> msgs = [];
        foreach (NIC m in Machines)
            msgs.AddRange(m.Run());

        foreach (Packet p in msgs)
        {
            if (p.Dest > Machines.Length)
            {
                NAT = p;
                continue;
            }
            Machines[p.Dest].Message(p);
        }

        return msgs.Count;
    }

    public (long, long) Simulate()
    {
        long first = 0;
        int n;
        int last = 0;
        while (true)
        {
            n = OneCycle();
            if (n == 0 && last == 0)
            {
                if (NAT.Y == LastNAT.Y)
                    return (first, NAT.Y);
                if (first == 0)
                    first = NAT.Y;
                Console.WriteLine($"NAT sends packet {NAT}");
                Machines[0].Message(NAT);
                LastNAT = NAT;
            }
            last = n;
        }
    }
}

class Program
{
    static void Main(string[] args)
    {
        string filename = @"D:\aoc\2019\23\puzzle.txt";

        Network n = new(filename, 50);
        (long first, long last) = n.Simulate();
        Console.WriteLine($"Part 1: {first}");
        Console.WriteLine($"Part 2: {last}");
    }
}