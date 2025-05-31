#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <regex>
#include <unordered_map>

using namespace std;

template <typename T>
struct hash<std::vector<T>> {
    std::size_t operator()(const std::vector<T>& vec) const {
        std::size_t h = 0;
        for (const auto& elem : vec) {
            h ^= std::hash<T>{}(elem) + 0x9e3779b9 + (h << 6) + (h >> 2);
        }
        return h;
    }
};

static void countCycles(const vector<int>& data) {
    int count = 0;
    auto d = data;
    auto n = d.size();
    unordered_map<vector<int>, int> seen{};
    while (seen.count(d) == 0) {
        seen.emplace(d, count);
        int li = 0;
        int lv = 0;
        for (int i = 0; i < n; i++) {
            if (d[i] > lv) {
                lv = d[i];
                li = i;
            }
        }
        d[li] = 0;
        for (int i = 0; i < lv; i++) {
            d[(li + 1 + i) % n]++;
        }
        count++;
    }
    cout << count << endl;
    cout << (count - seen[d]) << endl;
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
    for (const auto& e : findAllMatches("\\-?\\d+", s)) {
        out.push_back(stoi(e));
    }
    return out;
}

static vector<int> load(const string& f) {
    auto d = readFile(f);
    return parseInts(d[0]);
}

int main(int argc, char* argv[])
{
    if (argc != 2) {
        cerr << "Usage: " << argv[0] << " <filename>" << endl;
        return 1;
    }
    auto data = load(argv[1]);
    countCycles(data);
}
