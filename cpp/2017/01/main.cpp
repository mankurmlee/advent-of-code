#include <vector>
#include <iostream>
#include <fstream>
#include <string>

using namespace std;

static vector<string> readFile(string filename) {
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

static int captcha(string s, bool half) {
    int sum = 0;
    size_t n = s.length();
    size_t offset = half ? n >> 1: 1;
    for (size_t i = 0; i < n; i++) {
        if (s[i] == s[(i + offset)%n]) {
            sum += s[i] - '0';
        }
    }
    return sum;
}

static void doAll(vector<string> data, bool half) {
    for (const auto& l : data) {
        cout << captcha(l, half) << endl;
    }
}

int main(int argc, char* argv[])
{
    if (argc != 2) {
        cerr << "Usage: " << argv[0] << " <filename>" << endl;
        return 1;
    }
    auto data = readFile(argv[1]);
    doAll(data, false);
    doAll(data, true);
}
