package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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

func load(s string) int {
	d := readFile(s)
	return atoi(d[0])
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	partOne(p)
	partTwo(p)
}

func makeTable(n int) []int {
	table := make([]int, n+1)
	for i := 1; i < n; i++ {
		table[i] = i + 1
	}
	table[n] = 1
	return table
}

func partOne(n int) {
	table := makeTable(n)
	i := 1
	j := table[i]
	for i != j {
		k := table[j]
		table[i] = k
		i = k
		j = table[i]
	}
	fmt.Println(i)
}

func partTwo(n int) {
	l := n >> 1
	table := makeTable(n)
	for i := n; i > 1; i-- {
		m := table[l]
		o := table[m]
		table[l] = o
		if i%2 == 1 {
			l = o
		}
	}
	fmt.Println(l)
}
