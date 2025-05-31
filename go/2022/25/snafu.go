package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"slices"
)

var (
	snafuDigit = map[rune]int{ '=': -2, '-': -1, '0': 0, '1': 1, '2': 2 }
	digitSnafu = map[int]string{ -2: "=", -1: "-", 0: "0", 1: "1", 2: "2" }
)

type Puzzle struct {
	snafu []string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Abs(n int) int {
	if n >= 0 {
		return n
	}
	return -n
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		p.snafu = append(p.snafu, s.Text())
	}
	check(s.Err())
	return p
}

func snafu2dec(s string) (dec int) {
	base := 1
	runes := []rune(s)
	slices.Reverse(runes)
	for _, r := range runes {
		dec += base * snafuDigit[r]
		base *= 5
	}
	return dec
}

func dec2snafu(dec int) string {
	if Abs(dec) < 3 {
		return digitSnafu[dec]
	}
	n := dec
	base := 1
	var power int
	for Abs(n) > 2 {
		if n > 0 {
			n += 2
		} else {
			n -= 2
		}
		n /= 5
		base *= 5
		power++
	}
	tail := dec2snafu(dec - n * base)
	for len(tail) < power {
		tail = "0" + tail
	}
	return digitSnafu[n] + tail
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	var sum int
	for _, s := range p.snafu {
		dec := snafu2dec(s)
		sum += dec
		test := dec2snafu(dec)
		fmt.Println(s, "=>", dec, "=>", test)
	}
	fmt.Println("Sum:", sum)
	fmt.Println("Part 1:", dec2snafu(sum))
}
