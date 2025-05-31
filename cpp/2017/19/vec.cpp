#include <cstddef>
#include <functional>

using namespace std;

struct Vec {
    int x, y;

    Vec operator+(const Vec& o) const {
        return Vec{x + o.x, y + o.y};
    }

    bool operator==(const Vec& o) const {
        return x == o.x && y == o.y;
    }
};

struct VecHash {
    size_t operator()(const Vec& v) const {
        return hash<int>()(v.x) ^ (hash<int>()(v.y) << 1);
    }
};
