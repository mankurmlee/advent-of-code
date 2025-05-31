package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"cmp"
)

type Range struct {
	u, v int
}

type Pair struct {
	a, b Range
}

type Pairs []Pair

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func load(filename string) (l Pairs) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	r := strings.NewReplacer("-", " ", ",", " ")

	for s.Scan() {
		data := strings.Fields(r.Replace(s.Text()))
		l = append(l, Pair{
			Range{Atoi(data[0]), Atoi(data[1])},
			Range{Atoi(data[2]), Atoi(data[3])},
		})
	}

	check(s.Err())
	return l
}

func (ps Pairs) CountFullOverlap() (count int) {
	for _, p := range ps {
		l := cmp.Compare(p.a.u, p.b.u)
		h := cmp.Compare(p.a.v, p.b.v)
		if l != h || l == 0 {
			count++
		}
	}
	return count
}

func (ps Pairs) CountOverlap() (count int) {
	for _, p := range ps {
		if  p.a.u < p.b.u && p.a.v < p.b.u ||
			p.a.u > p.b.v && p.a.v > p.b.v {
			continue
		}
		count++
	}
	return count
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	fmt.Println("Part 1:", p.CountFullOverlap())
	fmt.Println("Part 2:", p.CountOverlap())
}
