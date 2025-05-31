package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

var adjacent = []Vertex{
	Vertex{-1, -1}, Vertex{ 0, -1}, Vertex{ 1, -1},
	Vertex{-1,  0},                 Vertex{ 1,  0},
	Vertex{-1,  1}, Vertex{ 0,  1}, Vertex{ 1,  1},
}

type Vertex struct { x, y int }

type Puzzle struct {
	width  int
	height int
	grid   []int
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
		if p.width == 0 {
			p.width = len(line)
		}
		p.height++
		for _, ch := range line {
			p.grid = append(p.grid, int(ch - '0'))
		}
	}
	check(s.Err())
	return p
}

func (p Puzzle) String() string {
	w, h := p.width, p.height
	rows := make([]string, h)
	for y := range rows {
		line := make([]string, w)
		for x := range line {
			line[x] = strconv.Itoa(p.grid[y * w + x])
		}
		rows[y] = strings.Join(line, "")
	}
	return strings.Join(rows, "\n")
}

func (p Puzzle) Update() int {
	var flashes []int
	flashed := map[int]bool{}
	w, h := p.width, p.height

	for i := range p.grid {
		p.grid[i]++
		if p.grid[i] < 10 { continue }
		flashes = append(flashes, i)
		flashed[i] = true
	}

	for len(flashes) > 0 {
		i := flashes[0]
		flashes = flashes[1:]
		x, y := i % w, i / w
		for _, v := range adjacent {
			vx, vy := x + v.x, y + v.y
			if vx < 0 || vy < 0 || vx >= w || vy >= h { continue }
			j := vy * w + vx
			if flashed[j] { continue }
			p.grid[j]++
			if p.grid[j] < 10 { continue }
			flashes = append(flashes, j)
			flashed[j] = true
		}
	}

	for i := range flashed {
		p.grid[i] = 0
	}
	return len(flashed)
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	var p1, i, flashes int
	for i = 0; i < 100; i++ {
		p1 += p.Update()
	}
	fmt.Println("Part 1:", p1)

	for ; flashes < 100; i++ {
		flashes = p.Update()
	}
	fmt.Println("Part 2:", i)
}
