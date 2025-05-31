package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
)

var directions = []Vec{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

type Vec struct {
	x, y int
}

type Grid struct {
	w, h int
	data map[Vec]rune
}

type Region []Vec

func (r Region) GetFences() [][]Vec {
	g := make(map[Vec]bool)
	for _, v := range r {
		g[v] = true
	}

	f := make([][]Vec, 4)
	for i, d := range directions {
		for _, v := range r {
			v1 := Vec{v.x + d.x, v.y + d.y}
			if !g[v1] {
				f[i] = append(f[i], v)
			}
		}
	}
	return f
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

func load(filename string) Grid {
	f := readFile(filename)
	data := make(map[Vec]rune)
	for y, l := range f {
		for x, r := range l {
			data[Vec{x, y}] = r
		}
	}
	return Grid{len(f[0]), len(f), data}
}

func (g Grid) GetRegions() (regions []Region) {
	seen := make(map[Vec]bool)
	for v := range g.data {
		if seen[v] {
			continue
		}
		r := g.GetRegion(v)
		for _, v1 := range r {
			seen[v1] = true
		}
		regions = append(regions, r)
	}
	return regions
}

func (g Grid) GetRegion(start Vec) Region {
	seen := make(map[Vec]bool)
	plant := g.data[start]
	seen[start] = true
	r := Region{start}
	q := []Vec{start}
	for len(q) > 0 {
		end := len(q) - 1
		v := q[end]
		q = q[:end]
		for _, d := range directions {
			v1 := Vec{v.x + d.x, v.y + d.y}
			if seen[v1] {
				continue
			}
			seen[v1] = true
			if g.data[v1] != plant {
				continue
			}
			r = append(r, v1)
			q = append(q, v1)
		}
	}
	return r
}

func countSides(fences [][]Vec) (out int) {
	for i, d := range fences {
		out += countDirection(d, i%2 == 0)
	}
	return out
}

func countDirection(d []Vec, ns bool) (sides int) {
	m := make(map[int][]int)
	var u int
	var v int
	for _, f := range d {
		if ns {
			u, v = f.x, f.y
		} else {
			u, v = f.y, f.x
		}
		m[v] = append(m[v], u)
	}
	for _, l := range m {
		slices.Sort(l)
		last := 0
		for i, v := range l {
			if i == 0 || last+1 != v {
				sides++
			}
			last = v
		}
	}
	return sides
}

func countFences(fences [][]Vec) (out int) {
	for _, f := range fences {
		out += len(f)
	}
	return out
}

func main() {
	flag.Parse()
	grid := load(flag.Arg(0))

	var p1 int
	var p2 int
	for _, r := range grid.GetRegions() {
		area := len(r)
		fences := r.GetFences()
		p1 += area * countFences(fences)
		p2 += area * countSides(fences)
	}
	fmt.Println(p1)
	fmt.Println(p2)
}
