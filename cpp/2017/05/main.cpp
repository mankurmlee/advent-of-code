#include <vector>
#include <iostream>
#include <fstream>
#include <string>

using namespace std;

static void countSteps(const vector<int>& data, bool part2) {
    auto maze = data;
    int i = 0;
    int o = 0;
    int offset = 0;
    int steps = 0;
    auto n = data.size();
    while (i >= 0 && i < n) {
        o = i;
        i += maze[o];
        offset = part2 && offset >= 3 ? offset - 1 : offset + 1;
        maze[o] = offset;
        steps++;
    }
    cout << steps << endl;
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

static vector<int> load(const string& f) {
    vector<int> out;
    for (const string& l : readFile(f)) {
        out.push_back(stoi(l));
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
    countSteps(data, false);
    countSteps(data, true);
}
