package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"slices"
	"strings"
	"strconv"
	"runtime"
)

type Stack[T any] []T
func (s *Stack[T]) Push(v T) { *s = append(*s, v) }
func (s *Stack[T]) Pop() (v T) {
	old := *s
	n := len(old)
	v = old[n-1]
	*s = old[:n-1]
	return v
}

type Vec struct {
	x, y    int
}

type Move struct {
	pos     Vec
	steps   int
	last    *Move
	puzzle  *Puzzle
}

func (m Move) Next(directed bool) (next Move) {
	var adj []Vec
	for _, v := range m.puzzle.GetAdj(m.pos, directed) {
		if v != m.last.pos {
			adj = append(adj, v)
		}
	}
	n := len(adj)
	if n == 0 {
		return next
	}
	return Move{
		adj[0],
		m.steps + 1,
		&m,
		m.puzzle,
	}
}

func (m Move) IsReversible(directed bool) bool {
	var curr, next Move
	if !directed {
		return true
	}
	p := *m.puzzle
	w := p.width
	slopes := []rune{'<', '^', '>', 'v'}
	curr = m
	for curr.steps > 0 {
		next = *curr.last
		tile := p.grid[next.pos.y * w + next.pos.x]
		if slices.Contains(slopes, tile) {
			return false
		}
		curr = next
	}
	return true
}

func (m Move) FormLink(directed bool) (l Link) {
	var (
		u, v int
		curr, next Move
		f, b bool
		ok bool
	)
	p := m.puzzle
	u, ok = p.nodeIdx[m.last.pos]
	next = m
	ok = false
	for !ok {
		curr = next
		next = curr.Next(directed)
		if next.steps == 0 {
			return l
		}
		v, ok = p.nodeIdx[next.pos]
	}
	if u == v {
		return l
	}
	f = true
	b = next.IsReversible(directed)
	if u > v {
		u, v, f, b = v, u, b, f
	}
	return Link{ u, v, next, f, b }
}

type Node struct {
	pos     Vec
}

type LinkKey struct {
	u, v      int
}

type Link struct {
	u, v      int
	path      Move
	forwards  bool
	backwards bool
}

type Step struct {
	node int
	cost int
	last *Step
}

func (s *Step) Contains(v int) bool {
	curr := &Step{ last: s }
	for curr.last != nil {
		curr = curr.last
		if curr.node == v {
			return true
		}
	}
	return false
}

type Puzzle struct {
	grid            []rune
	start, exit     Vec
	width, height   int
	nodes           []Node
	links           []Link
	nodeIdx         map[Vec]int
	linkCost        map[LinkKey]int
	adjs            map[int][]int
}

func (p Puzzle) Walker(i int, queue chan []*Step, results chan []*Step, done chan Step) {
	var (
		longest Step
		u, v, c int
	)

	for batch := range queue {
		var new []*Step
		for _, s := range batch {
			u = s.node
			if u == 1 {
				if s.cost > longest.cost {
					longest = *s
				}
				continue
			}
			for _, v = range p.adjs[u] {
				if s.Contains(v) {
					continue
				}
				c, _ = p.linkCost[LinkKey{u, v}]
				step := Step{ v, s.cost + c, s }
				new = append(new, &step)
			}
		}
		results <- new
	}

	done <- longest
}

func (p Puzzle) ALongWalk() (longest Step) {
	ncpus   := runtime.NumCPU()
	jobs    := make(chan []*Step, ncpus)
	results := make(chan []*Step, ncpus)
	done    := make(chan Step,    ncpus)
	defer close(results)
	defer close(done)

	// create worker pool
	for i := 0; i < ncpus; i++ {
		go p.Walker(i, jobs, results, done)
	}

	queue := []*Step{ &Step{} }

	var numtasks int
	var pops int
	for len(queue) > 0 {

		qlen := len(queue)
		if qlen < ncpus {
			// queue is too small send everything in one batch
			jobs <- queue
			numtasks++
		} else {
			// split queue into batches and send them to the workers
			batchsize := qlen / ncpus
			for i := 0; i < ncpus; i++ {
				end := (i+1)*batchsize
				if i == ncpus - 1 {
					end = qlen
				}
				jobs <- queue[i*batchsize:end]
				numtasks++
			}
		}
		pops += qlen

		// collate results and requeue
		newqueue := []*Step{}
		for i := 0; i < numtasks; i++ {
			new := <- results
			newqueue = append(newqueue, new...)
		}
		queue = newqueue
		numtasks = 0
	}

	// Tell workers to stop
	close(jobs)

	for i := 0; i < ncpus; i++ {
		s := <- done
		if s.cost > longest.cost {
			longest = s
		}
	}

	fmt.Println("Pops:", pops)

	return longest
}

