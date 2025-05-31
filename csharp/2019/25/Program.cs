using System.Text.RegularExpressions;

class Room(string name, string desc, string[] doors, List<string> items)
{
    public string Name = name;
    public string Desc = desc;
    public string[] Doors = doors;
    public List<string> Items = items;
}

record Link(string Room, string Door);

class Path(string room, string from, Path? last)
{
    public static readonly Dictionary<string,string> OPPOSITE = new()
    {
        {"north", "south"}, {"east", "west"}, {"south", "north"}, {"west", "east"},
    };
    public string Room = room;
    public string From = from;
    public Path? Last = last;
    public List<string> ToDoorList()
    {
        List<string> d = [];
        Path p = this;
        while (p.From != "" && p.Last != null)
        {
            d.Add(OPPOSITE[p.From]);
            p = p.Last;
        }
        d.Reverse();
        return d;
    }
}

class Graph
{
    public static readonly Dictionary<string,string> OPPOSITE = new()
    {
        {"north", "south"}, {"east", "west"}, {"south", "north"}, {"west", "east"},
    };
    public Dictionary<string,Room> Rooms = [];
    public Dictionary<Link,string> Links = [];

    public void AddRoom(Room r)
    {
        Rooms[r.Name] = r;
    }

    public void AddLink(Link l, string dest)
    {
        Links[l] = dest;
        string oppo = OPPOSITE[l.Door];
        Links[new(dest, oppo)] = l.Room;
    }

    public Path? FindPath(string start, string dest)
    {
        HashSet<string> seen = [];
        Queue<Path> q = new();
        q.Enqueue(new(start, "", null));
        while (q.Count > 0)
        {
            Path p = q.Dequeue();
            if (p.Room == dest)
                return p;
            foreach (string door in Rooms[p.Room].Doors)
            {
                string adj = Links[new(p.Room, door)];
                if (seen.Contains(adj)) continue;
                seen.Add(adj);
                q.Enqueue(new(adj, OPPOSITE[door], p));
            }
        }
        return null;
    }
}

class Droid
{
    public static readonly Dictionary<string,string> OPPOSITE = new()
    {
        {"north", "south"}, {"east", "west"}, {"south", "north"}, {"west", "east"},
    };

    public static readonly string[] AVOID = [
        "escape pod",
        "giant electromagnet",
        "infinite loop",
        "molten lava",
        "photons",
    ];

    private readonly Intcode Prog;
    public Path MyPath;
    public Graph MyGraph = new();
    public List<string> Inventory = [];

    public Droid(string filename)
    {
        Prog = new(filename);
        Prog.Execute();
        Room r = ParseRoom(false);
        MyPath = new(r.Name, "", null);
    }

    public void SendCommand(string cmd)
    {
        Console.WriteLine(cmd);
        Prog.WriteLine(cmd);
        Prog.Execute();
    }

    public Room ParseRoom(bool echo)
    {
        string room = Prog.ReadString();
        string[] data = room.Split(
            "\n\n",
            StringSplitOptions.RemoveEmptyEntries|
            StringSplitOptions.TrimEntries);

        if (data.Length > 4 && MyPath.Last != null)
            MyPath = MyPath.Last;

        if (echo)
            Console.WriteLine(room);

        string[] s = data[0].Split('\n');

        string pattern = @"==\s*(.*?)\s*==";
        Match match = Regex.Match(s[0], pattern);

        string name = match.Groups[1].Value;
        string desc = s[1];

        s = data[1].Split('\n');
        int n = s.Length;
        string[] doors = new string[n - 1];
        for (int i = 1; i < n; i++)
            doors[i-1] = s[i][2..];

        s = data[2].Split('\n');
        n = s.Length;
        List<string> items = [];
        for (int i = 1; i < n; i++)
            items.Add(s[i][2..]);

        Room r = new(name, desc, doors, items);
        MyGraph.AddRoom(r);
        return r;
    }

    public void CollectAll()
    {
        HashSet<Link> seen = [];
        Stack<Link> s = new();
        Room r = ParseRoom(true);
        foreach (string door in r.Doors)
        {
            Link a = new(r.Name, door);
            seen.Add(a);
            s.Push(a);
        }

        while (s.Count > 0)
        {
            Link l = s.Pop();
            Retrace(l);
            Move(l.Door);
            r = ParseRoom(true);
            MyGraph.AddLink(l, r.Name);
            string from = OPPOSITE[l.Door];
            MyPath = new(r.Name, from, MyPath);
            seen.Add(new(r.Name, from));
            TakeAll(r);
            foreach (string door in r.Doors)
            {
                Link a = new(r.Name, door);
                if (seen.Contains(a)) continue;
                seen.Add(a);
                s.Push(a);
            }
        }
    }

    private void TakeAll(Room r)
    {
        foreach (string item in r.Items)
        {
            if (AVOID.Contains(item)) continue;
            SendCommand($"take {item}");
            Inventory.Add(item);
            Console.WriteLine(Prog.ReadString());
        }
    }

    private void Retrace(Link o)
    {
        while (!o.Room.Equals(MyPath.Room) && MyPath.Last != null)
        {
            Move(MyPath.From);
            MyPath = MyPath.Last;
            Console.WriteLine($"Retraced back to {MyPath.Room}");
        }
    }

    private void Move(string door)
    {
        SendCommand(door);
    }

    public void Goto(string dest)
    {
        string start = MyPath.Room;
        Path? p = MyGraph.FindPath(start, dest);
        if (p == null) {
            Console.WriteLine("Couldn't find path!");
            return;
        }

        foreach (string door in p.ToDoorList())
            Move(door);
        MyPath = p;
    }

    // - jam
    // - bowl of rice
    // - antenna
    // - manifold
    // - hypercube
    // - dehydrated water
    // - candy cane
    // - dark matter
    public void TrySecurity()
    {
        SendCommand("drop bowl of rice");
        // SendCommand("drop hypercube");
        SendCommand("drop jam");
        // SendCommand("drop antenna");
        SendCommand("drop manifold");
        // SendCommand("drop dehydrated water");
        // SendCommand("drop candy cane");
        SendCommand("drop dark matter");
        Move("west");
        Console.WriteLine(Prog.ReadString());
    }
}

class Program
{
    static void Main(string[] args)
    {
        string filename = @"D:\aoc\2019\25\puzzle.txt";
        Droid d = new(filename);

        d.CollectAll();
        d.Goto("Security Checkpoint");
        d.TrySecurity();
    }
}