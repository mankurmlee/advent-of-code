package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func fileinput() (out []string) {
	flag.Parse()
	f, err := os.Open(flag.Arg(0))
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		out = append(out, s.Text())
	}
	check(s.Err())
	return out
}

type Vec struct {
	X int
	Y int
}

type Lights struct {
	Size int
	Grid map[Vec]struct{}
}

func newLights() Lights {
	lines := fileinput()
	size := 0
	grid := map[Vec]struct{}{}
	for y, l := range lines {
		size = max(size, len(l))
		for x, v := range []byte(l) {
			if v == '#' {
				grid[Vec{x, y}] = struct{}{}
			}
		}
	}
	return Lights{size, grid}
}

func (l Lights) Update(partTwo bool) Lights {
	out := map[Vec]struct{}{}
	for y := range l.Size {
		for x := range l.Size {
			pos := Vec{x, y}
			_, exists := l.Grid[pos]
			a := l.CountNeighbours(x, y)
			if a == 3 || exists && a == 2 {
				out[pos] = struct{}{}
			}
		}
	}
	if partTwo {
		Corners(out, l.Size-1)
	}
	return Lights{l.Size, out}
}

func Corners(grid map[Vec]struct{}, n int) {
	grid[Vec{0, 0}] = struct{}{}
	grid[Vec{n, 0}] = struct{}{}
	grid[Vec{0, n}] = struct{}{}
	grid[Vec{n, n}] = struct{}{}
}

func (l Lights) CountNeighbours(x, y int) (count int) {
	for dy := -1; dy < 2; dy++ {
		for dx := -1; dx < 2; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			if _, exists := l.Grid[Vec{x + dx, y + dy}]; exists {
				count++
			}
		}
	}
	return count
}

func main() {
	l := newLights()
	for range 100 {
		l = l.Update(false)
	}
	fmt.Println("Part 1:", len(l.Grid))

	l = newLights()
	Corners(l.Grid, l.Size-1)
	for range 100 {
		l = l.Update(true)
	}
	fmt.Println("Part 2:", len(l.Grid))
}
