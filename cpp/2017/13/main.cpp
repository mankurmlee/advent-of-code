#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <regex>
#include <unordered_map>

using namespace std;
using Firewall = unordered_map<int, int>;

static pair<int, int> simulate(const Firewall& fw, int start) {
	int severity = 0;
	int caught = 0;
	for (const auto& [depth, range] : fw) {
		if ((start + depth) % (range + range - 2) == 0) {
			severity += depth * range;
			caught++;
		}
	}
	return {severity, caught};
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

static vector<int> parseInts(const string& s) {
	vector<int> out{};
	for (const string& s : findAllMatches("\\d+", s)) {
		out.push_back(stoi(s));
	}
	return out;
}

static Firewall load(const string& f) {
	Firewall fw{};
	for (const string& l : readFile(f)) {
		auto d = parseInts(l);
		fw.emplace(d[0], d[1]);
	}
	return fw;
}

int main(int argc, char* argv[]) {
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto fw = load(argv[1]);

	int delay = 0;
	auto p = simulate(fw, delay);
	cout << p.first << endl;

	while (p.second > 0) {
		p = simulate(fw, delay);
		delay++;
	}
	cout << (delay - 1) << endl;
}
