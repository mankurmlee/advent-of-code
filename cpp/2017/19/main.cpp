#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <unordered_map>
#include "vec.cpp"

using namespace std;
using Grid = unordered_map<Vec,char,VecHash>;

struct Pipes {
	Vec start;
	Grid grid;
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

static Pipes load(const string& f) {
	Pipes p;
	int y = 0;
	for (const string& s : readFile(f)) {
		int x = 0;
		for (char ch : s) {
			if (ch != ' ') {
				p.grid.emplace(Vec{x, y}, ch);
				if (y == 0) {
					p.start = {x, y};
				}
			}
			x++;
		}
		y++;
	}
	return p;
}

static constexpr Vec DIRS[4] = { {0, 1}, {1, 0}, {0, -1}, {-1, 0} };

struct State {
	Vec pos;
	int facing;
	int cost;
};

static void search(const Pipes& p) {
	string letters;
	int cost = 0;
	vector<State> q{{p.start, 0, 1}};
	while (q.size() > 0) {
		State s = q.back();
		q.pop_back();
		if (s.cost > cost) {
			cost = s.cost;
		}
		char c = p.grid.at(s.pos);
		if (c != '+' && c != '|' && c != '-') {
			letters += c;
		}
		Vec pos1 = s.pos + DIRS[s.facing];
		if (p.grid.count(pos1)) {
			q.push_back(State{pos1, s.facing, s.cost + 1});
			continue;
		}
		if (c != '+') {
			continue;
		}
		for (int o : {1, 3}) {
			int f1 = (s.facing + o) % 4;
			pos1 = s.pos + DIRS[f1];
			if (p.grid.count(pos1)) {
				q.push_back(State{pos1, f1, s.cost + 1});
				continue;
			}
		}
	}
	cout << letters << endl;
	cout << cost << endl;
}

int main(int argc, char* argv[]) {
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto p = load(argv[1]);
	search(p);
}
