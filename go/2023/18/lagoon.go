package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"slices"
	"cmp"
	"runtime"
)

type Vec struct {
	x, y int
}

type Link struct {
	u, v Vec
}

type Lagoon []Link

type Puzzle struct {
	lagoon []Lagoon
}

func (lag Lagoon) GetRowArea(y int) (area int) {
	var linx []Link
	for _, l := range lag {
		if l.u.y <= y && l.v.y >= y {
			linx = append(linx, l)
		}
	}

	var lastx int
	var up, down bool
	for _, l := range linx {
		if l.u.y == l.v.y {
			// Horizontal link
			area += l.v.x - l.u.x - 1
			continue
		}
		if up && down {
			area += l.u.x - lastx - 1
		}
		area++
		if l.u.y < y {
			up = !up
		}
		if l.v.y > y {
			down = !down
		}
		lastx = l.v.x
	}
	return area
}

func (lag Lagoon) GetAreaByRange(minY, maxY int, ch chan int) {
	var area int
	for y := minY; y <= maxY; y++ {
		area += lag.GetRowArea(y)
	}
	ch <- area
}

func (lag Lagoon) GetArea() (area int) {
	minY, maxY := lag.YRange()
	height := maxY - minY + 1

	ncpus := runtime.NumCPU()
	if height < ncpus {
		ncpus = 1
	}

	ch := make(chan int, ncpus)
	defer close(ch)

	for i := 0; i < ncpus; i++ {
		start := minY + i       * height / ncpus
		end   := minY + (i + 1) * height / ncpus - 1
		if i == ncpus - 1 {
			end = maxY
		}
		go lag.GetAreaByRange(start, end, ch)
	}

	for i := 0; i < ncpus; i++ {
		area += <- ch
	}

	return area
}

func (lag Lagoon) YRange() (minY, maxY int) {
	minY = lag[0].u.y
	maxY = minY
	for _, l := range lag {
		if l.u.y < minY {
			minY = l.u.y
		}
		if l.v.y > maxY {
			maxY = l.v.y
		}
	}
	return minY, maxY
}

func (v Vec) Apply(dir string, dist int) (pos Vec, l Link) {
	pos = v
	switch dir {
	case "L":
		pos.x -= dist
		l.u = pos
		l.v = v
	case "R":
		pos.x += dist
		l.u = v
		l.v = pos
	case "U":
		pos.y -= dist
		l.u = pos
		l.v = v
	case "D":
		pos.y += dist
		l.u = v
		l.v = pos
	}
	return pos, l
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	r := strings.NewReplacer("#", " ", "(", " ", ")", " ")

	var pos0, pos1 Vec
	var link0, link1 Link
	var set0, set1 []Link

	for s.Scan() {
		linkData := strings.Fields(r.Replace(s.Text()))

		pos0, link0 = pos0.Apply(linkData[0], Atoi(linkData[1]))
		set0 = append(set0, link0)

		dir1, dist  := decode(linkData[2])
		pos1, link1 = pos1.Apply(dir1, dist)
		set1 = append(set1, link1)
	}
	check(s.Err())

	slices.SortFunc(set0, func(a, b Link) int {
		if n := cmp.Compare(a.u.x, b.u.x); n != 0 {
			return n
		}
		return cmp.Compare(a.v.x, b.v.x)
	})
	slices.SortFunc(set1, func(a, b Link) int {
		if n := cmp.Compare(a.u.x, b.u.x); n != 0 {
			return n
		}
		return cmp.Compare(a.v.x, b.v.x)
	})

	p.lagoon = append(p.lagoon, set0)
	p.lagoon = append(p.lagoon, set1)

	return p
}

func decode(colData string) (dir string, dist int) {
	switch colData[5] {
	case '0':
		dir = "R"
	case '1':
		dir = "D"
	case '2':
		dir = "L"
	case '3':
		dir = "U"
	}
	d, err := strconv.ParseInt(colData[0:5], 16, 64)
	dist = int(d)
	check(err)

	return dir, dist
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
	flag.Parse()
	p := load(flag.Arg(0))

	area0 := p.lagoon[0].GetArea()
	fmt.Println("Part1 answer is:", area0)

	area1 := p.lagoon[1].GetArea()
	fmt.Println("Part2 answer is:", area1)
}
