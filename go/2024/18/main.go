package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var DIR = [4]Vec{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func atoi(a string) int {
	i, err := strconv.Atoi(a)
	check(err)
	return i
}

func findInts(s string) (nums []int) {
	r := regexp.MustCompile(`\-?\d+`)
	for _, v := range r.FindAllString(s, -1) {
		nums = append(nums, atoi(v))
	}
	return nums
}

func readFile(filename string) (lines []string) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	check(s.Err())
	return lines
}

func main() {
	flag.Parse()
	g := load(flag.Arg(0))

	walls := g.Walls()

	c := g.FindExit(walls)
	fmt.Println(c)

	var v Vec
	for i := g.bytes; c > 0; i++ {
		v = g.data[i]
		walls[v] = true
		c = g.FindExit(walls)
	}
	fmt.Println(v)
}

func load(filename string) Grid {
	d := readFile(filename)
	dim := findInts(d[0])
	var data []Vec
	for _, l := range d[1:] {
		v := findInts(l)
		data = append(data, Vec{v[0], v[1]})
	}
	return Grid{dim[0], dim[1], dim[2], data}
}

type Vec struct {
	x, y int
}

func (v Vec) Add(o Vec) Vec {
	return Vec{v.x + o.x, v.y + o.y}
}

func (v Vec) Equals(o Vec) bool {
	return v.x == o.x && v.y == o.y
}

type State struct {
	Vec
	steps int
}

type Grid struct {
	w, h, bytes int
	data        []Vec
}

func (g Grid) FindExit(walls map[Vec]bool) int {
	dest := Vec{g.w - 1, g.h - 1}
	been := make(map[Vec]int)
	been[Vec{0, 0}] = 0
	q := []State{{}}
	for len(q) > 0 {
		s := q[0]
		q = q[1:]
		if s.Equals(dest) {
			return s.steps
		}
		q = append(q, g.Expand(s, walls, been)...)
	}
	return 0
}

func (g Grid) Expand(s State, walls map[Vec]bool, been map[Vec]int) (out []State) {
	for _, d := range DIR {
		s1 := State{s.Vec.Add(d), s.steps + 1}
		if s1.x < 0 || s1.y < 0 || s1.x >= g.w || s1.y >= g.h {
			continue
		}
		if walls[s1.Vec] {
			continue
		}
		b, ok := been[s1.Vec]
		if ok && b <= s1.steps {
			continue
		}
		been[s1.Vec] = s.steps
		out = append(out, s1)
	}
	return out
}

func (g Grid) Walls() map[Vec]bool {
	out := make(map[Vec]bool)
	for _, v := range g.data[:g.bytes] {
		out[v] = true
	}
	return out
}
