

using System.Collections.Concurrent;
using System.Globalization;

class Puzzle
{
    public int Lo { get; }
    public int Hi { get; }

    public Puzzle(string filename)
    {
        string[] lines = File.ReadAllLines(filename);
        var data = lines[0].Split('-');
        Lo = int.Parse(data[0]);
        Hi = int.Parse(data[1]);
    }

    public int CountMatches(bool exact)
    {
        int count = 0;
        for (int i = Lo; i <= Hi; i++)
        {
            if (MeetsCriteria(i, exact))
            {
                count++;
            }
        }
        return count;
    }

    private static bool MeetsCriteria(int n, bool exact)
    {
        int[] counts = new int[10];
        int last = -1;
        foreach (var digit in Digits(n))
        {
            if (digit < last)
            {
                return false;
            }
            counts[digit]++;
            last = digit;
        }
        foreach (var i in counts)
        {
            if (i == 2 || !exact && i > 2)
            {
                return true;
            }
        }
        return false;
    }

    private static int[] Digits(int n)
    {
        string numstr = n.ToString();
        int[] digits = new int[numstr.Length];
        for (int i = 0; i < numstr.Length; i++)
        {
            digits[i] = int.Parse(numstr[i].ToString());
        }
        return digits;
    }
}

class Program
{
    static void Main(string[] args)
    {
        Puzzle p = new(args[0]);
        Console.WriteLine($"Part 1: {p.CountMatches(false)}");
        Console.WriteLine($"Part 2: {p.CountMatches(true)}");
    }
}
