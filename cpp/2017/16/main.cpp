#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <regex>
#include <unordered_set>

using namespace std;

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

class Puzzle {
private:
	vector<string> code;
	vector<char> data;
	int start;

public:
	Puzzle(size_t n, const vector<string>& prog) : code(prog), start(0) {
		for (int i = 0; i < n; i++) {
			data.push_back('a' + i);
		}
	}

	string toString() const {
		string s;
		auto n = data.size();
		for (int i = start; i < start + n; i++) {
			s += data.at(i%n);
		}
		return s;
	}

	void spin(string s) {
		auto d = findAllMatches("\\d+", s);
		auto x = stoi(d[0]);
		start -= x;
		auto n = data.size();
		while (start < 0) start += n;
		start %= n;
	}

	void exchange(string s) {
		auto d = findAllMatches("\\d+", s);
		auto i = start + stoi(d[0]);
		auto j = start + stoi(d[1]);
		auto n = data.size();
		swap(data[i%n], data[j%n]);
	}

	void partner(string s) {
		auto a = s.at(1);
		auto b = s.at(3);
		auto n = data.size();
		int p = -1;
		int q = -1;
		for (size_t i = 0; i < n; i++) {
			char v = data.at(i);
			if (v == a) {
				p = i;
			} else 	if (v == b) {
				q = i;
			}
		}
		swap(data[p], data[q]);
	}

	void exec(string s) {
		char c = s.at(0);
		if (c == 's') {
			spin(s);
		} else if (c == 'x') {
			exchange(s);
		} else if (c == 'p') {
			partner(s);
		}
	}

	void run() {
		for (string stmt : code) {
			exec(stmt);
		}
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

static Puzzle load(const string& f) {
	auto d = readFile(f);
	auto code = findAllMatches("[0-9a-z/]+", d[0]);
	Puzzle p(stoi(d[1]), code);
	return move(p);
}

int main(int argc, char* argv[]) {
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto puzzle = load(argv[1]);

	puzzle.run();
	string s = puzzle.toString();
	cout << s << endl;;

	vector<string> order;
	unordered_set<string> seen;
	while (seen.count(s) == 0) {
		order.push_back(s);
		seen.insert(s);
		puzzle.run();
		s = puzzle.toString();
	}
	auto n = order.size();
	int i = ((1000000000 % order.size()) + n - 1) % n;
	cout << order.at(i) << endl;
}
