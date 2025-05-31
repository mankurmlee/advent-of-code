package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"slices"
)

type Vertex struct {x, y int}

type Puzzle struct {
	width  int
	height int
	grid   []int
}

var adjacent = []Vertex{
	Vertex{0, -1}, Vertex{ 1, 0},
	Vertex{0,  1}, Vertex{-1, 0},
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
		line := s.Text()
		p.height++
		if p.width == 0 {
			p.width = len(line)
		}
		for _, r := range line {
			p.grid = append(p.grid, int(r - '0'))
		}
	}
	check(s.Err())
	return p
}

func (p Puzzle) GetLowPoints() (lowpoints []Vertex) {
	w, h := p.width, p.height
	for i, alt := range p.grid {
		isLowPoint := true
		x := i % w
		y := i / w
		for _, v := range adjacent {
			vx, vy := x + v.x, y + v.y
			if vx < 0 || vy < 0 || vx >= w || vy >= h {
				continue
			}
			if p.grid[vy * w + vx] <= alt {
				isLowPoint = false
				break
			}
		}
		if isLowPoint {
			lowpoints = append(lowpoints, Vertex { x, y })
		}
	}
	return lowpoints
}

func (p Puzzle) PartTwo(lps []Vertex) int {
	var sizes []int
	for _, lp := range lps {
		sizes = append(sizes, p.GetBasinSize(lp))
	}
	slices.Sort(sizes)
	slices.Reverse(sizes)
	return sizes[0] * sizes[1] * sizes[2]
}

func (p Puzzle) GetBasinSize(lp Vertex) int {
	w, h := p.width, p.height
	been := map[Vertex]bool{ lp: true }
	q := []Vertex{ lp }
	for len(q) > 0 {
		s := q[0]
		q = q[1:]
		for _, v := range adjacent {
			adj := Vertex{ s.x + v.x, s.y + v.y }
			if adj.x < 0 || adj.y < 0 || adj.x >= w || adj.y >= h {
				continue
			}
			if p.grid[adj.y * w + adj.x] == 9 {
				continue
			}
			if been[adj] {
				continue
			}
			been[adj] = true
			q = append(q, adj)
		}
	}
	return len(been)
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	lps := p.GetLowPoints()

	var p1 int
	for _, v := range lps {
		p1 += p.grid[v.y * p.width + v.x] + 1
	}
	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p.PartTwo(lps))
}
