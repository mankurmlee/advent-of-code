package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Vec [3]int

type Cuboid struct {
	lo      Vec
	hi      Vec
	turnOn  bool
}

func (p Cuboid) Intersects(q Cuboid) bool {
	if  p.hi[0] < q.lo[0] || q.hi[0] < p.lo[0] ||
		p.hi[1] < q.lo[1] || q.hi[1] < p.lo[1] ||
		p.hi[2] < q.lo[2] || q.hi[2] < p.lo[2] {
		return false
	}
	return true
}

func (c Cuboid) Volume() int {
	vol := c.hi[0] - c.lo[0] + 1
	vol *= c.hi[1] - c.lo[1] + 1
	vol *= c.hi[2] - c.lo[2] + 1
	return vol
}

func (c Cuboid) Diff(other Cuboid) []Cuboid {
	if !c.Intersects(other) {
		return []Cuboid{c}
	}

	dflt := c
	out := make([]Cuboid, 0, 6)
	for i := range 3 {
		if c.lo[i] < other.lo[i] {
			new := dflt
			new.hi[i] = other.lo[i] - 1
			out = append(out, new)
		}
		if c.hi[i] > other.hi[i] {
			new := dflt
			new.lo[i] = other.hi[i] + 1
			out = append(out, new)
		}

		dflt.lo[i] = max(dflt.lo[i], other.lo[i])
		dflt.hi[i] = min(dflt.hi[i], other.hi[i])
	}

	return out
}

type Reactor []Cuboid

func (r Reactor) Count() (vol int) {
	for _, c := range r {
		vol += c.Volume()
	}
	return vol
}

func (r Reactor) RunStep(c Cuboid) (out Reactor) {
	for _, other := range r {
		out = append(out, other.Diff(c)...)
	}
	if c.turnOn {
		out = append(out, c)
	}
	return out
}

func (r Reactor) ClipTo(clip Cuboid) (out Reactor) {
	for _, c := range r {
		if !c.Intersects(clip) {continue}
		for i := range 3 {
			c.lo[i] = max(c.lo[i], clip.lo[i])
			c.hi[i] = min(c.hi[i], clip.hi[i])
		}
		out = append(out, c)
	}
	return out
}

func atoi(a string) int {
	i, err := strconv.Atoi(a)
	check(err)
	return i
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load(filename string) (p Reactor) {
	r := strings.NewReplacer(".", " ", "=", " ", ",", " ")
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		d := strings.Fields(r.Replace(s.Text()))
		p = append(p, Cuboid{
			Vec{ atoi(d[2]), atoi(d[5]), atoi(d[8]) },
			Vec{ atoi(d[3]), atoi(d[6]), atoi(d[9]) },
			d[0] == "on",
		})
	}
	check(s.Err())
	return p
}

func main() {
	flag.Parse()

	var r Reactor
	for _, c := range load(flag.Arg(0)) {
		r = r.RunStep(c)
	}
	c := r.ClipTo(Cuboid{Vec{-50, -50, -50}, Vec{50, 50, 50}, false})

	fmt.Println("Cuboids:", len(r))
	fmt.Println("Part 1:", c.Count())
	fmt.Println("Part 2:", r.Count())
}
