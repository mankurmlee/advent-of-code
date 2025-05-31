#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <regex>

using namespace std;

template <typename T>
using Matrix = vector<vector<T>>;

static vector<string> readFile(const string& filename) {
    vector<string> lines{};
    ifstream file(filename);
    if (!file) {
        cerr << "Error: Could not open file " << filename << endl;
        terminate();
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
    for (const string& s : findAllMatches("\\d+", s)) {
        out.push_back(stoi(s));
    }
    return out;
}

static Matrix<int> load(const string& f) {
    Matrix<int> out{};
    for (const string& l : readFile(f)) {
        out.push_back(parseInts(l));
    }
    return out;
}

static int biggestDiff(const vector<int>& row) {
    int lo = row[0];
    int hi = lo;
    for (const auto& v : row) {
        if (v < lo) {
            lo = v;
        }
        if (v > hi) {
            hi = v;
        }
    }
    return hi - lo;
}

static int findQuotient(const vector<int>& row) {
    for (const auto& top : row) {
        for (const auto& bot : row) {
            if (top != bot && top % bot == 0) {
                return top / bot;
            }
        }
    }
    throw runtime_error("Could not find div!");
}

static void checksum(int (*func)(const vector<int>& row), const Matrix<int>& m) {
    int sum = 0;
    for (const auto& r : m) {
        sum += func(r);
    }
    cout << sum << endl;
}

int main(int argc, char* argv[])
{
    if (argc != 2) {
        cerr << "Usage: " << argv[0] << " <filename>" << endl;
        return 1;
    }
    auto data = load(argv[1]);
    checksum(biggestDiff, data);
    checksum(findQuotient, data);
}
