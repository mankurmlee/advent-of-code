package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Vertex struct {x, y int}

type Command struct {
	op string
	n  int
}

type Course []Command

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

func load(filename string) (p Course) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		data := strings.Fields(s.Text())
		p = append(p, Command{ data[0], atoi(data[1]) })
	}
	check(s.Err())
	return p
}

func (p Course) PartOne() int {
	var pos Vertex
	for _, cmd := range p {
		switch cmd.op {
		case "forward":
			pos.x += cmd.n
		case "down":
			pos.y += cmd.n
		case "up":
			pos.y -= cmd.n
		}
	}
	return pos.x * pos.y
}

func (p Course) PartTwo() int {
	var (
		pos Vertex
		aim int
	)
	for _, cmd := range p {
		switch cmd.op {
		case "forward":
			pos.x += cmd.n
			pos.y += aim * cmd.n
		case "down":
			aim += cmd.n
		case "up":
			aim -= cmd.n
		}
	}
	return pos.x * pos.y
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	fmt.Println("Part 1:", p.PartOne())
	fmt.Println("Part 2:", p.PartTwo())
}
