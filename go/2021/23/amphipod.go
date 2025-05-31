package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"cmp"
	"slices"
)

const ROOMS = 16

var (
	stepcost = map[byte]int{'A': 1, 'B': 10, 'C': 100, 'D': 1000}
	depth = 3
)

type Vec [2]int

type PathKey [2]Vec

func (k PathKey) Ordered() PathKey {
	if k[0][1] > k[1][1] || k[0][1] == k[1][1] && k[0][0] > k[1][0] {
		k[0], k[1] = k[1], k[0]
	}
	return k
}

type Node struct {
	sort    byte
	adjs    []Vec
}

type Prawn struct {
	pos     Vec
	goal    byte
}

func (p Prawn) String() string {
	return fmt.Sprintf("%c{%d %d}", p.goal, p.pos[0], p.pos[1])
}

func cmpPrawn(a, b Prawn) int {
	if n := cmp.Compare(a.goal, b.goal); n != 0 {
		return n
	}
	if n := cmp.Compare(a.pos[1], b.pos[1]); n != 0 {
		return n
	}
	return cmp.Compare(a.pos[0], b.pos[0])
}

type State [ROOMS]Prawn

type Path struct {
	steps int
	stops []Vec
}

type Puzzle struct {
	grid        [][]byte
	stops       []Vec
	graph       map[Vec]Node
	start       State
	paths       map[PathKey]Path
	cost        map[State]int
}

type Car struct {
	pos  Vec
	last *Car
}

func (p *Puzzle) GetPath(k PathKey) (path Path) {
	if cache, ok := p.paths[k]; ok {
		return cache
	}
	been := map[Vec]bool{ k[0]: true }
	q := []Car{Car{pos: k[0]}}
	for len(q) > 0 {
		c := q[0]
		q[0].last = nil
		q = q[1:]
		if c.pos == k[1] {
			for c.last != nil {
				if c.pos != k[1] {
					switch p.graph[c.pos].sort {
					case 'H', 'A', 'B', 'C', 'D':
						path.stops = append(path.stops, c.pos)
					}
				}
				c = *c.last
				path.steps++
			}
			break
		}
		for _, v := range p.graph[c.pos].adjs {
			if been[v] {
				continue
			}
			been[v] = true
			q = append(q, Car{v, &c})
		}
	}
	p.paths[k] = path
	return path
}

func (p *Puzzle) GetCost(state State) (cost int) {
	slices.SortFunc(state[:], cmpPrawn)
	if cache, ok := p.cost[state]; ok {
		return cache
	}

	prawns := make(map[Vec]byte)
	for _, o := range state {
		prawns[o.pos] = o.goal
	}

	var fromSort, toSort byte
	var c int
	var hasObs bool
	cost = -1
	for i, o := range state {
		fromSort = p.graph[o.pos].sort
		if o.goal == fromSort {
			skip := true
			for y := o.pos[1] + 1; y <= depth; y++ {
				other, _ := prawns[Vec{o.pos[0], y}]
				if other != o.goal {
					skip = false
					break
				}
			}
			if skip {
				continue
			}
		}
		for _, v := range p.stops {
			if _, ok := prawns[v]; ok {
				continue
			}

			// Check target tile type is allowed (Hallway/Room)
			toSort = p.graph[v].sort
			if fromSort == toSort || toSort != o.goal && toSort != 'H' {
				continue
			}
			if toSort == o.goal {
				skip := false
				for y := v[1] + 1; y <= depth; y++ {
					other, ok := prawns[Vec{v[0], y}]
					if !ok || o.goal != other {
						skip = true
						break
					}
				}
				if skip {
					continue
				}
			}

			// Get path
			path := p.GetPath(PathKey{o.pos, v}.Ordered())

			// Check path for obstructions
			hasObs = false
			for _, stop := range path.stops {
				if _, ok := prawns[stop]; ok {
					hasObs = true
					break
				}
			}
			if hasObs {
				continue
			}

			child := state
			child[i].pos = v
			c = p.GetCost(child)
			if c < 0 {
				continue
			}
			if cost < 0 || cost > c + stepcost[o.goal] * path.steps {
				cost = c + stepcost[o.goal] * path.steps
			}
		}
	}

	//~ fmt.Println("GetCost(): Caching", state, "costs", cost)
	p.cost[state] = cost
	return cost
}

func (p Puzzle) GetAdjs(x, y int) (out []Vec) {
	h := len(p.grid)
	for _, v := range [4]Vec{ Vec{0, -1}, Vec{0, 1}, Vec{-1, 0}, Vec{1, 0} } {
		nx := x + v[0]
		ny := y + v[1]
		if nx < 0 || ny < 0 || ny >= h || nx >= len(p.grid[ny]) {
			continue
		}
		switch p.grid[ny][nx] {
		case '.', 'A', 'B', 'C', 'D':
			out = append(out, Vec{nx, ny})
		}
	}
	return out
}

func (p *Puzzle) Parse() {
	p.graph = make(map[Vec]Node)
	var i int
	for y, r := range p.grid {
		for x, t := range r {
			v := Vec{x, y}
			n := Node{adjs: p.GetAdjs(x, y) }
			switch t {
			case '.':
				n.sort = 'H'
				if len(n.adjs) > 2 {
					n.sort = 'J'
				}
			case 'A', 'B', 'C', 'D':
				switch x {
				case 3:
					n.sort = 'A'
				case 5:
					n.sort = 'B'
				case 7:
					n.sort = 'C'
				case 9:
					n.sort = 'D'
				default:
					fmt.Println("Ignoring...")
					continue
				}
				p.start[i] = Prawn{v, t}
				i++
			default:
				continue
			}
			p.graph[v] = n
			if n.sort != 'J' {
				p.stops = append(p.stops, v)
			}
		}
	}
	slices.SortFunc(p.start[:], cmpPrawn)
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
	for s.Scan() {
		p.grid = append(p.grid, []byte(s.Text()))
	}
	check(s.Err())
	p.Parse()
	p.paths = make(map[PathKey]Path)
	depth = len(p.grid) - 2

	var k State
	var r int
	var b byte
	for x := 3; x < 10; x += 2 {
		switch x {
		case 3:
			b = 'A'
		case 5:
			b = 'B'
		case 7:
			b = 'C'
		case 9:
			b = 'D'
		}
		for y := 2; y <= depth; y++ {
			k[r] = Prawn{Vec{x, y}, b}
			r++
		}
	}

	p.cost = map[State]int{ k: 0 }
	return p
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	cost := p.GetCost(p.start)
	fmt.Println("Least energy cost:", cost)
	fmt.Println("Cost cache size:", len(p.cost))
	fmt.Println("Path cache size:", len(p.paths))
}
