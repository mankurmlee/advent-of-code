class Image
{
    private readonly Layer[] layers;

    public Image(int w, int h, string filename)
    {
        string[] lines = File.ReadAllLines(filename);
        char[] data = lines[0].ToCharArray();

        int size = w * h;
        int n = data.Length / size;
        layers = new Layer[n];

        for (int i = 0; i < n; i++)
        {
            int p = i * size;
            int q = p + size;
            layers[i] = new(w, h, data[p..q]);
        }
    }

    public int PartOne()
    {
        Layer l = FewestLayer('0');
        return l.Count('1') * l.Count('2');
    }

    private Layer FewestLayer(char ch)
    {
        int min = -1;
        int idx = -1;

        int n = layers.Length;
        for (int i = 0; i < n; i++)
        {
            int count = layers[i].Count(ch);
            if (min < 0 || min > count)
            {
                min = count;
                idx = i;
            }
        }

        return layers[idx];
    }

    public Layer Flatten()
    {
        int n = layers.Length;
        Layer layer = new(layers[n-1]);
        for (int i = n - 2; i >= 0; i--)
        {
            layer.Merge(layers[i]);
        }
        return layer;
    }
}

class Layer(int w, int h, char[] data)
{
    private int width = w;
    private int height = h;
    private char[] data = data;

    public Layer(Layer other) : this(
        other.width,
        other.height,
        new char[other.data.Length])
    {
        Array.Copy(other.data, this.data, other.data.Length);
    }

    public int Count(char ch)
    {
        int count = 0;
        foreach (char c in data)
        {
            if (c == ch) count++;
        }
        return count;
    }

    public override string ToString()
    {
        List<string> result = [];
        for (int y = 0; y < height; y++)
        {
            int u = y * width;
            int v = u + width;
            string txt = new string(data[u..v])
                .Replace('0', ' ')
                .Replace('1', '#')
                .Replace('2', '.');
            result.Add(txt);
        }
        return string.Join(Environment.NewLine, result);
    }

    public void Merge(Layer other)
    {
        int n = data.Length;
        for (int i = 0; i < n; i++)
        {
            char v = other.data[i];
            if (v != '2')
            {
                data[i] = v;
            }
        }
    }
}

class Program
{
    static void Main(string[] args)
    {
        Image img = new(25, 6, args[0]);
        Console.WriteLine($"Part 1: {img.PartOne()}");

        Console.WriteLine("Part 2:");
        Console.WriteLine(img.Flatten());
    }
}
