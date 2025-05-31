#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <regex>

using namespace std;

struct Vec {
	int x, y;

	int dist() {
		bool sx = x >= 0;
		bool sy = y >= 0;
		int ax = abs(x);
		int ay = abs(y);
		return sx == sy ? max(ax, ay) : ax + ay;
	}
};

static void simulate(vector<string> dirs) {
	int furthest = 0;
	Vec pos{0, 0};
	for (const string& s : dirs) {
		if (s == "n") {
			pos.y += 1;
		} else if (s == "s") {
			pos.y -= 1;
		} else if (s == "nw") {
			pos.x -= 1;
		} else if (s == "se") {
			pos.x += 1;
		} else if (s == "sw") {
			pos.x -= 1;
			pos.y -= 1;
		} else if (s == "ne") {
			pos.x += 1;
			pos.y += 1;
		}
		int d = pos.dist();
		if (d > furthest) {
			furthest = d;
		}
	}
	cout << pos.dist() << endl;
	cout << furthest << endl;
}

static vector<string> findAllMatches(const string& pattern, const string& s) {
	vector<string> out{};
	regex re(pattern);
	sregex_iterator begin(s.begin(), s.end(), re);
	sregex_iterator end;
	for (auto it = begin; it != end; it++) {
		out.push_back(it->str());
	}
	return out;
}

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

static vector<string> load(const string& f) {
	auto d = readFile(f);
	return findAllMatches("\\w+", d[0]);
}

int main(int argc, char* argv[])
{
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto dirs = load(argv[1]);
	simulate(dirs);
}
