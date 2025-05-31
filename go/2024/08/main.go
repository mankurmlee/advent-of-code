package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type Vec struct {
	x, y int
}

func (v Vec) Add(o Vec) Vec {
	return Vec{v.x + o.x, v.y + o.y}
}

func (v Vec) Sub(o Vec) Vec {
	return Vec{v.x - o.x, v.y - o.y}
}

type Puzzle struct {
	w, h     int
	antennae map[rune][]Vec
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

func load(filename string) Puzzle {
	a := make(map[rune][]Vec)
	d := readFile(filename)
	for y, l := range d {
		for x, r := range l {
			if r == '.' {
				continue
			}
			a[r] = append(a[r], Vec{x, y})
		}
	}
	return Puzzle{len(d[0]), len(d), a}
}

func main() {
	flag.Parse()
	puzzle := load(flag.Arg(0))

	partOne(puzzle, basic)
	partOne(puzzle, advanced)
}

func partOne(p Puzzle, findNodes func(Puzzle, Vec, Vec) []Vec) {
	nodes := make(map[Vec]bool)
	for _, v := range p.antennae {
		n := len(v)
		for i := range n - 1 {
			for j := i + 1; j < n; j++ {
				for _, node := range findNodes(p, v[i], v[j]) {
					nodes[node] = true
				}
			}
		}
	}
	fmt.Println(len(nodes))
}

func advanced(p Puzzle, a, b Vec) (nodes []Vec) {
	d := b.Sub(a)
	n := a
	for inBounds(p, n) {
		nodes = append(nodes, n)
		n = n.Sub(d)
	}
	n = b
	for inBounds(p, n) {
		nodes = append(nodes, n)
		n = n.Add(d)
	}
	return nodes
}

func basic(p Puzzle, a, b Vec) (nodes []Vec) {
	d := b.Sub(a)
	a_prime := a.Sub(d)
	if inBounds(p, a_prime) {
		nodes = append(nodes, a_prime)
	}
	b_prime := b.Add(d)
	if inBounds(p, b_prime) {
		nodes = append(nodes, b_prime)
	}
	return nodes
}

func inBounds(p Puzzle, v Vec) bool {
	return v.x >= 0 && v.y >= 0 && v.x < p.w && v.y < p.h
}
