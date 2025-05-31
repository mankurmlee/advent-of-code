// Probably makes the following assumptions
// Input is:
// * a square map
// * has odd number of cols/rows
// * start position is dead centre
// * first and last rows and columns are clear
//
package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"slices"
	"container/heap"
	"strconv"
)

const (
	DIR_N  = 0
	DIR_E  = 1
	DIR_S  = 2
	DIR_W  = 3
	DIR_NE = 4
	DIR_SE = 5
	DIR_SW = 6
	DIR_NW = 7
)

type PriorityQueue []*State
func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].guess < pq[j].guess }
func (pq *PriorityQueue) Push(x any) { *pq = append(*pq, x.(*State)) }
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return item
}

type Garden struct {
	width, height int
	start         Vec
	tiles         []rune
	dir           []Direction
	samePlots     int
	diffPlots     int
}

type Vec struct {
	x, y int
}

type State struct {
	pos         Vec
	steps       int
	guess       int
	last        *State
}

type Direction struct {
	fill      int
	offset    int
	searchers []State
	offset0   int
	search0   []State
}

func (g Garden) GetDiagCount(dir Direction, steps int) (count int) {
	w, h := g.width, g.height
	if steps <= dir.offset {
		return 0
	}
	var fill int
	if steps >= dir.offset + dir.fill {
		fill = 1 + (steps - dir.offset - dir.fill) / h
		evens := (fill + 1) / 2
		evens *= evens
		odds  := fill / 2
		odds  += odds * odds
		count += g.samePlots * evens + g.diffPlots * odds
	}
	reach := 1 + (steps - dir.offset - 1) / h

	for i := fill; i < reach; i++ {
		count += g.CountPlots(dir.searchers, steps - dir.offset - w * i) * (i + 1)
	}

	fmt.Println("Diag:", count, "Reach", reach, "Fill", fill)
	return count
}

func (g Garden) GetOrthoCount(dir Direction, steps int) (count int) {
	h := g.height
	var reach, fill int

	if steps >= dir.offset + 1 {
		reach = 1 + (steps - dir.offset - 1) / h
	} else if steps > dir.offset0 {
		reach = 1
	}

	if steps >= dir.offset + dir.fill {
		fill = 1 + (steps - dir.offset - dir.fill) / h
	}

	for i := fill; i < reach; i++ {
		if i == 0 {
			count += g.CountPlots(dir.search0, steps)
		} else {
			count += g.CountPlots(dir.searchers, steps - dir.offset - h * i)
		}
	}

	oddRows  := (fill + 1) / 2
	evenRows := fill / 2
	count += evenRows * g.samePlots + oddRows * g.diffPlots

	fmt.Println("Ortho:", count, "Reach:", reach, "Fill:", fill)
	return count
}


func (g Garden) GetCount(steps int) (count int) {

	count += g.GetOrthoCount(g.dir[DIR_N], steps)
	count += g.GetOrthoCount(g.dir[DIR_E], steps)
	count += g.GetOrthoCount(g.dir[DIR_S], steps)
	count += g.GetOrthoCount(g.dir[DIR_W], steps)
	count += g.GetDiagCount(g.dir[DIR_NE], steps)
	count += g.GetDiagCount(g.dir[DIR_SE], steps)
	count += g.GetDiagCount(g.dir[DIR_SW], steps)
	count += g.GetDiagCount(g.dir[DIR_NW], steps)

	// Central garden
	c := g.CountPlots([]State{State{ pos: g.start }}, steps)
	fmt.Println("Central:", c)
	count += c

	return count
}

func (g Garden) CountPlots(start []State, steps int) (count int) {
	seen := make(map[Vec]int)

	var queue []State
	for _, s := range start {
		if s.steps > steps {
			continue
		}
		queue = append(queue, s)
		seen[s.pos] = s.steps
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if curr.steps == steps {
			continue
		}
		for _, s := range g.Explore(curr) {
			_, ok := seen[s.pos]
			if ok && s.steps >= seen[s.pos] {
				continue
			}
			seen[s.pos] = s.steps
			queue = append(queue, s)
		}
	}

	evenodd := steps % 2
	for _, v := range seen {
		if evenodd == v % 2 {
			count++
		}
	}

	return count
}

func (g Garden) Clip(v Vec) Vec {
	x := v.x % g.width
	y := v.y % g.height
	if x < 0 {
		x += g.width
	}
	if y < 0 {
		y += g.height
	}
	return Vec{x, y}
}

// Calculate offsets required to do infinite calculations
func (g *Garden) GetStats(steps int) {
	w := g.width
	h := g.height

	// North
	g.dir = append(g.dir, g.GetOrthoStats(
		Vec{0, -1},
		Vec{0, -h - 1},
		Vec{1, 0}),
	)

	// East
	g.dir = append(g.dir, g.GetOrthoStats(
		Vec{w, 0},
		Vec{w * 2, 0},
		Vec{0, 1}),
	)

	// South
	g.dir = append(g.dir, g.GetOrthoStats(
		Vec{0, h},
		Vec{0, h * 2},
		Vec{1, 0}),
	)

	// West
	g.dir = append(g.dir, g.GetOrthoStats(
		Vec{-1, 0},
		Vec{-w - 1, 0},
		Vec{0, 1}),
	)

	// North East
	g.dir = append(g.dir, g.GetDiagStats(Vec{w, -1}))

	// South East
	g.dir = append(g.dir, g.GetDiagStats(Vec{w, h}))

	// South West
	g.dir = append(g.dir, g.GetDiagStats(Vec{-1, h}))

	// North West
	g.dir = append(g.dir, g.GetDiagStats(Vec{-1, -1}))

	// Count odd & even plots
	mid := []State{State{ pos: g.start }}
	g.samePlots = g.CountPlots(mid, steps)
	g.diffPlots = g.CountPlots(mid, steps+1)

	for _, d := range g.dir {
		fmt.Print(d.offset, d.fill, "; ")
	}
	fmt.Println()
	fmt.Println("Same Plots:", g.samePlots, "Diff Plots:", g.diffPlots)
}

