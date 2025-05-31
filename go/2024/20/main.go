package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var DIR = [4]Vec{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
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

func load(filename string) Puzzle {
	d := readFile(filename)
	w := make(map[Vec]bool)
	var start Vec
	var end Vec
	for y, l := range d {
		for x, v := range l {
			switch v {
			case '#':
				w[Vec{x, y}] = true
			case 'S':
				start = Vec{x, y}
			case 'E':
				end = Vec{x, y}
			}
		}
	}
	p := Puzzle{
		len(d[0]), len(d),
		start, end,
		w,
		nil,
		nil,
	}
	p.costTo = p.PriceMap(start)
	p.costFrom = p.PriceMap(end)
	return p
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	p.Checksum(2)
	p.Checksum(20)
}

type Vec struct {
	x, y int
}

func (v Vec) TaxiTo(o Vec) int {
	return abs(o.x-v.x) + abs(o.y-v.y)
}

func (v Vec) Add(o Vec) Vec {
	return Vec{v.x + o.x, v.y + o.y}
}

func (v Vec) Equals(o Vec) bool {
	return v.x == o.x && v.y == o.y
}

type State struct {
	Vec
	cost int
}

type Puzzle struct {
	w, h             int
	start, end       Vec
	walls            map[Vec]bool
	costTo, costFrom map[Vec]int
}

func (p Puzzle) Checksum(cheat int) {
	tot := 0
	for k, v := range p.Save(cheat) {
		if k >= 100 {
			tot += v
		}
	}
	fmt.Println(tot)
}

func (p Puzzle) Save(cheat int) map[int]int {
	saved := make(map[int]int)
	base := p.costTo[p.end]
	for s, c1 := range p.costTo {
		if c1+2 >= base {
			continue
		}
		for e, c2 := range p.costFrom {
			if c1+c2+2 >= base {
				continue
			}
			t := s.TaxiTo(e)
			if t > cheat {
				continue
			}
			c := c1 + c2 + t
			if c >= base {
				continue
			}
			saved[base-c]++
		}
	}
	return saved
}

func (p Puzzle) PriceMap(start Vec) map[Vec]int {
	best := make(map[Vec]int)
	best[start] = 0
	q := []State{{start, 0}}
	for len(q) > 0 {
		s := q[0]
		q = q[1:]
		for _, d := range DIR {
			s1 := State{s.Vec.Add(d), s.cost + 1}
			if p.walls[s1.Vec] {
				continue
			}
			b, ok := best[s1.Vec]
			if ok && b <= s1.cost {
				continue
			}
			best[s1.Vec] = s1.cost
			q = append(q, s1)
		}
	}
	return best
}
