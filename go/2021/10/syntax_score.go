package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"slices"
)

var (
	pairs   = map[rune]rune{ '(': ')', '[': ']', '{': '}', '<': '>' }
	points  = map[rune]int{
		')': 3, ']': 57, '}': 1197, '>': 25137,
		'(': 1, '[': 2, '{': 3, '<': 4,
	}
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

func (p Puzzle) PartOne() (score int) {
	for _, line := range p {
		score += corruptScore(line)
	}
	return score
}

func corruptScore(line string) int {
	var stack []rune
	for _, ch := range line {
		if pairs[ch] > 0 {
			stack = append(stack, ch)
			continue
		}
		n := len(stack)
		o := stack[n-1]
		stack = stack[:n-1]
		expected := pairs[o]
		if ch != expected {
			fmt.Printf("%s - Expected %s, but found %s instead.\n",
				line, string(expected), string(ch))
			return points[ch]
		}
	}
	return 0
}

func (p Puzzle) PartTwo() int {
	var scores []int
	for _, line := range p {
		score := autocomplete(line)
		if score > 0 {
			scores = append(scores, score)
		}
	}
	slices.Sort(scores)
	n := len(scores)
	return scores[n/2]
}

func autocomplete(line string) (score int) {
	var stack []rune
	for _, ch := range line {
		if pairs[ch] > 0 {
			stack = append(stack, ch)
			continue
		}
		n := len(stack)
		o := stack[n-1]
		stack = stack[:n-1]
		expected := pairs[o]
		if ch != expected {
			return 0
		}
	}
	slices.Reverse(stack)
	for _, ch := range stack {
		score = score * 5 + points[ch]
	}
	return score
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	fmt.Println("Part 1:", p.PartOne())
	fmt.Println("Part 2:", p.PartTwo())
}
