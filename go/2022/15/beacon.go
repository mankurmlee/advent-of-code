package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"slices"
	"cmp"
)

type Vec struct { x, y int }

type Sensor struct { pos, beacon Vec }

type Puzzle []Sensor

type Range struct { u, v int }

type Ranges []Range

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func Abs(n int) int {
	if n >= 0 {
		return n
	}
	return -n
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	r := strings.NewReplacer("=", " ", ":", " ", ",", " ")
	for s.Scan() {
		data := strings.Fields(r.Replace(s.Text()))
		p = append(p, Sensor{
			Vec{ Atoi(data[3]), Atoi(data[5])},
			Vec{ Atoi(data[11]), Atoi(data[13])},
		})
	}
	check(s.Err())
	return p
}

func (v Vec) Dist(v1 Vec) int {
	return Abs(v1.x - v.x) + Abs(v1.y - v.y)
}

func (p Puzzle) GetExclusionZone(y int) (o Ranges) {
	var ranges Ranges
	for _, s := range p {
		v := s.pos
		reach := v.Dist(s.beacon)
		dist := Abs(y - v.y)
		if reach < dist {
			continue
		}
		extra := reach - dist
		ranges = append(ranges, Range{ v.x - extra, v.x + extra + 1 })
	}
	slices.SortFunc(ranges, func(a, b Range) int {
		return cmp.Compare(a.u, b.u)
	})
	n := len(ranges)
	curr := ranges[0]
	for i := 1; i < n; i++ {
		next := ranges[i]
		if curr.v >= next.u {
			curr.v = max(curr.v, next.v)
		} else {
			o = append(o, curr)
			curr = next
		}
	}
	o = append(o, curr)
	return o
}

func (rs Ranges) Contains(x int) bool {
	for _, r := range rs {
		if x >= r.u && x < r.v {
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	y := Atoi(flag.Arg(1))

	z := p.GetExclusionZone(y)
	fmt.Println(z)

	var count int
	for _, r := range z {
		count += r.v - r.u
	}

	seen := make(map[int]struct{})
	for _, s := range p {
		if s.beacon.y == y && z.Contains(s.beacon.x) {
			seen[s.beacon.x] = struct{}{}
		}
	}
	for x, _ := range seen {
		fmt.Println("Beacon at", x, ",", y)
	}
	count -= len(seen)
	fmt.Println("Part1:", count)

	max := y * 2
	for y = 0; y < max; y++ {
		z = p.GetExclusionZone(y)
		hasGap := true
		for _, r := range z {
			if r.u <= 0 && r.v > max {
				hasGap = false
				break
			}
		}
		if hasGap {
			break
		}
	}
	fmt.Println("Z:", z)
	x := z[0].v
	fmt.Println("x:", x, "y:", y)
	fmt.Println("Part 2:", x * 4000000 + y)
}
