package main

import (
	"flag"
	"os"
	"bufio"
	"strconv"
	"fmt"
	"container/heap"
	"slices"
	"runtime"
)

type Direction int

const (
	None  Direction = iota
	North
	East
	South
	West
)

type Context struct {
	grid           [][]int
	width, height  int
	start, exit    Vec
	minrun, maxrun int
	maxthreads     int
}

type Vec struct {
	x, y int
}

type Move struct {
	pos  Vec
	cost int
	last Trail
}

type Trail *Move

type Key struct {
	pos      Vec
	lastdir  Direction
	runlen   int
}

type Searcher struct {
	ctx      *Context
	move     Move
	lastdir  Direction
	runlen   int
	trail    Trail
}

func (s Searcher) GetNeighbours(ch chan []*Searcher) {
	var searchers []*Searcher
	x, y := s.move.pos.x, s.move.pos.y
	w, h := s.ctx.width, s.ctx.height

	for _, v := range [4]Vec{
		Vec{x - 1, y},
		Vec{x + 1, y},
		Vec{x, y - 1},
		Vec{x, y + 1},
	} {
		// Check we're still on the grid
		if  v.x < 0 || v.x >= w ||
			v.y < 0 || v.y >= h {
			continue
		}

		// Check we're not going back to the last position
		if s.trail.last != nil && s.trail.last.pos == v {
			continue
		}

		d := None
		switch {
		case v.y < y:
			d = North
		case v.x > x:
			d = East
		case v.y > y:
			d = South
		case v.x < x:
			d = West
		}

		if s.runlen < s.ctx.minrun && s.lastdir != None && d != s.lastdir {
			continue
		}

		runlen := 1
		if d == s.lastdir {
			runlen = s.runlen + 1
		}

		if runlen > s.ctx.maxrun {
			continue
		}

		// Add to output
		m := Move{
			v,
			s.move.cost + s.ctx.grid[v.y][v.x],
			s.trail,
		}

		searchers = append(searchers, &Searcher{
			s.ctx,
			m,
			d,
			runlen,
			&m,
		})
	}
	ch <- searchers
	return
}

func (ctx Context) DoSearch() (path Trail) {
	// Keep a table of best scores
	best := make(map[Key]int)

	// Create the starting move
	move := Move{ pos: ctx.start }

	// Add starting move to priority queue
	pq := PriorityQueue{&Searcher{
		ctx  : &ctx,
		move : move,
		trail: &move,
	}}
	heap.Init(&pq)

	ch := make(chan []*Searcher, ctx.maxthreads)
	defer close(ch)

	for pq.Len() > 0 {
		var tasks []*Searcher
		for i := 0; i < ctx.maxthreads && pq.Len() > 0; i++ {
			current := heap.Pop(&pq).(*Searcher)

			// Have we reached the exit?
			if  current.move.pos == ctx.exit &&
				current.runlen >= ctx.minrun {
				return current.trail
			}

			tasks = append(tasks, current)
		}

		// Get the neighbours
		for _, t := range tasks {
			go t.GetNeighbours(ch)
		}

		var neighbours []*Searcher
		for range tasks {
			neighbours = append(neighbours, <-ch...)
		}

		for _, s := range neighbours {
			// Check if better than current best
			key := Key{
				s.move.pos,
				s.lastdir,
				s.runlen,
			}
			cost, ok := best[key]
			if ok && s.move.cost >= cost {
				continue
			}
			best[key] = s.move.cost

			// Add to pq
			heap.Push(&pq, s)
		}
	}

	return path
}

func (ctx Context) DrawTrail(t Trail) {
	var path []Vec
	for {
		path = append(path, t.pos)
		if t.last == nil {
			break
		}
		t = t.last
	}
	for y, row := range ctx.grid {
		for x, i := range row {
			s := strconv.Itoa(i)
			if slices.Contains(path, Vec{x, y}) {
				s = "*"
			}
			fmt.Print(s)
		}
		fmt.Println()
	}
}

func load(filename string) (g [][]int) {
	f, err := os.Open(filename)
	check(err)
	s := bufio.NewScanner(f)
	for s.Scan() {
		var r []int
		for _, ch := range s.Text() {
			c, err := strconv.Atoi(string(ch))
			check(err)
			r = append(r, c)
		}
		g = append(g, r)
	}
	check(s.Err())
	return g
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	g := load(flag.Arg(0))
	w, h := len(g[0]), len(g)

	var ctx Context
	ctx.grid     = g
	ctx.width    = w
	ctx.height   = h
	ctx.start    = Vec{0, 0}
	ctx.exit     = Vec{w-1, h-1}
	ctx.maxrun   = 3
	ctx.maxthreads = runtime.NumCPU() * 30

	path1 := ctx.DoSearch()
	fmt.Println("Part 1")
	ctx.DrawTrail(path1)
	answer1 := path1.cost

	ctx.minrun = 4
	ctx.maxrun = 10

	path2 := ctx.DoSearch()
	fmt.Println("Part 2")
	ctx.DrawTrail(path2)
	answer2 := path2.cost

	fmt.Println("Answer to part 1 is:", answer1)
	fmt.Println("Answer to part 2 is:", answer2)
}

type PriorityQueue []*Searcher

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].move.cost < pq[j].move.cost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(*Searcher))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}
