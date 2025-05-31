package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"container/heap"
)

var adjacent = []Vertex{
	Vertex{0, -1}, Vertex{0, 1}, Vertex{-1, 0}, Vertex{1, 0},
}

type Vertex struct {x, y int}

type Node struct {
	pos      Vertex
	cost     int
	estimate int
	last     *Node
}

type Puzzle struct {
	width  int
	height int
	risk   []int
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func abs(n int) int {
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
		line := s.Text()
		if p.width == 0 {
			p.width = int(len(line))
		}
		p.height++
		for _, r := range line {
			p.risk = append(p.risk, int(r - '0'))
		}
	}
	check(s.Err())
	return p
}

func (p Puzzle) FindPath() (zero Node) {
	best := map[Vertex]int{}
	exit := Vertex{p.width - 1, p.height - 1}
	pq := &PriorityQueue{Node{pos: Vertex{0, 0}}}
	heap.Init(pq)
	for len(*pq) > 0 {
		s := heap.Pop(pq).(Node)
		if s.pos == exit {
			return s
		}
		adjs := p.GetAdjs(s)
		for _, adj := range adjs {
			b, ok := best[adj.pos]
			if ok && b <= adj.cost {
				continue
			}
			best[adj.pos] = adj.cost
			adj.estimate = adj.cost + exit.Dist(adj.pos)
			heap.Push(pq, adj)
		}
	}
	return zero
}

func (p Puzzle) GetAdjs(s Node) (out []Node) {
	w, h := p.width, p.height
	for _, v := range adjacent {
		x := s.pos.x + v.x
		y := s.pos.y + v.y
		if x < 0 || y < 0 || x >= w || y >= h {
			continue
		}
		cost := s.cost + p.risk[y * w + x]
		out = append(out, Node{Vertex{x, y}, cost, 0, &s})
	}
	return out
}

func (v1 Vertex) Dist(v2 Vertex) int {
	return abs(v1.x - v2.x) + abs(v1.y - v2.y)
}

func (p Puzzle) Enlarge() Puzzle {
	var x, y int
	w, h := p.width, p.height
	risk := p.risk
	ow, oh := w * 5, h * 5
	orisk := make([]int, ow * oh)
	for y = 0; y < oh; y++ {
		for x = 0; x < ow; x++ {
			u := x % w
			v := y % h
			d := x / w + y / h
			r := risk[v * w + u]
			orisk[y * ow + x] = (r + d - 1) % 9 + 1
		}
	}
	return Puzzle{ow, oh, orisk}
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	best := p.FindPath()
	fmt.Println(best)
	fmt.Println("Part 1:", best.cost)

	q := p.Enlarge()
	p2 := q.FindPath()
	fmt.Println(p2)
	fmt.Println("Part 2:", p2.cost)
}

type PriorityQueue []Node
func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
func (pq *PriorityQueue) Push(x any) { *pq = append(*pq, x.(Node)) }
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}
