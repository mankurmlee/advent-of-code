#include <string>
#include <vector>
#include <cstdint>
#include <deque>
#include <unordered_map>

using namespace std;

template <typename T>
using Matrix = vector<vector<T>>;

class Thread {
private:
	bool duet;
	bool yield;
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

	void jnz(vector<string> stmt) {
		if (stov(stmt.at(1)) != 0) {
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
		} else if (cmd == "sub") {
			regs[stmt.at(1)] -= stov(stmt.at(2));
		} else if (cmd == "mul") {
			regs[stmt.at(1)] *= stov(stmt.at(2));
			regs["mul"]++;
		} else if (cmd == "mod") {
			regs[stmt.at(1)] %= stov(stmt.at(2));
		} else if (cmd == "rcv") {
			rcv(stmt);
			if (yield) return;
		} else if (cmd == "jgz") {
			jgz(stmt);
		} else if (cmd == "jnz") {
			jnz(stmt);
		}
		regs["ip"]++;
	}

public:
	Thread(Matrix<string> p, int id, bool d) :
		prog(p),
		duet(d),
		yield(false)
	{
		regs["p"] = id;
	}

	void run() {
		auto n = prog.size();
		yield = false;
		while (!yield && regs["ip"] < n) {
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
};
