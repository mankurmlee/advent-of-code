#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <regex>

using namespace std;

class Puzzle {
private:
	vector<int> buf{};
	vector<int> data{};
	vector<int> ascii{};
	bool part2 = false;
	size_t pos = 0;
	size_t skip = 0;

public:
	Puzzle(int n, vector<int> l, vector<int> a) : data(l), ascii(a) {
		for (int i = 0; i < n; i++) {
			buf.push_back(i);
		}
	}

	void reverse(int start, int end) {
		size_t n = buf.size();
		for (size_t i = start, j = end - 1; i < j; i++, j--) {
			swap(buf[i%n], buf[j%n]);
		}
	}

	void shuffle() {
		auto d = (part2 ? ascii : data);
		size_t n = buf.size();
		for (size_t len : d) {
			if (len > n) {
				continue;
			}
			reverse(pos, pos + len);
			pos = (pos + len + skip) % n;
			skip++;
		}
		if (!part2) {
			cout << (buf[0] * buf[1]) << endl;
		}
	}

	void partTwo() {
		part2 = true;
		size_t n = buf.size();
		if (n != 256) {
			return;
		}
		for (int i = 0; i < n; i++) {
			buf[i] = i;
		}
		pos = 0;
		skip = 0;
		for (int i = 0; i < 64; i++) {
			shuffle();
		}
		for (int i = 0; i < 16; i++) {
			int v = 0;
			for (int j = 0; j < 16; j++) {
				v ^= buf.at((i << 4) + j);
			}
			printf("%02x", v);
		}
		cout << endl;
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

static Puzzle load(string f) {
	auto d = readFile(f);
	vector<int> ascii{};
	for (char c : d[1]) {
		ascii.push_back(c);
	}
	for (int v : {17, 31, 73, 47, 23}) {
		ascii.push_back(v);
	}
	return move(Puzzle(stoi(d[0]), parseInts(d[1]), ascii));
}

int main(int argc, char* argv[])
{
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto p = load(argv[1]);
	p.shuffle();
	p.partTwo();
}
