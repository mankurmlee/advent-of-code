package main

import (
	"bufio"
	"cmp"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
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

func load(s string) (o []Range) {
	for _, l := range readFile(s) {
		d := parseInts(l)
		o = append(o, Range{d[0], d[1]})
	}
	slices.SortFunc(o, func(a, b Range) int {
		return cmp.Compare(a.lo, b.lo)
	})
	return o
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	partOne(p)
	partTwo(p)
}

func partTwo(rs []Range) {
	rs = mergeAll(rs)
	rem := 1 << 32
	for _, r := range rs {
		rem -= 1 + r.hi - r.lo
	}
	fmt.Println(rem)
}

func mergeAll(rs []Range) (o []Range) {
	c := rs[0]
	for _, r := range rs[1:] {
		if r.lo > c.hi {
			o = append(o, c)
			c = r
		}
		c.hi = max(c.hi, r.hi)
	}
	o = append(o, c)
	return o
}

func partOne(rs []Range) {
	pots := make(map[int]struct{})
	for _, r := range rs {
		pots[r.hi+1] = struct{}{}
	}
	for _, r := range rs {
		for c := range pots {
			if r.lo <= c && c <= r.hi {
				delete(pots, c)
			}
		}
	}
	lowest := 1 << 32
	for k := range pots {
		if k < lowest {
			lowest = k
		}
	}
	fmt.Println(lowest)
}

type Range struct {
	lo, hi int
}
