#include <vector>
#include <iostream>
#include <fstream>
#include <regex>
#include <unordered_map>

using namespace std;

struct Component {
	int port[2];

	inline bool hasPort(int p) const {
		return port[0] == p || port[1] == p;
	}

	inline int strength() const {
		return port[0] + port[1];
	}
};

struct Config {
	vector<Component> box;
};

struct VectorHash {
    size_t operator()(const vector<int>& vec) const {
        size_t hash = 0;
        for (int num : vec) {
            hash ^= std::hash<int>()(num) + 0x9e3779b9 + (hash << 6) + (hash >> 2);
        }
        return hash;
    }
};

struct VectorEqual {
    bool operator()(const vector<int>& vec1, const vector<int>& vec2) const {
        return vec1 == vec2;
    }
};

struct Bridge {
	int length, strength;
};

class BridgeBuilder {
private:
	vector<Component> components;
	unordered_map<vector<int>,int,VectorHash,VectorEqual> strmemo;
	unordered_map<vector<int>,Bridge,VectorHash,VectorEqual> lenmemo;

	vector<int> removeValue(const vector<int>& input, int value) const {
		vector<int> out;
		for (int v : input) {
			if (v != value) {
				out.push_back(v);
			}
		}
		return move(out);
	}

	int findStrong(int port, const vector<int>& bag) {
		vector<int> key = bag;
		key.push_back(port);
		auto it = strmemo.find(key);
		if (it != strmemo.end()) {
			return it->second;
		}
		int best = 0;
		for (int v : bag) {
			auto c = components.at(v);
			if (!c.hasPort(port)) {
				continue;
			}
			int str = c.strength();
			if (bag.size() > 1) {
				int port1 = port == c.port[0] ? c.port[1] : c.port[0];
				str += findStrong(port1, move(removeValue(bag, v)));
			}
			if (str > best) {
				best = str;
			}
		}
		strmemo.emplace(key, best);
		return best;
	}

	Bridge findLongest(int port, const vector<int>& bag) {
		vector<int> key = bag;
		key.push_back(port);
		auto it = lenmemo.find(key);
		if (it != lenmemo.end()) {
			return it->second;
		}
		Bridge best{0, 0};
		for (int v : bag) {
			auto c = components.at(v);
			if (!c.hasPort(port)) {
				continue;
			}
			Bridge b{1, c.strength()};
			if (bag.size() > 1) {
				int port1 = port == c.port[0] ? c.port[1] : c.port[0];
				Bridge child = findLongest(port1, move(removeValue(bag, v)));
				b.length += child.length;
				b.strength += child.strength;
			}
			if (b.length > best.length || b.length == best.length && b.strength > best.strength) {
				best = b;
			}
		}
		lenmemo.emplace(key, best);
		return best;
	}

public:
	int strongest(const Config& cfg) {
		components = cfg.box;
		strmemo.clear();
		size_t n = components.size();
		vector<int> bag;
		for (size_t i = 0; i < n; i++) {
			bag.push_back(i);
		}
		return findStrong(0, bag);
	}

	int longest(const Config& cfg) {
		components = cfg.box;
		lenmemo.clear();
		size_t n = components.size();
		vector<int> bag;
		for (size_t i = 0; i < n; i++) {
			bag.push_back(i);
		}
		return findLongest(0, bag).strength;
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

static Config load(const string& f) {
	Config cfg;
	for (const string& l : readFile(f)) {
		auto nums = parseInts(l);
		cfg.box.push_back({{nums[0], nums[1]}});
	}
	return cfg;
}

int main(int argc, char* argv[]) {
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto cfg = load(argv[1]);

	BridgeBuilder b;
	cout << b.strongest(cfg) << endl;
	cout << b.longest(cfg) << endl;
}
