package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"slices"
)

type State struct {
	x    int
	cost int
}

type Puzzle []int

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

func abs(n int) int {
	if n > 0 {
		return n
	}
	return -n
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Scan()
	check(s.Err())
	for _, x := range strings.Split(s.Text(), ",") {
		p = append(p, atoi(x))
	}
	slices.Sort(p)
	return p
}

func (p Puzzle) GetCost(align int) (cost int) {
	for _, x := range p {
		cost += abs(align - x)
	}
	return cost
}

func (p Puzzle) TriangleCost(align int) (cost int) {
	for _, x := range p {
		n := abs(align - x)
		cost += (n*n + n) / 2
	}
	return cost
}

func (p Puzzle) BestAlignment(costFunc func(int)int) State {
	x := p[len(p)/2]
	last := State{ x,   costFunc(x)   }
	curr := State{ x+1, costFunc(x+1) }
	step := 1
	if curr.cost > last.cost {
		step = -1
		last, curr = curr, last
	}
	for curr.cost < last.cost {
		last = curr
		curr = State{ last.x + step, costFunc(last.x + step) }
	}
	return last
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	p1 := p.BestAlignment(p.GetCost)
	fmt.Println("Part 1:", p1.cost)

	p2 := p.BestAlignment(p.TriangleCost)
	fmt.Println("Part 2:", p2.cost)
}
