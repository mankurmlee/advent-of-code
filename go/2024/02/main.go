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

func abs(n int) int {
	if n >= 0 {
		return n
	}
	return -n
}

func sgn(n int) int {
	if n > 0 {
		return 1
	} else if n < 0 {
		return -1
	}
	return 0
}

func load(filename string) (seqs [][]int) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		var seq []int
		for _, v := range strings.Fields(s.Text()) {
			seq = append(seq, atoi(v))
		}
		seqs = append(seqs, seq)
	}
	check(s.Err())
	return seqs
}

func main() {
	flag.Parse()
	seqs := load(flag.Arg(0))
	fmt.Println(seqs)

	fmt.Println(partOne(seqs))
	fmt.Println(partTwo(seqs))
}

func partTwo(seqs [][]int) (n int) {
	for _, seq := range seqs {
		if safe(seq) || almostSafe(seq) {
			n++
		}
	}
	return n
}

func almostSafe(seq []int) bool {
	for i := range seq {
		var s []int
		for j, v := range seq {
			if i == j {
				continue
			}
			s = append(s, v)
		}
		if safe(s) {
			return true
		}
	}
	return false
}

func partOne(seqs [][]int) (n int) {
	for _, seq := range seqs {
		if safe(seq) {
			n++
		}
	}
	return n
}

func safe(seq []int) bool {
	n := len(seq)
	if n <= 1 {
		return true
	}
	x0 := seq[0]
	j := 0
	for i := 1; i < n; i++ {
		x := seq[i]
		diff := x - x0
		j += sgn(diff)
		if abs(j) != i {
			return false
		}
		d := abs(diff)
		if d < 1 || d > 3 {
			return false
		}
		x0 = x
	}
	return true
}
