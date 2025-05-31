class SpringDroid(string filename)
{
    private readonly Intcode Backup = new(filename);

    public string Walk()
    {
        Intcode Prog = Backup.Clone();
        // !(a && b && c) && d
        Prog.WriteLine("NOT A J");
        Prog.WriteLine("NOT J J");
        Prog.WriteLine("AND B J");
        Prog.WriteLine("AND C J");
        Prog.WriteLine("NOT J J");
        Prog.WriteLine("AND D J");
        Prog.WriteLine("WALK");
        Prog.Execute();

        long last = Prog.Read().Last();
        if (last < 256)
            return Prog.ReadString();
        return last.ToString();
    }

    public string Run()
    {
        Intcode Prog = Backup.Clone();
        // !(a && b && c) && d && (e || h)
        Prog.WriteLine("NOT A J");
        Prog.WriteLine("NOT J J");
        Prog.WriteLine("AND B J");
        Prog.WriteLine("AND C J");
        Prog.WriteLine("NOT J J");
        Prog.WriteLine("AND D J");
        Prog.WriteLine("NOT E T");
        Prog.WriteLine("NOT T T");
        Prog.WriteLine("OR H T");
        Prog.WriteLine("AND T J");
        Prog.WriteLine("RUN");
        Prog.Execute();

        long last = Prog.Read().Last();
        if (last < 256)
            return Prog.ReadString();
        return last.ToString();
    }
}

class Program
{
    static void Main(string[] args)
    {
        string filename = @"D:\aoc\2019\21\puzzle.txt";
        SpringDroid d = new(filename);

        string res = d.Walk();
        Console.WriteLine($"Part 1: {res}");

        res = d.Run();
        Console.WriteLine($"Part 2: {res}");
    }
}