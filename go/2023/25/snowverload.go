package main

import (
	"flag"
	"os"
	"bufio"
	"strings"
	"slices"
	"fmt"
	"container/heap"
)

type Node string

type Link struct {
	a, b Node
}

func (l Link) In(slice []Link) bool {
	return slices.Contains(slice, l.Order())
}

func (l Link) Order() Link {
	if l.a > l.b {
		return Link{l.b, l.a}
	}
	return l
}

type Move struct {
	pos   Node
	steps int
	last  *Move
}

func (m Move) GetLinks() (links []Link) {
	curr := m
	var n Move
	var b Node
	a := m.pos
	for i := 0; i < m.steps; i++ {
		n = *curr.last
		b = n.pos
		links = append(links, Link{a, b}.Order())
		curr = n
		a = b
	}
	return links
}

type Puzzle struct {
	nodes []Node
	adjs  map[Node][]Node
}

func (p Puzzle) CountNodes(start Node, ignore []Link) int {
	open := []Node{start}
	seen := make(map[Node]bool)
	for len(open) > 0 {
		n := len(open)
		a := open[n-1]
		open = open[:n-1]
		for _, b := range p.adjs[a] {
			l := Link{a, b}
			if l.In(ignore) {
				continue
			}
			_, ok := seen[b]
			if ok {
				continue
			}
			seen[b] = true
			open = append(open, b)
		}
	}

	return len(seen)
}

func (p Puzzle) FindCriticalLinks(num int, check []Link, cut []Link) []Link {
	var m Move
	var moreCut []Link
	for _, l := range check {
		//~ fmt.Println("Num:", num, "Link:", l, "Cut:", cut)

		// Add another incision
		moreCut = append(cut, l)

		// Try and find an alternate path between the two nodes
		m = p.FindPath(l, moreCut)

		// If a path doesn't exist then it means that severing this link
		// will cause a split
		if m.steps == 0 {
			return []Link{l}
		}

		if num > 1 {
			res := p.FindCriticalLinks(num - 1, m.GetLinks(), moreCut)
			if len(res) == num - 1 {
				return append(res, l)
			}
		}
	}
	return []Link{}
}

func (p Puzzle) GetLinks() (links []Link) {
	for a, nodes := range p.adjs {
		for _, b := range nodes {
			if a < b {
				links = append(links, Link{a, b})
			}
		}
	}
	return links
}

func (p Puzzle) GetAdjacent(m Move, ignore []Link) (moves []Move) {
	for _, n := range p.adjs[m.pos] {
		if m.steps > 0 && m.last.pos == n {
			continue
		}
		l := Link{m.pos, n}
		if l.In(ignore) {
			continue
		}
		moves = append(moves, Move{ n, m.steps + 1, &m })
	}
	return moves
}

func (p Puzzle) FindPath(l Link, ignore []Link) (trail Move) {
	best := make(map[Node]int)

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, Move{ pos: l.a })
	for len(*pq) > 0 {
		move := heap.Pop(pq).(Move)
		if move.pos == l.b {
			return move
		}
		for _, m := range p.GetAdjacent(move, ignore) {
			b, ok := best[m.pos]
			if ok && m.steps > b {
				continue
			}
			best[m.pos] = m.steps
			heap.Push(pq, m)
		}
	}

	return trail
}

type PriorityQueue []Move
func (h  PriorityQueue) Len() int           { return len(h) }
func (h  PriorityQueue) Less(i, j int) bool { return h[i].steps < h[j].steps }
func (h  PriorityQueue) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *PriorityQueue) Push(x any)         { *h = append(*h, x.(Move)) }
func (h *PriorityQueue) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	s := bufio.NewScanner(f)
	r := strings.NewReplacer(":", "")
	p.adjs = make(map[Node][]Node)
	for s.Scan() {
		wiringData := strings.Fields(r.Replace(s.Text()))
		for _, nData := range wiringData {
			n := Node(nData)
			if !slices.Contains(p.nodes, n) {
				p.nodes = append(p.nodes, n)
			}
		}
		for _, bData := range wiringData[1:] {
			a := Node(wiringData[0])
			b := Node(bData)
			p.adjs[a] = append(p.adjs[a], b)
			p.adjs[b] = append(p.adjs[b], a)
		}
	}
	check(s.Err())
	return p
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	fmt.Println("Number of nodes:", len(p.nodes))

	result := p.FindCriticalLinks(3, p.GetLinks(), []Link{})
	fmt.Println("Result:", result)

	// [{ltn trh} {fdb psj} {nqh rmt}]
	i := p.CountNodes(result[0].a, result)
	j := len(p.nodes) - i
	res := i * j
	fmt.Println("Split between", i, "and", j, "nodes:", res)
}
