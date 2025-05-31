package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strconv"
)

type Puzzle struct {
	measurements []int
}

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

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		p.measurements = append(p.measurements, atoi(s.Text()))
	}
	check(s.Err())
	return p
}

func (p Puzzle) CountSlidingIncr(window int) (num int) {
	var (
		last int
		curr int
		h int
	)

	for h = 0; h < window; h++ {
		curr += p.measurements[h]
	}

	for n := len(p.measurements); h < n; h++ {
		last = curr
		curr += p.measurements[h] - p.measurements[h - window]
		if curr > last {
			num++
		}
	}
	return num
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	fmt.Println("Part 1:", p.CountSlidingIncr(1))
	fmt.Println("Part 2:", p.CountSlidingIncr(3))
}
