package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var directions = []Vec{
	{0, -1}, {1, 0}, {0, 1}, {-1, 0},
}

type Vec struct {
	x, y int
}

func (v Vec) Add(other Vec) Vec {
	return Vec{v.x + other.x, v.y + other.y}
}

type Guard struct {
	pos Vec
	dir int
}

type Puzzle struct {
	startpos Vec
	w, h     int
	walls    map[Vec]bool
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

func load(filename string) (puzzle Puzzle) {
	puzzle.walls = make(map[Vec]bool)
	data := readFile(filename)
	for y, line := range data {
		for x, r := range line {
			switch r {
			case '#':
				puzzle.walls[Vec{x, y}] = true
			case '^':
				puzzle.startpos = Vec{x, y}
			}
		}
	}
	puzzle.h = len(data)
	puzzle.w = len(data[0])
	return puzzle
}

func main() {
	flag.Parse()
	puzzle := load(flag.Arg(0))

	trail := puzzle.PartOne()
	puzzle.PartTwo(trail)
}

func (p Puzzle) PartOne() map[Vec]bool {
	pos := p.startpos
	d := 0
	trail := make(map[Vec]bool)
	for pos.x >= 0 && pos.y >= 0 && pos.x < p.w && pos.y < p.h {
		trail[pos] = true
		next := pos.Add(directions[d])
		if p.walls[next] {
			d = (d + 1) % 4
		} else {
			pos = next
		}
	}
	fmt.Println(len(trail))
	return trail
}

func (p Puzzle) PartTwo(trail map[Vec]bool) {
	count := 0
	for v := range trail {
		if v != p.startpos && p.WillLoop(v) {
			count++
		}
	}
	fmt.Println(count)
}

func (p Puzzle) WillLoop(block Vec) bool {
	walls := make(map[Vec]bool, len(p.walls))
	for v := range p.walls {
		walls[v] = true
	}
	walls[block] = true

	been := make(map[Guard]struct{})
	g := Guard{p.startpos, 0}
	for g.pos.x >= 0 && g.pos.y >= 0 && g.pos.x < p.w && g.pos.y < p.h {
		been[g] = struct{}{}
		pos := g.pos.Add(directions[g.dir])
		for walls[pos] {
			g.dir = (g.dir + 1) % 4
			pos = g.pos.Add(directions[g.dir])
		}
		g.pos = pos

		if _, ok := been[g]; ok {
			return true
		}
	}
	return false
}
