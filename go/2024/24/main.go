package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func findWords(s string) (out []string) {
	re := regexp.MustCompile(`\w+`)
	return re.FindAllString(s, -1)
}

func readChunkedFile(filename string) (chunks [][]string) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	var chunk []string
	for s.Scan() {
		l := s.Text()
		if l != "" {
			chunk = append(chunk, l)
			continue
		}
		if len(chunk) == 0 {
			continue
		}
		chunks = append(chunks, chunk)
		chunk = []string{}
	}
	check(s.Err())
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks
}

func load(filename string) (p Puzzle) {
	c := readChunkedFile(filename)
	for _, l := range c[0] {
		d := findWords(l)
		p.wires = append(p.wires, Wire{d[0], d[1] == "1"})
	}
	for _, l := range c[1] {
		d := findWords(l)
		p.gates = append(p.gates, Gate{d[1], d[0], d[2], d[3]})
	}
	return p
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	c := p.ToCircuit()
	c.PartOne()
}

type Wire struct {
	name  string
	value bool
}

type Gate struct {
	logic     string
	i0, i1, o string
}

func (g Gate) Print() {
	fmt.Println(g.i0, g.logic, g.i1, "->", g.o)
}

type Puzzle struct {
	wires []Wire
	gates []Gate
}

func (p Puzzle) ToCircuit() (c Circuit) {
	c.wires = make(map[string]bool)
	for _, w := range p.wires {
		c.wires[w.name] = w.value
	}
	c.gates = make(map[string]Gate)
	for _, g := range p.gates {
		c.gates[g.o] = g
	}
	return c
}

type Circuit struct {
	wires map[string]bool
	gates map[string]Gate
}

func (c Circuit) PartOne() {
	re := regexp.MustCompile(`^[xyz]\d\d$`)
	var z []string
	for _, g := range c.gates {
		if g.o[0] == 'z' {
			z = append(z, g.o)
		}
		if g.logic == "XOR" {
			if !re.MatchString(g.i0) && !re.MatchString(g.i1) && !re.MatchString(g.o) {
				g.Print()
			}
		}
	}
	slices.Sort(z)
	var tot int
	for i, n := range z {
		if c.Eval(n) {
			tot += 1 << i
		}
		c.Explain(n)
	}
	fmt.Println(tot)
}

func (c Circuit) Explain(n string) {
	if n == "z45" {
		return
	}
	g := c.gates[n]
	if g.logic != "XOR" {
		g.Print()
		return
	}
	i0 := c.gates[g.i0]
	i1 := c.gates[g.i1]
	if i0.logic == "AND" && i0.i1 != "x00" {
		i0.Print()
	}
	if i1.logic == "AND" && i0.i1 != "x00" {
		i1.Print()
	}
}

func (c Circuit) Eval(n string) (v bool) {
	if b, ok := c.wires[n]; ok {
		return b
	}
	g := c.gates[n]
	a := c.Eval(g.i0)
	b := c.Eval(g.i1)
	switch g.logic {
	case "AND":
		v = a && b
	case "OR":
		v = a || b
	case "XOR":
		v = (a || b) && !(a && b)
	}
	c.wires[n] = v
	return v
}
