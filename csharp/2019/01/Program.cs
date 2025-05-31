class Puzzle
{
    public List<int> Masses { get; }

    public Puzzle(string filename)
    {
        string[] lines = File.ReadAllLines(filename);
        List<int> masses = [];
        foreach (string line in lines)
        {
            masses.Add(int.Parse(line));
        }

        Masses = masses;
    }

    public override string ToString()
    {
        return "[" + string.Join(" ", Masses) + "]";
    }

    public int FuelRequirement(bool inclFuelMass)
    {
        int sum = 0;

        foreach (int m in Masses)
        {
            int c = m;
            int r = 0;
            while (true)
            {
                r = c / 3 - 2;
                if (r <= 0) break;
                // Console.WriteLine("Adding " + r + " fuel for " + m);
                sum += r;
                if (!inclFuelMass) break;
                c = r;
            }
        }

        return sum;
    }
}

class Program
{
    static void Main(string[] args)
    {
        string filename = args[0];
        Puzzle p = new(filename);
        Console.WriteLine(p);

        Console.WriteLine("Part 1: " + p.FuelRequirement(false));
        Console.WriteLine("Part 2: " + p.FuelRequirement(true));
    }
}
