#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <regex>
#include <unordered_map>

using namespace std;

template <typename T>
using Matrix = vector<vector<T>>;

class Cpu {
private:
	unordered_map<string, int> regs;
	int largestIR = 0;

public:
	Cpu() {}

	bool test(const vector<string>& stmt) {
		string r = stmt.at(4);
		string op = stmt.at(5);
		int v = stoi(stmt.at(6));
		if (op == "==") {
			return regs[r] == v;
		} else if (op == "!=") {
			return regs[r] != v;
		} else if (op == "<") {
			return regs[r] < v;
		} else if (op == ">") {
			return regs[r] > v;
		} else if (op == "<=") {
			return regs[r] <= v;
		} else if (op == ">=") {
			return regs[r] >= v;
		}
		cerr << "Unexpected operator: " << op << endl;
		throw runtime_error("Unexpected operator");
	}

	void exec(const vector<string>& stmt) {
		if (test(stmt)) {
			string r = stmt.at(0);
			string op = stmt.at(1);
			int v = stoi(stmt.at(2));
			if (op == "inc") {
				regs[r] += v;
			} else if (op == "dec") {
				regs[r] -= v;
			}
			if (regs[r] > largestIR) {
				largestIR = regs[r];
			}
		}
	}

	void run(const Matrix<string>& prog) {
		for (const auto& stmt : prog) {
			exec(stmt);
		}
		int out = 0;
		for (const auto& [_, v] : regs) {
			if (v > out) {
				out = v;
			}
		}
		cout << out << endl;
		cout << largestIR << endl;
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

static Matrix<string> load(const string& f) {
	Matrix<string> p{};
	for (const string& l : readFile(f)) {
		auto stmt = findAllMatches("\\S+", l);
		p.push_back(stmt);
	}
	return p;
}

int main(int argc, char* argv[])
{
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto puzzle = load(argv[1]);
	Cpu cpu;
	cpu.run(puzzle);
}
