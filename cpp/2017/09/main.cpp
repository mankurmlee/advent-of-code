#include <vector>
#include <iostream>
#include <fstream>
#include <string>

using namespace std;

static void parse(const string& s) {
	int gc = 0;
	int score = 0;
	int depth = 0;
	bool garbage = false;
	size_t n = s.size();
	size_t i = 0;
	while (i < n) {
		char v = s[i];
		if (garbage) {
			if (v == '!') {
				i++;
			} else if (v == '>') {
				garbage = false;
			} else {
				gc++;
			}
			i++;
			continue;
		}
		if (v == '{') {
			depth++;
		} else if (v == '}') {
			score += depth;
			depth--;
		} else if (v == '<') {
			garbage = true;
		}
		i++;
	}
	cout << score << endl;
	cout << gc << endl;
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

int main(int argc, char* argv[])
{
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}

	for (const string& s : readFile(argv[1])) {
		parse(s);
	}
}
