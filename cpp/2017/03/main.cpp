#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <unordered_map>

using namespace std;

struct Vec {
    int x;
    int y;

    Vec(int x, int y) : x(x), y(y) {}

    Vec operator+(const Vec& o) const {
        return Vec{x + o.x, y + o.y};
    }

    bool operator==(const Vec& o) const {
        return x == o.x && y == o.y;
    }
};

struct VecHash {
    size_t operator()(const Vec& v) const {
        return hash<int>()(v.x) ^ (hash<int>()(v.y) << 1);
    }
};

static Vec DIR[4] = {{0, -1}, {-1, 0}, {0, 1}, {1, 0}};
static Vec SUR[8] = {{0, -1}, {-1, 0}, {0, 1}, {1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}};

class Memory {
private:
	unordered_map<Vec,int,VecHash> grid;
    Vec pos;
    int facing;

    void move() {
        int f1 = (facing + 1) % 4;
        Vec p1 = pos + DIR[f1];
        if (grid.count(p1)) {
            f1 = facing;
            p1 = pos + DIR[f1];
        }
        facing = f1;
        pos = p1;
    }

    int value() const {
        int sum = 0;
        for (const Vec& v : SUR) {
            Vec p1 = pos + v;
            auto it = grid.find(p1);
            if (it != grid.end()) {
                sum += it->second;
            }
        }
        return sum;
    }

public:
    Memory() : pos({0, 0}), facing(2) {
        grid.emplace(pos, 1);
    }

    int next() {
        move();
        int v = value();
        grid.emplace(pos, v);
        return v;
    }
};

static vector<string> readFile(const string& filename) {
    vector<string> lines{};
    ifstream file(filename);
    if (!file) {
        cerr << "Error: Could not open file " << filename << endl;
        throw runtime_error("File not found!");
    }
    string line;
    while (getline(file, line)) {
        lines.push_back(line);
    }
    file.close();
    return lines;
}

static int load(const string& f) {
    auto data = readFile(f);
    return stoi(data[0]);
}

static void partOne(const int& num) {
    int i = 1;
    while (i*i < num) {i += 2;}
    int i1 = i - 1;
    int i2 = i - 2;
    int rem = (num - i2*i2) % i1;
    cout << ((rem >= i1 >> 1) ? rem : i1 - rem) << endl;
}

static void partTwo(const int& num) {
    Memory m;
    int val = 1;
    while (val < num) {
        val = m.next();
    }
    cout << val << endl;
}

int main(int argc, char* argv[])
{
    if (argc != 2) {
        cerr << "Usage: " << argv[0] << " <filename>" << endl;
        return 1;
    }
    auto data = load(argv[1]);
    partOne(data);
    partTwo(data);
}
