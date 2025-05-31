package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
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

func parseInts(s string) (nums []int) {
	r := regexp.MustCompile(`\-?\d+`)
	for _, v := range r.FindAllString(s, -1) {
		nums = append(nums, atoi(v))
	}
	return nums
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

func main() {
	flag.Parse()
	discs := load(flag.Arg(0))
	solve(discs)
	n := len(discs)
	discs = append(discs, Disc{n + 1, 11, 0})
	solve(discs)
}

func solve(discs []Disc) {
	prod := 1
	var x int
	for _, d := range discs {
		for (x+d.level+d.pos)%d.slots != 0 {
			x += prod
		}
		prod *= d.slots
	}
	fmt.Println(x)
}

func load(filename string) (out []Disc) {
	for _, l := range readFile(filename) {
		nums := parseInts(l)
		out = append(out, Disc{nums[0], nums[1], nums[3]})
	}
	return out
}

type Disc struct {
	level int
	slots int
	pos   int
}
