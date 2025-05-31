package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var DIR = [4]Vec{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

func check(err error) {
	if err != nil {
		panic(err)
	}
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
	var m Maze
	m.Load(flag.Arg(0))
	m.Solve()
}

type Vec struct {
	x, y int
}

func (v Vec) Equals(o Vec) bool {
	return v.x == o.x && v.y == o.y
}

func (v Vec) Add(o Vec) Vec {
	return Vec{v.x + o.x, v.y + o.y}
}

type Deer struct {
	Vec
	dir int
}

type DeerState struct {
	Deer
	score int
	last  *DeerState
}

type Maze struct {
	w, h       int
	start, end Vec
	data       []rune
}

func (m Maze) Solve() {
	lowest := -1
	var paths []DeerState
	var s DeerState
	s.Vec = m.start
	best := make(map[Deer]int)
	best[s.Deer] = 0
	walls := m.Walls()
	pq := NewPriorityQueue[DeerState]()
	pq.Enqueue(s, s.score)
	for pq.Len() > 0 {
		v := pq.Dequeue()
		if v.Vec.Equals(m.end) {
			if v.score == lowest {
				paths = append(paths, v)
			}
			if v.score < lowest || lowest < 0 {
				lowest = v.score
				paths = []DeerState{v}
			}
			continue
		}
		v1 := DeerState{
			Deer{v.Add(DIR[v.dir]), v.dir},
			v.score + 1,
			&v,
		}
		if !walls[v1.Vec] {
			bs, ok := best[v1.Deer]
			if !ok || bs >= v1.score {
				best[v1.Deer] = v1.score
				pq.Enqueue(v1, v1.score)
			}
		}
		for _, d1 := range []int{1, 3} {
			v2 := DeerState{
				Deer{v.Vec, (v.dir + d1) % 4},
				v.score + 1000,
				&v,
			}
			bs, ok := best[v2.Deer]
			if !ok || bs >= v2.score {
				best[v2.Deer] = v2.score
				pq.Enqueue(v2, v2.score)
			}
		}
	}
	fmt.Println("Part 1:", lowest)
	uniq := make(map[Vec]bool)
	for _, p := range paths {
		v := &p
		for v.last != nil {
			uniq[v.Vec] = true
			v = v.last
		}
		uniq[v.Vec] = true
	}
	fmt.Println("Part 2:", len(uniq))
}

func (m Maze) Walls() map[Vec]bool {
	w := make(map[Vec]bool)
	for y := range m.h {
		for x := range m.w {
			if m.data[y*m.w+x] == '#' {
				w[Vec{x, y}] = true
			}
		}
	}
	return w
}

func (m *Maze) Load(filename string) {
	d := readFile(filename)
	m.h = len(d)
	m.w = len(d[0])
	for y, line := range d {
		for x, v := range line {
			m.data = append(m.data, v)
			if v == 'S' {
				m.start = Vec{x, y}
			}
			if v == 'E' {
				m.end = Vec{x, y}
			}
		}
	}
}
