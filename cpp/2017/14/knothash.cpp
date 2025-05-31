#include <vector>
#include <string>

using namespace std;

class KnotHash {
private:
	vector<int> buf;
	size_t pos = 0;
	size_t skip = 0;

	void reverse(int start, int end) {
		for (size_t i = start, j = end - 1; i < j; i++, j--) {
			swap(buf[i%256], buf[j%256]);
		}
	}

	void shuffle(const vector<int>& d) {
		for (size_t len : d) {
			reverse(pos, pos + len);
			pos = (pos + len + skip) % 256;
			skip++;
		}
	}

public:
	KnotHash() {};

	vector<int> hash(const string& input) {
		vector<int> out;
		vector<int> data;
		for (char c : input) {
			data.push_back(c);
		}
		for (int i : {17, 31, 73, 47, 23}) {
			data.push_back(i);
		}
		buf.clear();
		for (int i = 0; i < 256; i++) {
			buf.push_back(i);
		}
		pos = 0;
		skip = 0;
		for (int i = 0; i < 64; i++) {
			shuffle(data);
		}
		for (int i = 0; i < 16; i++) {
			int v = 0;
			for (int j = 0; j < 16; j++) {
				v ^= buf.at((i << 4) + j);
			}
			out.push_back(v);
		}
		return out;
	}
};

