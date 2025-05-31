#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <unordered_set>
#include "knothash.cpp"

using namespace std;

struct Vec {
    int x, y;

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

using Region = unordered_set<Vec, VecHash>;

class Maze {
private:
	Region used;
	static constexpr Vec DIRS[] = {{0, -1}, {1, 0}, {0, 1}, {-1, 0}};

public:
	Maze(string key) {
		KnotHash knot;
		for (int y = 0; y < 128; y++) {
			int i = 0;
			for (int v : knot.hash(key + "-" + to_string(y))) {
				int mask = 1 << 7;
				for (int j = 0; j < 8; j++) {
					if (v & mask) {
						used.insert({i + j, y});
					}
					mask >>= 1;
				}
				i += 8;
			}
		}
	}

	int size() const {
		return used.size();
	}

	Region getRegion(Vec start) const {
		Region seen{start};
		vector<Vec> q{start};
		while (q.size() > 0) {
			auto s = q.back();
			q.pop_back();
			for (auto d : DIRS) {
				auto s1 = s + d;
				if (!used.count(s1)) continue;
				if (seen.count(s1)) continue;
				seen.insert(s1);
				q.push_back(s1);
			}
		}
		return seen;
	}

	int regionCount() const {
		int c = 0;
		Region seen;
		for (Vec v : used) {
			if (seen.count(v)) continue;
			Region r = getRegion(v);
			seen.insert(r.begin(), r.end());
			c++;
		}
		return c;
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

static string load(const string& f) {
	auto d = readFile(f);
	return d[0];
}

int main(int argc, char* argv[]) {
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}

	Maze m(load(argv[1]));
	cout << m.size() << endl;
	cout << m.regionCount() << endl;
}
