#include <vector>
#include <iostream>
#include <fstream>
#include <string>

using namespace std;

class Spinlock {
private:
	int step;
	int pos;
	vector<int> buf;

public:
	Spinlock(int s) : step(s), pos(0), buf({0}) {}

	void insert() {
		auto n = buf.size();
		pos = (pos + step) % n;
		if (pos + 1 == n) {
			buf.push_back(n);
		} else {
			buf.insert(buf.begin() + pos + 1, n);
		}
		pos = (pos + 1) % buf.size();
	}

	void partOne(int n) {
		for (int i = 0; i < n; i++) {
			insert();
		}
		cout << buf.at((pos + 1) % buf.size()) << endl;
	}

	void partTwo(int n) {
		int out = 0;
		pos = 0;
		int len = 1;
		while (len <= n) {
			if (pos == 1) {
				out = len - 1;
			}
			int ins = (len - pos + step - 1) / step;
			pos = ((pos + ins * step) % len) + 1;
			len += ins;
		}
		cout << out << endl;
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

static int load(const string& f) {
	auto d = readFile(f);
	return stoi(d[0]);
}

int main(int argc, char* argv[]) {
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	int step = load(argv[1]);
	Spinlock s(step);
	s.partOne(2017);
	s.partTwo(50000000);
}
