#include <vector>
#include <iostream>
#include <fstream>
#include <string>
#include <unordered_set>
#include <regex>

using namespace std;
using Passphrase = vector<string>;

static bool partTwo(Passphrase pass) {
    unordered_set<string> seen;
    for (const auto& w : pass) {
        string s = w;
        sort(s.begin(), s.end());
        if (seen.find(s) != seen.end()) {
            return false;
        }
        seen.insert(s);
    }
    return true;
}

static bool partOne(Passphrase pass) {
    unordered_set<string> seen;
    for (const auto& w : pass) {
        if (seen.find(w) != seen.end()) {
            return false;
        }
        seen.insert(w);
    }
    return true;
}

static void checksum(bool (*func)(Passphrase), vector<Passphrase> data) {
    int sum = 0;
    for (const auto& p : data) {
        if (func(p)) {
            sum++;
        }
    }
    cout << sum << endl;
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

static vector<string> parseWords(string s) {
    vector<string> out;
    for (const string& s : findAllMatches("\\w+", s)) {
        out.push_back(s);
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

static vector<Passphrase> load(const string& f) {
    vector<Passphrase> out;
    for (const string& l : readFile(f)) {
        out.push_back(parseWords(l));
    }
    return out;
}

int main(int argc, char* argv[])
{
    if (argc != 2) {
        cerr << "Usage: " << argv[0] << " <filename>" << endl;
        return 1;
    }
    auto data = load(argv[1]);
    checksum(partOne, data);
    checksum(partTwo, data);
}
