package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
)

type Vertex struct {
	x, y int
}

type Rect struct {
	lo, hi Vertex
}

type Puzzle struct {
	grove   map[Vertex]bool
}

var (
	dirOrder = []Vertex{
		Vertex{ 0, -1},
		Vertex{ 0,  1},
		Vertex{-1,  0},
		Vertex{ 1,  0},
	}
	dirChecks = map[Vertex][]Vertex{
		Vertex{ 0, -1}: []Vertex{ Vertex{-1, -1}, Vertex{ 0, -1}, Vertex{ 1, -1} },
		Vertex{ 0,  1}: []Vertex{ Vertex{-1,  1}, Vertex{ 0,  1}, Vertex{ 1,  1} },
		Vertex{-1,  0}: []Vertex{ Vertex{-1, -1}, Vertex{-1,  0}, Vertex{-1,  1} },
		Vertex{ 1,  0}: []Vertex{ Vertex{ 1, -1}, Vertex{ 1,  0}, Vertex{ 1,  1} },
	}
	adjacent = []Vertex{
		Vertex{-1, -1}, Vertex{0, -1}, Vertex{1, -1},
		Vertex{-1,  0},                Vertex{1,  0},
		Vertex{-1,  1}, Vertex{0,  1}, Vertex{1,  1},
	}
	numOrders = len(dirOrder)
)

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

	p.grove = make(map[Vertex]bool)

	var y int
	for s.Scan() {
		for x, r := range s.Text() {
			if r != '#' {
				continue
			}
			p.grove[Vertex{x, y}] = true
		}
		y++
	}
	check(s.Err())

	return p
}

func (v1 Vertex) Add(v2 Vertex) Vertex {
	return Vertex{v1.x + v2.x, v1.y + v2.y}
}

func (p Puzzle) IsAlone(elf Vertex) bool {
	for _, adj := range adjacent {
		if p.grove[elf.Add(adj)] {
			return false
		}
	}
	return true
}

func (p Puzzle) GetIntentions(now int) map[Vertex]int {
	intent := make(map[Vertex]int)
	for elf := range p.grove {
		if p.IsAlone(elf) {
			//~ fmt.Println(elf, "wants to stay still")
			intent[elf]++
		} else {
			nextMove := p.NextMove(elf, now)
			intent[nextMove]++
			//~ fmt.Println(elf, "wants to go to", nextMove)
		}
	}
	return intent
}

func (p Puzzle) NextMove(elf Vertex, now int) Vertex {
	nextMove := elf
	for i := 0; i < 4; i++ {
		dirIndex := (now + i) % numOrders
		dir := dirOrder[dirIndex]
		if p.IsDirectionClear(elf, dir) {
			return elf.Add(dir)
		}
	}
	return nextMove
}

func (p Puzzle) IsDirectionClear(elf Vertex, dir Vertex) bool {
	for _, v := range dirChecks[dir] {
		if p.grove[elf.Add(v)] {
			return false
		}
	}
	return true
}

func (p *Puzzle) Update(now int) (stable bool) {
	intent := p.GetIntentions(now)
	nextRound := make(map[Vertex]bool)
	stable = true
	for elf := range p.grove {
		if p.IsAlone(elf) {
			nextRound[elf] = true
		} else {
			nextMove := p.NextMove(elf, now)
			if intent[nextMove] == 1 {
				nextRound[nextMove] = true
			} else {
				nextRound[elf] = true
			}
			stable = false
		}
	}
	p.grove = nextRound
	return stable
}

func (p Puzzle) GetRect() (r Rect) {
	for elf := range p.grove {
		r.lo = elf
		break
	}
	r.hi = r.lo
	for elf := range p.grove {
		if elf.x < r.lo.x {
			r.lo.x = elf.x
		}
		if elf.x > r.hi.x {
			r.hi.x = elf.x
		}
		if elf.y < r.lo.y {
			r.lo.y = elf.y
		}
		if elf.y > r.hi.y {
			r.hi.y = elf.y
		}
	}
	return r
}

func (p Puzzle) CountSpace() (space int) {
	r := p.GetRect()
	for y := r.lo.y; y <= r.hi.y; y++ {
		for x := r.lo.x; x <= r.hi.x; x++ {
			if !p.grove[Vertex{x, y}] {
				space++
			}
		}
	}
	return space
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	t := 0
	for ; t < 10; t++ {
		p.Update(t)
	}

	space := p.CountSpace()
	fmt.Println("Part 1:", space)

	for !p.Update(t) {
		t++
	}
	fmt.Println("Part 2:", t + 1)
}
