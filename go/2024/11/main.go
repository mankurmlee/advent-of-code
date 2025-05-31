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

func change(v int) (res []int) {
	if v == 0 {
		return []int{1}
	}
	s := strconv.Itoa(v)
	n := len(s)
	if n%2 == 0 {
		m := n >> 1
		return []int{atoi(s[:m]), atoi(s[m:])}
	}
	return []int{v * 2024}
}

func count(puzzle []int, n int) {
	d := dict(puzzle)
	for range n {
		d = blink(d)
	}
	sum := 0
	for _, c := range d {
		sum += c
	}
	fmt.Println(sum)
}

func dict(puzzle []int) map[int]int {
	d := make(map[int]int)
	for _, v := range puzzle {
		d[v]++
	}
	return d
}

func blink(puzzle map[int]int) map[int]int {
	res := make(map[int]int)
	for v, c := range puzzle {
		for _, v1 := range change(v) {
			res[v1] += c
		}
	}
	return res
}

func main() {
	flag.Parse()
	puzzle := findInts(readFile(flag.Arg(0))[0])
	count(puzzle, 25)
	count(puzzle, 75)
}
