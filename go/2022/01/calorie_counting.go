package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strconv"
	"slices"
)

type Elf []int

type Elves []Elf

func (p Elves) SumTopN(n int) (sum int) {
	var cals []int
	for _, e := range p {
		cals = append(cals, e.Sum())
	}
	slices.Sort(cals)
	slices.Reverse(cals)
	for i := 0; i < n; i++ {
		sum += cals[i]
	}
	return sum
}

func (e Elf) Sum() (sum int) {
	for _, item := range e {
		sum += item
	}
	return sum
}

func load(filename string) (p Elves) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)

	var e Elf
	for s.Scan() {
		line := s.Text()
		if line == "" {
			p = append(p, e)
			e = Elf{}
			continue
		}
		e = append(e, Atoi(line))
	}
	p = append(p, e)

	check(s.Err())
	return p
}

func Atoi(A string) int {
	i, err := strconv.Atoi(A)
	check(err)
	return i
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	most := p.SumTopN(1)
	fmt.Println("Most calories:", most)

	top3 := p.SumTopN(3)
	fmt.Println("Top 3:", top3)
}
