package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
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

func abs(n int) int {
	if n >= 0 {
		return n
	}
	return -n
}

func load(filename string) (p [][]int) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	p = make([][]int, 2)
	s := bufio.NewScanner(f)
	for s.Scan() {
		d := strings.Fields(s.Text())
		p[0] = append(p[0], atoi(d[0]))
		p[1] = append(p[1], atoi(d[1]))
	}
	check(s.Err())
	return p
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	p1 := partOne(p)
	fmt.Println(p1)

	p2 := partTwo(p)
	fmt.Println(p2)
}

func partTwo(p [][]int) (tot int) {
	cs := counts(p[1])
	for _, v := range p[0] {
		c := cs[v]
		tot += v * c
	}
	return tot
}

func counts(l []int) (cs map[int]int) {
	cs = make(map[int]int)
	for _, v := range l {
		c := cs[v]
		cs[v] = c + 1
	}
	return cs
}

func partOne(p [][]int) (res int) {
	a := p[0]
	b := p[1]
	slices.Sort(a)
	slices.Sort(b)
	for i, v := range a {
		res += abs(b[i] - v)
	}
	return res
}
