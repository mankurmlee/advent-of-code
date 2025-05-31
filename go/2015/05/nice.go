package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
)

type Puzzle []string

var (
	isVowel = map[rune]bool{ 'a': true, 'e': true, 'i': true, 'o': true, 'u': true }
	badStrings = []string{"ab", "cd", "pq", "xy"}
)


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

func niceTest1(s string) bool {
	var (
		numVowels int
		hasDouble bool
		last      rune
	)

	for _, bad := range badStrings {
		if strings.Contains(s, bad) {
			return false
		}
	}

	for _, r := range s {
		if isVowel[r] {
			numVowels++
		}
		hasDouble = hasDouble || r == last
		last = r
	}

	return hasDouble && numVowels >= 3
}

func niceTest2(s string) bool {
	var (
		last, lastlast rune
		hasRepeatedPair, hasXYX bool
		pair string
	)

	pairs := map[string]int{}

	for i, r := range s {
		pair = string(last) + string(r)
		j, ok := pairs[pair]
		if ok {
			hasRepeatedPair = hasRepeatedPair || i > j + 1
		} else {
			pairs[pair] = i
		}
		hasXYX = hasXYX || r == lastlast
		lastlast = last
		last = r
	}

	return hasRepeatedPair && hasXYX
}

func main() {
	var p1, p2 int

	flag.Parse()
	p := load(flag.Arg(0))

	for _, s := range p {
		if niceTest1(s) {
			fmt.Println(s, "passes part 1 test")
			p1++
		}
		if niceTest2(s) {
			fmt.Println(s, "passes part 2 test")
			p2++
		}
	}

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)
}
