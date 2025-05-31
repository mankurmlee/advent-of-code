#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <unordered_set>
#include <unordered_map>
#include "vec.cpp"

using namespace std;

struct Config {
	Vec start;
	unordered_set<Vec,VecHash> infected;
};

struct Carrier {
	Vec pos;
	int facing;
	int infections;
};

static constexpr Vec DIRS[] = {{0, -1}, {1, 0}, {0, 1}, {-1, 0}};

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

static Config load(const string& f) {
	Config cfg;
	auto lines = readFile(f);
	int m = lines.size() >> 1;
	cfg.start = Vec{m, m};
	int y = 0;
	for (const string& l : lines) {
		int x = 0;
		for (char c : l) {
			if (c == '#') {
				cfg.infected.insert({x, y});
			}
			x++;
		}
		y++;
	}
	return move(cfg);
}

void simulate(const Config& cfg, int bursts, int step) {
	unordered_map<Vec,int,VecHash> grid;
	for (const Vec& v : cfg.infected) {
		grid.emplace(v,2);
	}
	Carrier c = {cfg.start, 0, 0};
	for (int i = 0; i < bursts; i++) {
		Carrier c1 = c;
		int v = grid[c.pos];
		int g1 = (grid[c.pos] + step) % 4;
		if (g1 == 2) {
			c1.infections++;
		}
		grid[c.pos] = g1;
		c1.facing = (c.facing + 3 + v) % 4;
		c1.pos = c.pos + DIRS[c1.facing];
		c = c1;
	}
	cout << c.infections << endl;
}

int main(int argc, char* argv[]) {
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto cfg = load(argv[1]);
	simulate(cfg, 10000, 2);
	simulate(cfg, 10000000, 1);
}
