package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Vertex struct {x, y int}

type Line struct { start, end Vertex }

type Puzzle []Line

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func atoi(a string) int {
	i, err := strconv.Atoi(a)
	check(err)
	return i
}

func sgn(n int) int {
	switch {
	case n > 0:
		return 1
	case n < 0:
		return -1
	default:
		return 0
	}
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	r := strings.NewReplacer(",", " ")
	for s.Scan() {
		data := strings.Fields(r.Replace(s.Text()))
		p = append(p, Line{
			Vertex{ atoi(data[0]), atoi(data[1]) },
			Vertex{ atoi(data[3]), atoi(data[4]) },
		})
	}
	check(s.Err())
	return p
}

func (v1 Vertex) GetStep(v2 Vertex) Vertex {
	return Vertex{ sgn(v2.x - v1.x), sgn(v2.y - v1.y) }
}

func (v1 Vertex) Add(v2 Vertex) Vertex {
	return Vertex{ v1.x + v2.x, v1.y + v2.y }
}

func (p Puzzle) CountOverlap(skipDiags bool) (count int) {
	surf := map[Vertex]int{}
	for _, line := range p {
		if skipDiags && line.start.x != line.end.x && line.start.y != line.end.y {
			continue
		}
		pos := line.start
		step := line.start.GetStep(line.end)
		for ; pos != line.end; pos = pos.Add(step) {
			surf[pos]++
		}
		surf[pos]++
	}
	for _, v := range surf {
		if v > 1 {
			count++
		}
	}
	return count
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	fmt.Println("Part 1:", p.CountOverlap(true))
	fmt.Println("Part 2:", p.CountOverlap(false))
}
