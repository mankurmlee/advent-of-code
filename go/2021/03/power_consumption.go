package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
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

func (p Puzzle) PartOne() int {
	var gamma, epsilon int
	oneCount := make([]int, len(p[0]))
	n := len(p)

	for _, bin := range p {
		for i, r := range bin {
			if r == '1' {
				oneCount[i]++
			}
		}
	}

	for _, ones := range oneCount {
		gamma   <<= 1
		epsilon <<= 1
		zeroes := n - ones
		if ones > zeroes {
			gamma++
		} else {
			epsilon++
		}
	}
	return gamma * epsilon
}

func (p Puzzle) GetRating(a, b byte) (y int) {
	w, h := len(p[0]), len(p)
	valid := make(map[int]struct{}, h)
	cache := strings.Join(p, "")

	for y := 0; y < h; y++ {
		valid[y] = struct{}{}
	}

	for x := 0; len(valid) > 1; x++ {
		ones := 0
		for y := range valid {
			if cache[y * w + x] == '1' {
				ones++
			}
		}
		zeroes := len(valid) - ones

		exclude := a
		if ones >= zeroes {
			exclude = b
		}

		for y := range valid {
			if cache[y * w + x] == exclude {
				delete(valid, y)
			}
		}
	}

	for y = range valid { break }
	rat, err := strconv.ParseInt(p[y], 2, 64)
	check(err)
	return int(rat)
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	fmt.Println("Part 1:", p.PartOne())
	fmt.Println("Part 2:", p.GetRating('1', '0') * p.GetRating('0', '1'))
}
