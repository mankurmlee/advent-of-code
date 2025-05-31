package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strconv"
)

type Puzzle []string

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
	for s.Scan() {
		p = append(p, s.Text())
	}
	check(s.Err())
	return p
}

func main() {
	var p1, p2 int
	flag.Parse()
	p := load(flag.Arg(0))

	for _, s := range p {
		r, err := strconv.Unquote(s)
		check(err)
		n := len(s)
		rn := len(r)
		fmt.Println(s, "has", n, "and", r, "has", rn)
		p1 += n - rn
	}
	fmt.Println("Part 1:", p1)

	for _, s := range p {
		q := strconv.Quote(s)
		n := len(s)
		qn := len(q)
		fmt.Println(s, "has", n, "and", q, "has", qn)
		p2 += qn - n
	}
	fmt.Println("Part 2:", p2)
}
