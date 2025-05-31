#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <regex>
#include <unordered_map>
#include <unordered_set>

using namespace std;

class Puzzle {
private:
	unordered_map<string, int> nodes{};
	unordered_map<string, vector<string>> childs{};
	unordered_map<string, int> weights{};

public:
	Puzzle() {}

	void add(const string& n, int w, const vector<string>& cs) {
		nodes.emplace(n, w);
		childs.emplace(n, cs);
	}

	string root() const {
		unordered_set<string> isChild{};
		for (const auto& [_, cs] : childs) {
			for (const string& c : cs) {
				isChild.insert(c);
			}
		}
		for (const auto& [k, _] : nodes) {
			if (isChild.count(k) == 0) {
				return k;
			}
		}
		return "";
	}

	int getWeight(const string& name) {
		auto it = weights.find(name);
		if (it != weights.end()) {
			return it->second;
		}
		int w = nodes[name];
		for (const string& c : childs[name]) {
			w += getWeight(c);
		}
		weights[name] = w;
		return w;
	}

	bool isUnbalanced(const string& name) {
		unordered_map<int, vector<string>> ws{};
		for (const auto& c : childs[name]) {
			if (isUnbalanced(c)) {
				return true;
			}
			int w = getWeight(c);
			ws[w].push_back(c);
		}
		if (ws.size() != 2) {
			return false;
		}
		int w, otherW;
		for (const auto& [k, v] : ws) {
			if (v.size() == 1) {
				w = nodes[v[0]] - k;
			} else {
				otherW = k;
			}
		}
		cout << (w + otherW) << endl;
		return true;
	}
};

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

static Puzzle load(const string& f) {
	Puzzle p{};
	for (const string& l : readFile(f)) {
		auto ws = findAllMatches("\\w+", l);
		vector<string> cs(ws.begin() + 2, ws.end());
		p.add(ws[0], stoi(ws[1]), cs);
	}
	return move(p);
}

int main(int argc, char* argv[])
{
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto puzzle = load(argv[1]);
	string root = puzzle.root();
	cout << root << endl;
	puzzle.isUnbalanced(root);
}