// Assumes square map
func (g Garden) GetOrthoStats(base0, base Vec, mul Vec) Direction {
	w, h := g.width, g.height

	var offset0 int
	var state0 []State
	for x := 0; x < w; x++ {
		v := mul.Scalar(x).Add(base0)
		p := g.AStar(v)
		state0 = append(state0, State{pos: g.Clip(p.pos), steps: p.steps})
		if offset0 == 0 || p.steps < offset0 {
			offset0 = p.steps
		}
	}
	offset0--

	var offset int
	var states []State
	for x := 0; x < w; x++ {
		v := mul.Scalar(x).Add(base)
		p := g.AStar(v)
		states = append(states, State{pos: g.Clip(p.pos), steps: p.steps})
		if offset == 0 || p.steps < offset {
			offset = p.steps
		}
	}
	for i := 0; i < len(states); i++ {
		states[i].steps -= offset - 1
	}
	offset -= h + 1
	fill := g.GetFillCost(states)
	return Direction{ fill, offset, states, offset0, state0 }
}

func (g Garden) GetDiagStats(base Vec) Direction {
	p := g.AStar(base)
	states := []State{ State{pos: g.Clip(p.pos), steps: 1} }
	offset := p.steps - 1
	fill := g.GetFillCost(states)
	return Direction{ fill, offset, states, offset, states }
}

func (v1 Vec) Add(v2 Vec) Vec {
	return Vec{v1.x + v2.x, v1.y + v2.y}
}

func (v Vec) Scalar(i int) Vec {
	return Vec{v.x * i, v.y * i}
}

func (g Garden) Explore(curr State) (next []State) {
	x := curr.pos.x
	y := curr.pos.y
	w, h := g.width, g.height
	for _, v := range []Vec{
		Vec{x - 1, y},
		Vec{x + 1, y},
		Vec{x, y - 1},
		Vec{x, y + 1},
	} {
		if curr.last != nil && curr.last.pos == v {
			continue
		}
		if v.x < 0 || v.y < 0 || v.x >= w || v.y >= h {
			continue
		}
		if g.tiles[v.y * w + v.x] == '#' {
			continue
		}
		next = append(next, State{v, curr.steps + 1, 0, &curr})
	}
	return next
}

func (g Garden) GetFillCost(start []State) (max int) {
	seen := make(map[Vec]int)

	queue := make([]State, len(start))
	copy(queue, start)

	for _, s := range queue {
		seen[s.pos] = s.steps
	}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		for _, s := range g.Explore(curr) {
			_, ok := seen[s.pos]
			if ok && s.steps >= seen[s.pos] {
				continue
			}
			seen[s.pos] = s.steps
			queue = append(queue, s)
		}
	}
	for _, v := range seen {
		if max < v {
			max = v
		}
	}
	return max
}

func (g Garden) AStar(target Vec) (path State) {
	seen := make(map[Vec]int)
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &State{
		pos: g.start,
		guess: target.DistTo(g.start),
	})
	for len(*pq) > 0 {
		curr := heap.Pop(pq).(*State)
		seen[curr.pos] = curr.steps
		if curr.pos == target {
			return *curr
		}
		for _, p := range g.AStarExplore(*curr) {
			steps, ok := seen[p.pos]
			if ok && steps <= p.steps {
				continue
			}
			seen[p.pos] = p.steps
			p.guess = target.DistTo(p.pos) + p.steps
			heap.Push(pq, p)
		}
	}
	return path
}

func (g Garden) AStarExplore(curr State) (next []*State) {
	x := curr.pos.x
	y := curr.pos.y
	for _, v := range []Vec{
		Vec{x - 1, y},
		Vec{x + 1, y},
		Vec{x, y - 1},
		Vec{x, y + 1},
	} {
		if curr.last != nil && curr.last.pos == v {
			continue
		}
		if g.Get(v) == '#' {
			continue
		}
		next = append(next, &State{ v, curr.steps + 1, 0, &curr })
	}
	return next
}

func (g Garden) Get(pos Vec) rune {
	x := pos.x % g.width
	y := pos.y % g.height
	if x < 0 {
		x += g.width
	}
	if y < 0 {
		y += g.height
	}
	return g.tiles[y * g.width + x]
}

func (v1 Vec) DistTo(v2 Vec) int {
	return Abs(v2.x - v1.x) + Abs(v2.y - v1.y)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func load(filename string) (g Garden) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)

	for s.Scan() {
		row := []rune(s.Text())
		if g.width == 0 {
			g.width = len(row)
		}
		g.tiles = append(g.tiles, row...)
		g.height++
	}
	check(s.Err())

	i := slices.Index(g.tiles, 'S')
	g.tiles[i] = '.'
	g.start = Vec{ i % g.width, i / g.width }

	return g
}

func Atoi(A string) int {
	i, err := strconv.Atoi(A)
	check(err)
	return i
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var steps int
	flag.Parse()
	filename := flag.Arg(0)
	stepData := flag.Arg(1)
	g := load(filename)
	if stepData == "" {
		steps = 26501365
	}

	g.GetStats(steps)

	count := g.GetCount(steps)
	fmt.Println("Plots for", steps, "steps is:", count)
}
