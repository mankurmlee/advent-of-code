package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"container/heap"
)

type Vec struct {
	x, y int
}

type Finder struct {
	pos    Vec
	steps  int
	last   *Finder
	puzzle *Puzzle
}

type Puzzle struct {
	grid   []int
	width  int
	height int
	start  Vec
	exit   Vec
}

type PriorityQueue []*Finder
func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].RoughCost() < pq[j].RoughCost()
}
func (pq *PriorityQueue) Push(x any) { *pq = append(*pq, x.(*Finder)) }
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return item
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)

	var w, h, z int
	for s.Scan() {
		row := s.Text()
		if w == 0 {
			w = len(row)
		}
		for x, r := range row {
			switch r {
			case 'S':
				p.start = Vec{x, h}
				z = 1
			case 'E':
				p.exit = Vec{x, h}
				z = 26
			default:
				z = int(r - 'a') + 1
			}
			p.grid = append(p.grid, z)
		}
		h++
	}
	check(s.Err())

	p.width, p.height = w, h
	return p
}

func Abs(n int) int {
	if n >= 0 {
		return n
	}
	return -n
}

func (f Finder) RoughCost() int {
	return Abs(f.puzzle.exit.x - f.pos.x) + Abs(f.puzzle.exit.y - f.pos.y) + f.steps
}

func (p Puzzle) AStar() Finder {
	seen := make(map[Vec]int)
	f := &Finder{ pos : p.start, puzzle : &p }
	pq := &PriorityQueue{ f }
	heap.Init(pq)

	for len(*pq) > 0 {
		f = heap.Pop(pq).(*Finder)
		if f.pos ==  p.exit {
			return *f
		}
		for _, g := range f.AdjacentUp() {
			if s, ok := seen[g.pos]; ok && g.steps >= s {
				continue
			}
			seen[g.pos] = g.steps
			heap.Push(pq, g)
		}
	}

	return *f
}

func (f Finder) AdjacentUp() (adjs []*Finder) {
	w, h := f.puzzle.width, f.puzzle.height
	x, y := f.pos.x, f.pos.y
	z := f.puzzle.grid[y * w + x]
	for _, v := range []Vec{
		Vec{x, y - 1},
		Vec{x + 1, y},
		Vec{x, y + 1},
		Vec{x - 1, y},
	} {
		if v.x < 0 || v.x >= w || v.y < 0 || v.y >= h {
			continue
		}
		if f.puzzle.grid[v.y * w + v.x] > z + 1 {
			continue
		}
		if f.last != nil && f.last.pos == v {
			continue
		}
		adjs = append(adjs, &Finder{ v, f.steps + 1, &f, f.puzzle })
	}
	return adjs
}

func (p Puzzle) BFS() Finder {
	w := p.width
	seen := make(map[Vec]int)
	f := &Finder{ pos : p.exit, puzzle : &p }
	q := []*Finder{ f }

	for len(q) > 0 {
		f = q[0]
		q = q[1:]
		if p.grid[f.pos.y * w + f.pos.x] == 1 {
			return *f
		}
		for _, g := range f.AdjacentDown() {
			if s, ok := seen[g.pos]; ok && g.steps >= s {
				continue
			}
			seen[g.pos] = g.steps
			q = append(q, g)
		}
	}

	return *f
}

func (f Finder) AdjacentDown() (adjs []*Finder) {
	w, h := f.puzzle.width, f.puzzle.height
	x, y := f.pos.x, f.pos.y
	z := f.puzzle.grid[y * w + x]
	for _, v := range []Vec{
		Vec{x, y - 1},
		Vec{x + 1, y},
		Vec{x, y + 1},
		Vec{x - 1, y},
	} {
		if v.x < 0 || v.x >= w || v.y < 0 || v.y >= h {
			continue
		}
		if f.puzzle.grid[v.y * w + v.x] < z - 1 {
			continue
		}
		if f.last != nil && f.last.pos == v {
			continue
		}
		adjs = append(adjs, &Finder{ v, f.steps + 1, &f, f.puzzle })
	}
	return adjs
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	up := p.AStar()
	fmt.Println("Part 1:", up.steps, "steps")

	down := p.BFS()
	fmt.Println("Part 2:", down.steps, "steps")
}
