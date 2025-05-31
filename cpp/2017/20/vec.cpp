#include <cmath>
#include <cstddef>
#include <functional>
#include <cstdint>

using namespace std;

struct Vec {
    int64_t x, y, z;

    Vec operator+(const Vec& o) const {
        return Vec{x + o.x, y + o.y, z + o.z};
    }

    bool operator==(const Vec& o) const {
        return x == o.x && y == o.y && z == o.z;
    }

    int64_t taxiDist(Vec o) const {
        return abs(o.x - x) + abs(o.y - y) + abs(o.z - z);
    }
};

struct VecHash {
    size_t operator()(const Vec& v) const {
        return hash<int>()(v.x) ^ (hash<int>()(v.y) << 1) ^ (hash<int>()(v.z) << 2);
    }
};
