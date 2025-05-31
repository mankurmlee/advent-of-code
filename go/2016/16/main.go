package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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

func load(filename string) (out Puzzle) {
	d := readFile(filename)
	out.length = atoi(d[0])
	out.input = make([]uint8, len(d[1]))
	for i, r := range d[1] {
		var v uint8
		if r == '1' {
			v = 1
		}
		out.input[i] = v
	}
	return out
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	p.Check(p.length)
	p.Check(35651584)
}

func extend(in []uint8) []uint8 {
	n := (len(in) << 1) + 1
	out := make([]uint8, n)
	copy(out, in)
	i := n
	for _, r := range in {
		i--
		if r == 0 {
			out[i] = 1
		}
	}
	return out
}

func checksum(s []uint8) []uint8 {
	v := make([]uint8, len(s))
	copy(v, s)
	for len(v)%2 == 0 {
		n := len(v)
		v1 := make([]uint8, n>>1)
		for i := 0; i < n; i += 2 {
			if v[i] == v[i+1] {
				v1[i>>1] = 1
			}
		}
		v = v1
	}
	return v
}

func toString(s []uint8) string {
	var sb strings.Builder
	for _, v := range s {
		if v == 0 {
			sb.WriteRune('0')
		} else {
			sb.WriteRune('1')
		}
	}
	return sb.String()
}

type Puzzle struct {
	length int
	input  []uint8
}

func (p Puzzle) Check(n int) {
	fmt.Println(toString(checksum(p.Fill(n))))
}

func (p Puzzle) Fill(n int) []uint8 {
	s := p.input
	for len(s) < n {
		s = extend(s)
	}
	return s[:n]
}
