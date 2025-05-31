package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Seq []int

type Puzzle struct {
	seqs []Seq
}

func (s Seq) Extrapolate() int {
	sub := s.GetSubSequence()
	if sub.IsZero() {
		return s[0]
	}
	return s[len(s)-1] + sub.Extrapolate()
}

func (s Seq) Backpolate() int {
	sub := s.GetSubSequence()
	if sub.IsZero() {
		return s[0]
	}
	return s[0] - sub.Backpolate()
}

func (s Seq) GetSubSequence() (sub Seq) {
	n := len(s)
	for i := 1; i < n; i++ {
		sub = append(sub, s[i] - s[i-1])
	}
	return sub
}

func (s Seq) IsZero() bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		var seq Seq
		for _, numData := range strings.Fields(s.Text()) {
			seq = append(seq, Atoi(numData))
		}
		p.seqs = append(p.seqs, seq)
	}
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

	var sum1 int
	for _, s := range p.seqs {
		sum1 += s.Extrapolate()
	}
	fmt.Println("Part 1 is:", sum1)

	var sum2 int
	for _, s := range p.seqs {
		sum2 += s.Backpolate()
	}
	fmt.Println("Part 2 is:", sum2)
}
