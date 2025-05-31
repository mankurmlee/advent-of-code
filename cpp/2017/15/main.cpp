#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <regex>
#include <cstdint>

using namespace std;

class Generator {
private:
	uint64_t factor;
	uint64_t value;

public:
	Generator(long f, long v) : factor(f), value(v) {}

	int next() {
		value = (value * factor) % 2147483647;
		return value;
	}
};

static void partOne(pair<int, int> config) {
	Generator a(16807, config.first);
	Generator b(48271, config.second);
	int count = 0;
	for (int i = 0; i < 40000000; i++) {
		uint32_t a1 = a.next();
		uint32_t b1 = b.next();
		if ((a1 & 0xffff) == (b1 & 0xffff)) {
			count++;
		}
	}
	cout << count << endl;
}

static void partTwo(pair<int, int> config) {
	Generator a(16807, config.first);
	Generator b(48271, config.second);
	int count = 0;
	for (int i = 0; i < 5000000; i++) {
		uint32_t a1;
		do {
			a1 = a.next();
		} while (a1 % 4 != 0);
		uint32_t b1;
		do {
			b1 = b.next();
		} while (b1 % 8 != 0);
		if ((a1 & 0xffff) == (b1 & 0xffff)) {
			count++;
		}
	}
	cout << count << endl;
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

static vector<int> parseInts(const string& s) {
	vector<int> out{};
	for (const string& s : findAllMatches("\\d+", s)) {
		out.push_back(stoi(s));
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

static pair<int,int> load(const string& f) {
	auto d = readFile(f);
	auto a = parseInts(d[0]);
	auto b = parseInts(d[1]);
	return {a[0], b[0]};
}

int main(int argc, char* argv[]) {
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto config = load(argv[1]);
	partOne(config);
	partTwo(config);
}
