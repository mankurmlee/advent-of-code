#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <unordered_map>
#include <algorithm>
#include <regex>

using namespace std;

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

class Surface {
private:
	int width;
	vector<char> buffer;

public:
	Surface(int _size) : width(_size) {
		buffer.resize(width*width);
	}

	Surface(int _size, const vector<char>& _buffer) :
		width(_size), buffer(_buffer) {}

	Surface(const string& s) {
		auto rows = findAllMatches("[.#]+", s);
		width = rows.size();
		buffer.reserve(width*width);
		for (const string& r : rows) {
			buffer.insert(buffer.end(), r.begin(), r.end());
		}
	}

	void flip() {
		for (int i = 0, j = width - 1; i < j; i++, j--) {
			int src = i * width;
			int dest = j * width;
			swap_ranges(
				buffer.begin() + src,
				buffer.begin() + src + width,
				buffer.begin() + dest
			);
		}
	}

	void rotate() {
		vector<char> out;
		int n = width * width;
		out.reserve(n);
		int u = 0;
		int v = width - 1;
		for (int i = 0; i < n; i++) {
			out.push_back(buffer.at(v * width + u));
			v--;
			if (v < 0) {
				v += width;
				u++;
			}
		}
		buffer = out;
	}

	int sum() const {
		int v = 0;
		for (char c : buffer) {
			v = (v << 1) | (c == '#' ? 1 : 0);
		}
		return v;
	}

	// Could speed this up if we cache this
	int id() const {
		int id = sum();
		Surface s1 = *this;
		for (int j = 0; j < 2; j++) {
			for (int i = 0; i < 4; i++) {
				s1.rotate();
				int sum1 = s1.sum();
				if (sum1 < id) {
					id = sum1;
				}
			}
			s1.flip();
		}
		return id;
	}

	inline int getWidth() const {
		return width;
	}

	Surface clip(int x, int y, int w) const {
		vector<char> b;
		b.reserve(w*w);
		int u = y * width + x;
		for (int i = 0; i < w; i++) {
			b.insert(b.end(), buffer.begin() + u, buffer.begin() + u + w);
			u += width;
		}
		return move(Surface{w, b});
	}

	void blit(int x, int y, const Surface& bmp) {
		auto n = bmp.buffer.size();
		int j = y * width + x;
		for (int i = 0; i < n; i += bmp.width) {
			copy_n(bmp.buffer.begin() + i, bmp.width, buffer.begin() + j);
			j += width;
		}
	}

	int count() const {
		int on = 0;
		for (char c : buffer) {
			if (c == '#') {
				on++;
			}
		}
		return on;
	}
};

class Volume {
private:
	int width;
	unordered_map<int,Surface> surfs;

public:
	Volume(int w) : width(w) {}

	void add(const Surface& k, const Surface& v) {
		int w = k.getWidth();
		if (width != w) {
			cout << "Unexpected key size: " << w << endl;
			throw runtime_error("Unexpected key size!");
		}
		surfs.emplace(k.id(), v);
	}

	Surface enhance(const Surface& grid) const {
		int n = grid.getWidth() / width;
		Surface out(n * (width + 1));
		for (int v = 0; v < n; v++) {
			for (int u = 0; u < n; u++) {
				Surface c = grid.clip(u * width, v * width, width);
				const Surface& p = surfs.at(c.id());
				out.blit(u * p.getWidth(), v * p.getWidth(), p);
			}
		}
		return move(out);
	}
};

class Rules {
private:
	unordered_map<int,Volume> volumes;

public:
	Volume& getVolume(int width) {
		auto it = volumes.find(width);
		if (it != volumes.end()) {
			return it->second;
		}
		volumes.emplace(width, Volume(width));
		return volumes.at(width);
	}

	Surface enhance(const Surface& grid) const {
		int size = grid.getWidth();
		int n = 0;
		if (size % 2 == 0) {
			n = 2;
		} else if (size % 3 == 0) {
			n = 3;
		}
		const Volume& v = volumes.at(n);
		return move(v.enhance(grid));
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

static Rules load(const string& f) {
	Rules out;
	for (const string& l : readFile(f)) {
		vector<string> m = findAllMatches("\\S+", l);
		Surface key(m[0]);
		Volume& v = out.getVolume(key.getWidth());
		v.add(key, Surface(m[2]));
	}
	return move(out);
}

int main(int argc, char* argv[]) {
	if (argc != 2) {
		cerr << "Usage: " << argv[0] << " <filename>" << endl;
		return 1;
	}
	auto r = load(argv[1]);

	Surface g(".#./..#/###");
	for (int i = 0; i < 18; i++) {
		g = r.enhance(g);
		if (i == 4) {
			cout << g.count() << endl;
		}
	}
	cout << g.count() << endl;
}
