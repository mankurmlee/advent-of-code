#include <vector>
#include <iostream>
#include <fstream>
#include <regex>
#include <unordered_map>

using namespace std;

struct State {
	char id;
	int value;

	bool operator==(const State& other) const {
        return id == other.id && value == other.value;
    }

	struct Hash {
        std::size_t operator()(const State& s) const {
            return std::hash<char>()(s.id) ^ (std::hash<int>()(s.value) << 1);
        }
    };
};

struct Action {
	int value;
	int inc;
	char next;
};

using Rules = unordered_map<State,Action,State::Hash>;

struct Config {
	char start;
	int steps;
	Rules rules;
};

class Tape {
private:
	int pos;
	unordered_map<int,int> data;

public:
	int checksum() const {
		int count = 0;
		for (const auto& [_, v] : data) {
			if (v == 1) {
				count++;
			}
		}
		return count;
	}

	int read() const {
		auto it = data.find(pos);
		if (it != data.end()) {
			return it->second;
		}
		return 0;
	}

	inline void write(int value) {
		data[pos] = value;
	}

	inline void move(int value) {
		pos += value;
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

static vector<vector<string>> readChunkedFile(const string& filename) {
	vector<vector<string>> chunks{};
	ifstream file(filename);
	if (!file) {
		cerr << "Error: Could not open file " << filename << endl;
		throw runtime_error("File not found!");
	}
	vector<string> chunk;
	string line;
	while (getline(file, line)) {
		if (line != "") {
			chunk.push_back(line);
			continue;
		}
		if (chunk.size() == 0) {
			continue;
		}
		chunks.push_back(chunk);
		chunk.clear();
	}
	file.close();
	if (chunk.size() > 0) {
		chunks.push_back(chunk);
	}
	return chunks;
}

static Config load(const string& f) {
	Config cfg;
	auto chunks = readChunkedFile(f);
	auto header = chunks.at(0);
	cfg.start = findAllMatches("\\w+", header[0]).back()[0];
	cfg.steps = parseInts(header[1])[0];
	auto n = chunks.size();
	for (int i = 1; i < n; i++) {
		auto d = chunks.at(i);
		char id = findAllMatches("\\w+", d[0]).back()[0];
		for (int o : {0, 4}) {
			int value = parseInts(d[1 + o])[0];
			int write = parseInts(d[2 + o])[0];
			int inc = findAllMatches("\\w+", d[3 + o]).back() == "left" ?  -1 : 1;
			char next = findAllMatches("\\w+", d[4 + o]).back()[0];
			cfg.rules.emplace(State{id, value}, Action{write, inc, next});
		}
	}
	return cfg;
}

static void run(const Config& cfg) {
	Tape tape;
	char id = cfg.start;
	for (size_t i = 0; i < cfg.steps; i++) {
		int value = tape.read();
		const Action& a = cfg.rules.at(State{id, value});
		if (value != a.value) {
			tape.write(a.value);
		}
		tape.move(a.inc);
		id = a.next;
	}
	cout << tape.checksum() << endl;
}

int main(int argc, char* argv[]) {
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto cfg = load(argv[1]);
	run(cfg);
}
