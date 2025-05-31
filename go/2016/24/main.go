package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Vec struct {
	x, y int
}

func (v Vec) Add(o Vec) Vec {
	return Vec{v.x + o.x, v.y + o.y}
}

func (v Vec) Equals(o Vec) bool {
	return v.x == o.x && v.y == o.y
}

var DIRS = [4]Vec{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

type State struct {
	pos  Vec
	cost int
}

type Memo struct {
	ret0  bool
	costs map[[2]rune]int
	cache map[string]int
}

func (m Memo) GetCost(key string) int {
	if cost, ok := m.cache[key]; ok {
		return cost
	}
	n := []rune(key)
	cost := -1
	for k, c := range m.costs {
		p, q := k[0], k[1]
		if p != n[0] || !slices.Contains(n[1:], q) {
			continue
		}
		key1 := makeKey(q, n[1:])
		c1 := c
		if len(key1) > 1 {
			c1 += m.GetCost(key1)
		} else if m.ret0 {
			c1 += m.costs[[2]rune{rune(key1[0]), '0'}]
		}
		if cost < 0 || c1 < cost {
			cost = c1
		}
	}
	m.cache[key] = cost
	return cost
}

type Puzzle struct {
	walls   map[Vec]bool
	nodes   map[Vec]rune
	nodePos map[rune]Vec
	costs   map[[2]rune]int
}

type RuneCost struct {
	node rune
	cost int
}

func (p Puzzle) Home() {
	seen := make(map[rune]int)
	seen['0'] = 0
	s0 := RuneCost{'0', 0}
	pq := NewPriorityQueue[RuneCost]()
	pq.Enqueue(s0, 0)
	for pq.Len() > 0 {
		s := pq.Dequeue()
		for k, c := range p.costs {
			p, q := k[0], k[1]
			if s.node != p {
				continue
			}
			c1 := s.cost + c
			best, ok := seen[q]
			if ok && best <= c1 {
				continue
			}
			seen[q] = c1
			pq.Enqueue(RuneCost{q, c1}, c1)
		}
	}
	for r, c := range seen {
		if r == '0' {
			continue
		}
		k := [2]rune{r, '0'}
		if _, ok := p.costs[k]; ok {
			continue
		}
		p.costs[k] = c
	}
}

func (p Puzzle) Solve(ret0 bool) {
	var nodes []rune
	for r := range p.nodePos {
		nodes = append(nodes, r)
	}
	memo := p.NewMemo(ret0)
	fmt.Println(memo.GetCost(makeKey('0', nodes)))
}

func (p Puzzle) NewMemo(ret0 bool) Memo {
	m := Memo{
		ret0,
		make(map[[2]rune]int),
		make(map[string]int),
	}
	for k, v := range p.costs {
		m.costs[k] = v
	}
	return m
}

func (p Puzzle) Verts(r0 rune) {
	v0 := p.nodePos[r0]
	seen := make(map[Vec]bool)
	seen[v0] = true
	q := []State{{v0, 0}}
	for len(q) > 0 {
		s := q[0]
		q = q[1:]
		r, ok := p.nodes[s.pos]
		if ok && !s.pos.Equals(v0) {
			p.costs[[2]rune{r0, r}] = s.cost
			p.costs[[2]rune{r, r0}] = s.cost
			continue
		}
		for _, d := range DIRS {
			v1 := s.pos.Add(d)
			if p.walls[v1] || seen[v1] {
				continue
			}
			seen[v1] = true
			q = append(q, State{v1, s.cost + 1})
		}
	}
}

func makeKey(r0 rune, nodes []rune) string {
	var sb strings.Builder
	sb.WriteRune(r0)
	slices.Sort(nodes)
	for _, r := range nodes {
		if r != r0 {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

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

func load(s string) Puzzle {
	walls := make(map[Vec]bool)
	nodes := make(map[Vec]rune)
	nodePos := make(map[rune]Vec)
	cost := make(map[[2]rune]int)
	p := Puzzle{walls, nodes, nodePos, cost}
	for y, l := range readFile(s) {
		for x, r := range l {
			switch r {
			case '#':
				walls[Vec{x, y}] = true
			case '.':
				// do nothing
			default:
				v := Vec{x, y}
				nodes[v] = r
				nodePos[r] = v
			}
		}
	}
	for c := range nodePos {
		p.Verts(c)
	}
	p.Home()
	return p
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	p.Solve(false)
	p.Solve(true)
}
