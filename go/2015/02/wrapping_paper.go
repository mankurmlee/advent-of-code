package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Present struct { l, w, h int }

type Puzzle []Present

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
	for s.Scan() {
		data := strings.Split(s.Text(), "x")
		p = append(p, Present{ Atoi(data[0]), Atoi(data[1]), Atoi(data[2]) })
	}
	check(s.Err())
	return p
}

func (p Present) WrappingArea() int {
	a := p.l * p.w
	b := p.w * p.h
	c := p.h * p.l
	return 2 * (a + b + c) + min(a, b, c)
}

func (p Present) RibbonLength() int {
	return 2 * (p.l + p.w + p.h - max(p.l, p.w, p.h)) + p.l * p.w * p.h
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	var total int
	for _, item := range p {
		area := item.WrappingArea()
		fmt.Println(item, "=", area)
		total += area
	}
	fmt.Println("Part 1:", total)

	var totallength int
	for _, item := range p {
		totallength += item.RibbonLength()
	}
	fmt.Println("Part 2:", totallength)
}
