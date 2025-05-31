package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"slices"
)

type Pair [2]byte

type Polymer map[Pair]int

type Puzzle struct {
	start  Polymer
	rules  map[Pair]byte
	tail   byte
}

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
	p.rules = map[Pair]byte{}
	for s.Scan() {
		data := strings.Fields(s.Text())
		switch len(data) {
		case 1:
			p.start = createPolymer(data[0])
			n := len(data[0])
			p.tail = data[0][n-1]
		case 3:
			k := data[0]
			v := data[2]
			p.rules[Pair{k[0], k[1]}] = v[0]
		}
	}
	check(s.Err())
	return p
}

func createPolymer(s string) Polymer {
	p := Polymer{}
	n := len(s)
	for i := 0; i < n - 1; i++ {
		p[Pair{s[i], s[i+1]}]++
	}
	return p
}

func (cfg Puzzle) Expand(in Polymer) Polymer {
	out := Polymer{}
	for p, c := range in {
		m := cfg.rules[p]
		out[Pair{ p[0], m }] += c
		out[Pair{ m, p[1] }] += c
	}
	return out
}

func (cfg Puzzle) GetAnswer(p Polymer) int {
	c := map[byte]int{}
	for k, v := range p {
		c[k[0]] += v
	}
	c[cfg.tail]++

	n := len(c)
	counts := make([]int, n)
	var i int
	for _, v := range c {
		counts[i] = v
		i++
	}
	slices.Sort(counts)
	return counts[n-1] - counts[0]
}

func main() {
	flag.Parse()
	cfg := load(flag.Arg(0))

	p := cfg.start
	for i := 0; i < 10; i++ {
		p = cfg.Expand(p)
	}
	fmt.Println("Part 1:", cfg.GetAnswer(p))

	for i := 0; i < 30; i++ {
		p = cfg.Expand(p)
	}
	fmt.Println("Part 2:", cfg.GetAnswer(p))
}
