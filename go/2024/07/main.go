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

func findInts(s string) (nums []int) {
	r := regexp.MustCompile(`\d+`)
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

func load(filename string) (puzzle [][]int) {
	for _, l := range readFile(filename) {
		puzzle = append(puzzle, findInts(l))
	}
	return puzzle
}

func main() {
	flag.Parse()
	puzzle := load(flag.Arg(0))
	total(puzzle, false)
	total(puzzle, true)
}

func total(puzzle [][]int, concat bool) {
	res := 0
	for _, eq := range puzzle {
		lhs := eq[0]
		rhs := eq[1:]
		if canMake(lhs, rhs, concat) {
			res += lhs
		}
	}
	fmt.Println(res)
}

func canMake(lhs int, rhs []int, concat bool) bool {
	if len(rhs) == 1 {
		return lhs == rhs[0]
	}
	i := len(rhs) - 1
	head := rhs[:i]
	last := rhs[i]
	if last > lhs {
		return false
	}
	if canMake(lhs-last, head, concat) {
		return true
	}
	if lhs%last == 0 && canMake(lhs/last, head, concat) {
		return true
	}
	if concat {
		if rem, ok := unconcat(lhs, last); ok && canMake(rem, head, true) {
			return true
		}
	}
	return false
}

func unconcat(lhs int, last int) (int, bool) {
	dec := 10
	for dec <= last {
		dec *= 10
	}
	if lhs%dec == last {
		return (lhs - last) / dec, true
	}
	return 0, false

}
