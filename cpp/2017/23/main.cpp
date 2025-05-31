#include <vector>
#include <iostream>
#include <fstream>
#include <regex>
#include "thread.cpp"

using namespace std;

template <typename T>
using Matrix = vector<vector<T>>;

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
	Matrix<string> prog;
	for (const string& l : readFile(f)) {
		prog.push_back(findAllMatches("\\S+", l));
	}
	return prog;
}

static void partOne(Matrix<string> prog) {
	Thread c(prog, 0, false);
	c.run();
	cout << c.get("mul") << endl;
}

static bool notPrime(int n) {
	for (int i = 2; i < n; i++) {
		if (n % i == 0) {
			return true;
		}
	}
	return false;
}

static void partTwo() {
	int count = 0;
	for (int t = 108400; t <= 125400; t += 17) {
		if (notPrime(t)) {
			count++;
		}
	}
	cout << count << endl;
}

int main(int argc, char* argv[]) {
	auto prog = load("puzzle.txt");
	partOne(prog);
	partTwo();
}
