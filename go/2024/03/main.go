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

func load(filename string) (lines []string) {
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
	lines := load(flag.Arg(0))
	fmt.Println(partOne(lines))
	fmt.Println(partTwo(lines))
}

func partTwo(lines []string) (tot int) {
	do := true
	r := regexp.MustCompile(`mul\(\d{1,3},\d{1,3}\)|do\(\)|don't\(\)`)
	for _, line := range lines {
		for _, m := range r.FindAllString(line, -1) {
			if m == `do()` {
				do = true
				continue
			} else if m == `don't()` {
				do = false
				continue
			}
			if !do {
				continue
			}
			n := extractNums(m)
			tot += n[0] * n[1]
		}
	}
	return tot
}

func partOne(lines []string) (tot int) {
	for _, line := range lines {
		tot += mulsum(line)
	}
	return tot
}

func mulsum(s string) (tot int) {
	r := regexp.MustCompile(`mul\(\d{1,3},\d{1,3}\)`)
	for _, m := range r.FindAllString(s, -1) {
		n := extractNums(m)
		tot += n[0] * n[1]
	}
	return tot
}

func extractNums(s string) (nums []int) {
	r := regexp.MustCompile(`\d+`)
	for _, m := range r.FindAllString(s, -1) {
		nums = append(nums, atoi(m))
	}
	return nums
}
