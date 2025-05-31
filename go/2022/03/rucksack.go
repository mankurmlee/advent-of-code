package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
)

type Rucksack string

type Puzzle []Rucksack

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		p = append(p, Rucksack(s.Text()))
	}
	check(s.Err())
	return p
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func priority(ch rune) int {
	char := byte(ch)
	if char >= 'a' && char <= 'z' {
		return int(char - 'a' + 1)
	}
	if char >= 'A' && char <= 'Z' {
		return int(char - 'A' + 27)
	}
	return 0
}

func (s Rucksack) GetMatch() (char rune) {
	items := make(map[rune]struct{})

	mid := len(s) / 2
	for _, ch := range s[:mid] {
		items[ch] = struct{}{}
	}

	for _, ch := range s[mid:] {
		if _, exists := items[ch]; exists {
			return ch
		}
	}

	return '0'
}

func findCommon(a, b, c Rucksack) (char rune) {
	inA := make(map[rune]struct{})
	for _, ch := range a {
		inA[ch] = struct{}{}
	}

	andB := make(map[rune]struct{})
	for _, ch := range b {
		if _, exists := inA[ch]; exists {
			andB[ch] = struct{}{}
		}
	}

	for _, ch := range c {
		if _, exists := andB[ch]; exists {
			return ch
		}
	}

	return '0'
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	var sum int
	for _, r := range p {
		c := r.GetMatch()
		sum += priority(c)
	}
	fmt.Println("Part 1:", sum)

	sum = 0
	n := len(p)
	for i := 0; i < n; i += 3 {
		c := findCommon(p[i], p[i+1], p[i+2])
		sum += priority(c)
	}
	fmt.Println("Part 2:", sum)
}
