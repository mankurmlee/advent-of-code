package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Vec struct {
	x, y int
}

type Puzzle [][]Vec

type Surface struct {
	grid []bool
	pos  Vec
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func Sgn(n int) int {
	if n > 0 {
		return 1
	} else if n < 0 {
		return -1
	}
	return 0
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	r := strings.NewReplacer("->", "", ",", " ")
	for s.Scan() {
		data := strings.Fields(r.Replace(s.Text()))
		n := len(data)
		var path []Vec
		for i := 0; i < n; i += 2 {
			path = append(path, Vec{Atoi(data[i]), Atoi(data[i+1])})
		}
		p = append(p, path)
	}
	check(s.Err())
	return p
}

func (q Puzzle) Draw(s *Surface) {
	for _, p := range q {
		for i, v := range p {
			if i == 0 {
				s.Move(v)
			} else {
				s.Draw(v)
			}
		}
	}
}

func (q Puzzle) DrawFloor(s *Surface) {
	var floor int
	for _, p := range q {
		for _, v := range p {
			if v.y > floor {
				floor = v.y
			}
		}
	}

	floor += 2
	s.Move(Vec{0, floor})
	s.Draw(Vec{999, floor})
}

func (v Vec) Add(v1 Vec) Vec {
	return Vec{ v1.x + v.x, v1.y + v.y }
}

func (v Vec) Sub(v1 Vec) Vec {
	return Vec{ v.x - v1.x, v.y - v1.y }
}

func (v Vec) Sgn() Vec {
	return Vec{ Sgn(v.x), Sgn(v.y) }
}

func (s *Surface) Move(target Vec) {
	s.pos = target
}

func (s *Surface) Draw(target Vec) {
	step := target.Sub(s.pos).Sgn()
	for curr := s.pos; curr != target; curr = curr.Add(step) {
		s.grid[curr.y * 1000 + curr.x] = true
	}
	s.grid[target.y * 1000 + target.x] = true
	s.pos = target
}

func (s *Surface) DropSand(src Vec) Vec {
	x, y := src.x, src.y
	for ; y < 199; y++ {
		if !s.grid[(y+1) * 1000 + x] {
			continue
		}
		if !s.grid[(y+1) * 1000 + x-1] {
			x--
			continue
		}
		if !s.grid[(y+1) * 1000 + x+1] {
			x++
			continue
		}
		s.grid[y * 1000 + x] = true
		return Vec{x, y}
	}
	return Vec{x, y}
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	s := Surface{ grid: make([]bool, 200000) }
	p.Draw(&s)
	var grains int
	for {
		v := s.DropSand(Vec{500, 0})
		if v.y > 190 {
			break
		}
		grains++
	}
	fmt.Println("Part 1:", grains)

	s = Surface{ grid: make([]bool, 200000) }
	p.Draw(&s)
	p.DrawFloor(&s)
	grains = 0
	for {
		v := s.DropSand(Vec{500, 0})
		grains++
		if v.y == 0 {
			break
		}
	}
	fmt.Println("Part 2:", grains)
}
