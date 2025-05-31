#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <regex>
#include <unordered_map>
#include <cstdint>
#include <deque>

using namespace std;

template <typename T>
using Matrix = vector<vector<T>>;

class Thread {
private:
	bool duet;
	bool yield;
	int sent;
	unordered_map<string,int64_t> regs;
	Matrix<string> prog;
	deque<int64_t> qin;
	vector<int64_t> qout;

	int stov(string s) {
		try {
			return stoi(s);
		} catch (const exception&) {
			return regs[s];
		}
	}

	void snd(vector<string> stmt) {
		int v = stov(stmt.at(1));
		if (duet) {
			qout.push_back(v);
			sent++;
		} else {
			regs["snd"] = v;
		}
	}

	void rcv(vector<string> stmt) {
		if (duet) {
			if (qin.size() == 0) {
				yield = true;
			} else {
				regs[stmt.at(1)] = qin.front();
				qin.pop_front();
			}
		} else if (stov(stmt.at(1))) {
			regs["rcv"] = regs["snd"];
			yield = true;
		}
	}

	void jgz(vector<string> stmt) {
		if (stov(stmt.at(1)) > 0) {
			regs["ip"] += stov(stmt.at(2)) - 1;
		}
	}

	void exec(vector<string> stmt) {
		string cmd = stmt.at(0);
		if (cmd == "snd") {
			snd(stmt);
		} else if (cmd == "set") {
			regs[stmt.at(1)] = stov(stmt.at(2));
		} else if (cmd == "add") {
			regs[stmt.at(1)] += stov(stmt.at(2));
		} else if (cmd == "mul") {
			regs[stmt.at(1)] *= stov(stmt.at(2));
		} else if (cmd == "mod") {
			regs[stmt.at(1)] %= stov(stmt.at(2));
		} else if (cmd == "rcv") {
			rcv(stmt);
			if (yield) return;
		} else if (cmd == "jgz") {
			jgz(stmt);
		}
		regs["ip"]++;
	}

public:
	Thread(Matrix<string> p, int id, bool d) :
		prog(p),
		duet(d),
		yield(false),
		sent(0) {
		regs["p"] = id;
	}

	void run() {
		auto n = prog.size();
		yield = false;
		while (!yield) {
			exec(prog.at(regs["ip"]));
		}
	}

	int64_t get(string s) {
		if (regs.count(s)) {
			return regs[s];
		}
		return 0;
	}

	void send(vector<int64_t> msg) {
		qin.insert(qin.end(), msg.begin(), msg.end());
	}

	vector<int64_t> recv() {
		auto out = move(qout);
		qout.clear();
		return move(out);
	}

	int getSent() {
		return sent;
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
	Matrix<string> prog;
	for (const string& l : readFile(f)) {
		prog.push_back(findAllMatches("\\S+", l));
	}
	return prog;
}

static void partOne(Matrix<string> prog) {
	Thread c(prog, 0, false);
	c.run();
	cout << c.get("rcv") << endl;
}

static void partTwo(Matrix<string> prog) {
	Thread c0(prog, 0, true);
	Thread c1(prog, 1, true);
	int nMsg = 1;
	while (nMsg > 0) {
		nMsg = 0;
		c0.run();
		auto msg = c0.recv();
		nMsg += msg.size();
		c1.send(msg);
		c1.run();
		msg = c1.recv();
		nMsg += msg.size();
		c0.send(msg);
	}
	cout << c1.getSent() << endl;
}

int main(int argc, char* argv[]) {
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto prog = load(argv[1]);
	partOne(prog);
	partTwo(prog);
}
