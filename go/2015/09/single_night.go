package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type LinkKey struct { a, b string }

type Puzzle struct {
	towns []string
	links map[LinkKey]int
}

type Node struct {
	town string
	cost int
	fuel int
	last *Node
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

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)

	towns := map[string]struct{}{}
	p.links = map[LinkKey]int{}

	for s.Scan() {
		data := strings.Fields(s.Text())
		dist := Atoi(data[4])
		towns[data[0]] = struct{}{}
		towns[data[2]] = struct{}{}
		p.links[LinkKey{ data[0], data[2] }] = dist
		p.links[LinkKey{ data[2], data[0] }] = dist
	}
	check(s.Err())

	for t := range towns {
		p.towns = append(p.towns, t)
	}
	return p
}

func (p Puzzle) FindPaths() (shortest, longest int) {
	var q []Node
	n := len(p.towns)
	for _, t := range p.towns {
		q = append(q, Node{ town: t, fuel: n - 1 })
	}
	shortest = 9999999
	for len(q) > 0 {
		n = len(q)
		s := q[n-1]
		q = q[:n-1]
		if s.fuel == 0 {
			if s.cost < shortest {
				shortest = s.cost
			}
			if s.cost > longest {
				longest = s.cost
			}
		} else {
			q = append(q, p.Next(s)...)
		}
	}
	return shortest, longest
}

func (p Puzzle) Next(s Node) (next []Node) {
	been := map[string]bool{}
	n := len(p.towns) - s.fuel - 1
	c := s
	been[c.town] = true
	for i := 0; i < n; i++ {
		c = *c.last
		been[c.town] = true
	}
	for _, t := range p.towns {
		if been[t] {
			continue
		}
		cost := s.cost + p.links[LinkKey{ s.town, t }]
		next = append(next, Node{ t, cost, s.fuel - 1, &s })
	}
	return next
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	s, l := p.FindPaths()
	fmt.Println("Part 1:", s)
	fmt.Println("Part 2:", l)
}
