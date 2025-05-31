#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <regex>
#include <cstdint>
#include <unordered_map>
#include "vec.cpp"

using namespace std;

class Particle {
public:
	Vec p, v, a;
	Particle(Vec _p, Vec _v, Vec _a) :
		p(_p), v(_v), a(_a) {}
};

class World {
private:
	unordered_map<int, Particle> particles;

public:
	World(const vector<Particle>& ps) {
		int i = 0;
		for (const Particle& p : ps) {
			particles.emplace(i, p);
			i++;
		}
	}

	void smash() {
		unordered_map<Vec, vector<int>, VecHash> space;
		unordered_map<int, Particle> p1;
		for (auto [id, p] : particles) {
			p.v = p.v + p.a;
			p.p = p.p + p.v;
			space[p.p].push_back(id);
			p1.emplace(id,p);
		}
		for (const auto& [_, ids] : space) {
			if (ids.size() > 1) {
				for (int id : ids) {
					p1.erase(id);
				}
			}
		}
		particles = p1;
	}

	size_t size() const {
		return particles.size();
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
	for (const string& s : findAllMatches("\\-?\\d+", s)) {
		out.push_back(stoi(s));
	}
	return out;
}

static vector<Particle> load(const string& f) {
	vector<Particle> ps;
	for (const string& l : readFile(f)) {
		auto d = parseInts(l);
		ps.push_back(Particle(
			{d.at(0), d.at(1), d.at(2)},
			{d.at(3), d.at(4), d.at(5)},
			{d.at(6), d.at(7), d.at(8)}
		));
	}
	return ps;
}

static void partOne(const vector<Particle>& ps) {
	int best;
	Vec zero{0,0,0};
	int aLo = 1 << 30;
	int vLo = 1 << 30;
	int i = 0;
	for (const Particle& p : ps) {
		int a = p.a.taxiDist(zero);
		if (a < aLo) {
			aLo = a;
			vLo = p.v.taxiDist(zero);
			best = i;
		} else if (a == aLo) {
			int v =  p.v.taxiDist(zero);
			if (v < vLo) {
				vLo = v;
				best = i;
			}
		}
		i++;
	}
	cout << best << endl;
}

void partTwo(const vector<Particle>& particles) {
	World w(particles);
	for (int i = 0; i < 100; i++) {
		w.smash();
	}
	cout << w.size() << endl;
}

int main(int argc, char* argv[]) {
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto ps = load(argv[1]);
	partOne(ps);
	partTwo(ps);
}
