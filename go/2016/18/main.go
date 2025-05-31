package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

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

func load(s string) (p Puzzle) {
	d := readFile(s)
	p.rows = atoi(d[0])
	p.first = []rune(d[1])
	return p
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	p.Solve(p.rows)
	p.Solve(400000)
}

type Puzzle struct {
	rows  int
	first []rune
}

func (p Puzzle) Solve(rows int) {
	var safecount int
	s := make([]rune, len(p.first))
	copy(s, p.first)
	for range rows - 1 {
		safecount += count(s)
		s1 := nextrow(s)
		s = s1
	}
	safecount += count(s)
	fmt.Println(safecount)
}

func count(s []rune) (n int) {
	for _, v := range s {
		if v == '.' {
			n++
		}
	}
	return n
}

func nextrow(s []rune) []rune {
	n := len(s)
	o := make([]rune, n)
	for i := range s {
		left := '.'
		if i >= 1 {
			left = s[i-1]
		}
		right := '.'
		if i+1 < n {
			right = s[i+1]
		}
		if left == right {
			o[i] = '.'
		} else {
			o[i] = '^'
		}
	}
	return o
}