func (p *Puzzle) UpdateAdjs() {
	p.adjs = make(map[int][]int)
	for _, l := range p.links {
		if l.forwards {
			links, ok := p.adjs[l.u]
			if !ok || !slices.Contains(links, l.v) {
				p.adjs[l.u] = append(p.adjs[l.u], l.v)
			}
		}

		if l.backwards {
			links, ok := p.adjs[l.v]
			if !ok || !slices.Contains(links, l.u) {
				p.adjs[l.v] = append(p.adjs[l.v], l.u)
			}
		}
	}
}

func (p *Puzzle) FindLinks(directed bool) {
	p.links = p.links[:0]
	p.linkCost = make(map[LinkKey]int)

	start := Move{ pos: p.start, puzzle: p }

	var open Stack[Move]
	open.Push(Move{ Vec{ p.start.x, p.start.y + 1 }, 1, &start, p})

	var seen []int
	seen = append(seen, 0, 1)

	for len(open) > 0 {
		s := open.Pop()
		l := s.FormLink(directed)
		if l.u == l.v {
			continue
		}
		if l.forwards {
			p.linkCost[LinkKey{l.u, l.v}] = l.path.steps
		}
		if l.backwards {
			p.linkCost[LinkKey{l.v, l.u}] = l.path.steps
		}
		p.links = append(p.links, l)

		pos := p.nodes[l.u].pos
		if !slices.Contains(seen, l.u) {
			node := Move{ pos: pos, puzzle: p }
			for _, v := range p.GetAdj(pos, directed) {
				open = append(open, Move{ v, 1, &node, p })
			}
			seen = append(seen, l.u)
		}

		pos = p.nodes[l.v].pos
		if !slices.Contains(seen, l.v) {
			node := Move{ pos: pos, puzzle: p }
			for _, v := range p.GetAdj(pos, directed) {
				open = append(open, Move{ v, 1, &node, p })
			}
			seen = append(seen, l.v)
		}
	}
}

func (p *Puzzle) FindNodes() {
	p.nodes = p.nodes[:0]
	p.nodeIdx = make(map[Vec]int)

	p.nodeIdx[p.start] = len(p.nodes)
	p.nodes = append(p.nodes, Node{ p.start })

	p.nodeIdx[p.exit] = len(p.nodes)
	p.nodes = append(p.nodes, Node{ p.exit })

	w := p.width
	for i, t := range p.grid {
		if t != '.' {
			continue
		}

		v := Vec{i % w, i / w}
		if len(p.GetAdj(v, false)) < 3 {
			continue
		}

		p.nodeIdx[v] = len(p.nodes)
		p.nodes = append(p.nodes, Node{ v })
	}
}

func (p *Puzzle) BuildGraph(directed bool) {
	p.FindNodes()
	p.FindLinks(directed)
	p.UpdateAdjs()
}

func (p Puzzle) GetAdj(pos Vec, directed bool) (adj []Vec) {
	w, h := p.width, p.height
	x, y := pos.x, pos.y

	for _, v := range []Vec{
		Vec{x, y-1},
		Vec{x+1, y},
		Vec{x, y+1},
		Vec{x-1, y},
	} {
		if v.x < 0 || v.y < 0 || v.x >= w || v.y >= h {
			continue
		}

		tile := p.grid[v.y * w + v.x]
		if tile == '#' {
			continue
		}

		if directed {
			if  tile == '^' && v.y >= y ||
				tile == '>' && v.x <= x ||
				tile == 'v' && v.y <= y ||
				tile == '<' && v.x >= x {
				continue
			}
		}
		adj = append(adj, v)
	}
	return adj
}

func (p Puzzle) DrawGrid() {
	var buf strings.Builder
	w := p.width
	for i, r := range p.grid {
		if i % w == 0 {
			buf.WriteRune('\n')
		}
		x := i % w
		y := i / w
		n, ok := p.nodeIdx[Vec{x, y}]
		if ok {
			s := strconv.Itoa(n)
			rs := []rune(s)
			nrs := len(rs)
			r = rs[nrs-1]
		}
		buf.WriteRune(r)
	}
	fmt.Println(buf.String())
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		row := []rune(s.Text())
		if p.width == 0 {
			p.width = len(row)
		}
		p.grid = append(p.grid, row...)
	}
	check(s.Err())
	p.height    = len(p.grid) / p.width
	p.start     = Vec{ 1, 0 }
	p.exit      = Vec{ p.width-2, p.height-1 }
	return p
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var path Step
	flag.Parse()
	p := load(flag.Arg(0))

	p.BuildGraph(true)
	p.DrawGrid()
	path = p.ALongWalk()
	fmt.Println("Part 1:", path)

	p.BuildGraph(false)
	//~ p.DrawGrid()
	path = p.ALongWalk()
	fmt.Println("Part 2:", path)
}
