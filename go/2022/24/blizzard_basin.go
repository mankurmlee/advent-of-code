package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"slices"
	"container/heap"
)

var dirs = []Vertex{
	Vertex{ 0, -1},
	Vertex{ 1,  0},
	Vertex{ 0,  1},
	Vertex{-1,  0},
	Vertex{ 0,  0},
}

type PriorityQueue []*State
func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].estimate < pq[j].estimate }
func (pq *PriorityQueue) Push(x any) { *pq = append(*pq, x.(*State)) }
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return item
}

type Vertex struct {
	x, y int
}

type Blizzard struct {
	pos, dir Vertex
}

type Blizzards struct {
	bslice []Blizzard
	bmap   map[Vertex]bool
}

type State struct {
	pos      Vertex
	cost     int
	estimate int
}

type SeenKey struct {
	pos      Vertex
	cost     int
}

type Puzzle struct {
	grid    [][]byte
	width   int
	height  int
	start   Vertex
	exit    Vertex
	blizs   Blizzards
}

func check(err error) {
	if err != nil {
		panic(err)
	}
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
	for s.Scan() {
		p.grid = append(p.grid, []byte(s.Text()))
	}
	check(s.Err())
	p.width  = len(p.grid[0])
	p.height = len(p.grid)
	p.start  = Vertex{ slices.Index(p.grid[0], '.'), 0 }
	p.exit   = Vertex{ slices.Index(p.grid[p.height-1], '.'), p.height-1 }
	for y := 0; y < p.height; y++ {
		for x := 0; x < p.width; x++ {
			switch p.grid[y][x] {
			case '^':
				p.blizs.bslice = append(p.blizs.bslice, Blizzard{Vertex{x, y}, dirs[0]})
			case '>':
				p.blizs.bslice = append(p.blizs.bslice, Blizzard{Vertex{x, y}, dirs[1]})
			case 'v':
				p.blizs.bslice = append(p.blizs.bslice, Blizzard{Vertex{x, y}, dirs[2]})
			case '<':
				p.blizs.bslice = append(p.blizs.bslice, Blizzard{Vertex{x, y}, dirs[3]})
			}
		}
	}
	p.blizs.CreateMap()
	return p
}

func (blizs *Blizzards) CreateMap() {
	bmap := make(map[Vertex]bool, len(blizs.bslice))
	for _, b := range blizs.bslice {
		bmap[b.pos] = true
	}
	blizs.bmap = bmap
}

func (v1 Vertex) Add (v2 Vertex) Vertex {
	return Vertex{v1.x + v2.x, v1.y + v2.y}
}

func (v1 Vertex) Dist (v2 Vertex) int {
	return Abs(v2.x - v1.x) + Abs(v2.y - v1.y)
}

func (p Puzzle) AStar() (State, Blizzards) {
	seen := make(map[SeenKey]bool)
	s := &State{
		pos     : p.start,
		estimate: p.start.Dist(p.exit),
	}
	var pq PriorityQueue
	heap.Init(&pq)
	heap.Push(&pq, s)

	blizstates := make(map[int]Blizzards)
	blizstates[1] = p.Update(p.blizs)

	for len(pq) > 0 {
		s = heap.Pop(&pq).(*State)
		if s.pos == p.exit {
			break
		}
		blizs, ok := blizstates[s.cost + 1]
		if !ok {
			blizs = p.Update(blizstates[s.cost])
			blizstates[s.cost + 1] = blizs
		}
		for _, adj := range p.NextMoves(s, blizs) {
			if seen[SeenKey{adj.pos, adj.cost}] {
				continue
			}
			seen[SeenKey{adj.pos, adj.cost}] = true
			heap.Push(&pq, adj)
		}
	}

	return *s, blizstates[s.cost]
}

func (p Puzzle) NextMoves(s *State, blizs Blizzards) (moves []*State) {
	w, h := p.width, p.height
	for _, dir := range dirs {
		v := s.pos.Add(dir)
		if v.x < 0 || v.y < 0 || v.x >= w || v.y >= h {
			continue
		}
		if p.grid[v.y][v.x] == '#' {
			continue
		}
		if blizs.bmap[v] {
			continue
		}
		est := s.cost + 1 + v.Dist(p.exit)
		move := &State{ v, s.cost + 1, est }
		moves = append(moves, move)
	}
	return moves
}

func (p Puzzle) Update(now Blizzards) (next Blizzards) {
	for _, b := range now.bslice {
		v := b.pos.Add(b.dir)
		if p.grid[v.y][v.x] == '#' {
			switch b.dir {
			case dirs[0]:
				v.y = p.height - 2
			case dirs[1]:
				v.x = 1
			case dirs[2]:
				v.y = 1
			case dirs[3]:
				v.x = p.width - 2
			}
		}
		next.bslice = append(next.bslice, Blizzard{v, b.dir})
	}
	next.CreateMap()
	return next
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	s, blizs := p.AStar()
	cost := s.cost
	fmt.Println("Part 1:", cost)

	p.start, p.exit = p.exit, p.start
	p.blizs = blizs
	s, blizs = p.AStar()
	cost += s.cost
	fmt.Println("Back:", s.cost)

	p.start, p.exit = p.exit, p.start
	p.blizs = blizs
	s, _ = p.AStar()
	cost += s.cost
	fmt.Println("And there again:", s.cost)
	fmt.Println("Part 2:", cost)
}
