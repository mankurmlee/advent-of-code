#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <regex>
#include <unordered_map>
#include <unordered_set>

using namespace std;

using Data = unordered_map<int, vector<int>>;

static unordered_set<int> getGroup(const Data& pipes, int start) {
	unordered_set<int> seen{start};
	vector<int> q{start};
	while (q.size() > 0) {
		int s = q.at(0);
		q.erase(q.begin());
		for (int conn : pipes.at(s)) {
			if (seen.count(conn)) {
				continue;
			}
			seen.insert(conn);
			q.push_back(conn);
		}
	}
	return seen;
}

static void countGroups(const Data& pipes) {
	unordered_set<int> seen;
	int id = 0;
	bool found = true;
	int num = 1;
	while (found) {
		unordered_set<int> seen1 = getGroup(pipes, id);
		if (id == 0) {
			cout << seen1.size() << endl;
		}
		seen.insert(seen1.begin(), seen1.end());
		found = false;
		for (const auto& [k, _] : pipes) {
			if (seen.count(k) == 0) {
				found = true;
				num++;
				id = k;
				break;
			}
		}
	}
	cout << num << endl;
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

static Data load(const string& f) {
	Data pipes{};
	for (const string& l : readFile(f)) {
		auto d = parseInts(l);
		vector<int> conns(d.begin() + 1, d.end());
		pipes[d[0]] = conns;
	}
	return pipes;
}

int main(int argc, char* argv[]) {
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto pipes = load(argv[1]);
	countGroups(pipes);
}
