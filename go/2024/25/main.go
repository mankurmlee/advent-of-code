package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readChunkedFile(filename string) (chunks [][]string) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	var chunk []string
	for s.Scan() {
		l := s.Text()
		if l != "" {
			chunk = append(chunk, l)
			continue
		}
		if len(chunk) == 0 {
			continue
		}
		chunks = append(chunks, chunk)
		chunk = []string{}
	}
	check(s.Err())
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks
}

func load(filename string) (p Puzzle) {
	for _, c := range readChunkedFile(filename) {
		var data [5]int
		for _, l := range c[1:6] {
			for i, v := range l {
				if v == '#' {
					data[i]++
				}
			}
		}
		if c[0][0] == '#' {
			p.locks = append(p.locks, data)
		} else {
			p.keys = append(p.keys, data)
		}
	}
	return p
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	p.Solve()
}

type Puzzle struct {
	keys  [][5]int
	locks [][5]int
}

func (p Puzzle) Solve() {
	var tot int
	for _, l := range p.locks {
		for _, k := range p.keys {
			if fits(l, k) {
				tot++
			}
		}
	}
	fmt.Println(tot)
}

func fits(l [5]int, k [5]int) bool {
	for i := range 5 {
		if l[i]+k[i] > 5 {
			return false
		}
	}
	return true
}
