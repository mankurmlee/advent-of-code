package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func fileinput() (out []string) {
	flag.Parse()
	f, err := os.Open(flag.Arg(0))
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		out = append(out, s.Text())
	}
	check(s.Err())
	return out
}

type FactorIndex struct {
	cache map[int]map[int]struct{}
}

func newFactorIndex() FactorIndex {
	return FactorIndex{map[int]map[int]struct{}{
		1: {1: struct{}{}},
	}}
}

func (z FactorIndex) PartOne(target int) int {
	t := target / 10
	i := 1
	sum := 0
	for {
		sum = 0
		for k := range z.Get(i) {
			sum += k
		}
		if sum >= t {
			return i
		}
		i++
	}
}

func (z FactorIndex) PartTwo(target, houses, scale int) int {
	t := int(math.Ceil(float64(target) / float64(scale)))
	i := 1
	sum := 0
	for {
		sum = 0
		for k := range z.Get(i) {
			if i <= k*houses {
				sum += k
			}
		}
		if sum >= t {
			return i
		}
		i++
	}
}

func (z FactorIndex) Get(n int) map[int]struct{} {
	if cached, e := z.cache[n]; e {
		return cached
	}
	p := nextFactor(n)
	lu := z.Get(n / p)
	out := cloneSet(lu)
	for k := range lu {
		out[k*p] = struct{}{}
	}
	z.cache[n] = out
	return out
}

func nextFactor(n int) int {
	if n%2 == 0 {
		return 2
	}
	for i := 3; i < int(math.Sqrt(float64(n)))+1; i += 2 {
		if n%i == 0 {
			return i
		}
	}
	return n
}

func cloneSet(inp map[int]struct{}) map[int]struct{} {
	out := make(map[int]struct{})
	for k := range inp {
		out[k] = struct{}{}
	}
	return out
}

func main() {
	target := atoi(fileinput()[0])
	s := newFactorIndex()
	fmt.Println("Part 1:", s.PartOne(target))
	fmt.Println("Part 2:", s.PartTwo(target, 50, 11))
}
