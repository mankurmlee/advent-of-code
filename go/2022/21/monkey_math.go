package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Puzzle struct {
	jobs map[string][]string
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
	r := strings.NewReplacer(":", " ")

	p.jobs = make(map[string][]string)
	for s.Scan() {
		data := strings.Fields(r.Replace(s.Text()))
		p.jobs[data[0]] = data[1:]
	}
	check(s.Err())
	return p
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

func sgn(x int) int {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	}
	return 0
}

func (p Puzzle) LazyEval(s string) int {
	data, ok := p.jobs[s]
	if !ok {
		fmt.Println("Puzzle.LazyEval: ERROR:", s, "not found!")
		return 0
	}
	if len(data) != 3 {
		return atoi(data[0])
	}
	a := p.LazyEval(data[0])
	b := p.LazyEval(data[2])
	switch data[1] {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	case "=":
		return a - b
	}
	fmt.Println("Puzzle.LazyEval: ERROR: Operation", data[1], "not recognised!")
	return 0
}

func (p Puzzle) Test(n int) int {
	p.jobs["humn"][0] = strconv.Itoa(n)
	return p.LazyEval("root")
}

func (p Puzzle) FindMatch() (match int) {
	p.jobs["root"][1] = "="
	min := 0
	max := 0
	minres := p.Test(0)
	maxres := minres

	jump := abs(minres)
	for sgn(minres) == sgn(maxres) {
		max += jump
		maxres = p.Test(max)
	}

	for {
		fmt.Println("Solution between", min, "and", max)
		match = (min + max) / 2
		res := p.Test(match)
		if res == 0 {
			break
		}
		if sgn(res) == sgn(minres) {
			min = match
		} else {
			max = match
		}
	}

	return match
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	root := p.LazyEval("root")
	fmt.Println("Part 1:", root)

	match := p.FindMatch()
	fmt.Println("Part 2:", match)
}
