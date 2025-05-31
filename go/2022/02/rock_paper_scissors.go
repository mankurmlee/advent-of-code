package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
)

type Round struct {
	oppo string
	resp string
}

type Puzzle []Round

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
		data := strings.Fields(s.Text())
		r := Round{ data[0], data[1] }
		p = append(p, r)
	}
	check(s.Err())
	return p
}

func (p Puzzle) Part1() (score int) {
	for _, r := range p {
		score += r.ShapeScore()
		score += r.OutcomeScore()
	}
	return score
}

func (r Round) ShapeScore() int {
	switch r.resp {
	case "X":
		return 1
	case "Y":
		return 2
	case "Z":
		return 3
	}
	return 0
}

func (r Round) OutcomeScore() int {
	if  r.oppo == "A" && r.resp == "X" ||
		r.oppo == "B" && r.resp == "Y" ||
		r.oppo == "C" && r.resp == "Z" {
		return 3
	}
	if  r.oppo == "A" && r.resp == "Y" ||
		r.oppo == "B" && r.resp == "Z" ||
		r.oppo == "C" && r.resp == "X" {
		return 6
	}
	return 0
}

func (p Puzzle) Part2() (score int) {
	for _, r := range p {
		score += r.Response()
	}
	return score
}

func oppoVal(shape string) int {
	switch shape {
	case "A":
		return 1
	case "B":
		return 2
	case "C":
		return 3
	}
	return 0
}

func (r Round) Response() (score int) {
	var me int
	oppo := oppoVal(r.oppo)
	switch r.resp {
	case "X":
		// lose
		me = oppo - 1
		if me <= 0 {
			me = 3
		}
	case "Y":
		// draw
		score += 3
		me = oppo
	case "Z":
		// win
		score += 6
		me = oppo + 1
		if me > 3 {
			me = 1
		}
	}
	score += me
	return score
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	fmt.Println("Part 1:", p.Part1())
	fmt.Println("Part 2:", p.Part2())
}
