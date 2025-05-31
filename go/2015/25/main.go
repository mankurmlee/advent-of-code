package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

func fileinput() (out string) {
	flag.Parse()
	f, err := os.Open(flag.Arg(0))
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Scan()
	check(s.Err())
	return s.Text()
}

func sequence() []int {
	seq := make([]int, 0, Mod)
	i := Start
	for {
		seq = append(seq, i)
		i = (i * Mul) % Mod
		if i == Start {
			return seq
		}
	}
}

func getInput() int {
	re := regexp.MustCompile("[,.]")
	line := re.ReplaceAllLiteralString(fileinput(), "")
	data := strings.Fields(line)
	n := len(data)
	row := atoi(data[n-3])
	col := atoi(data[n-1])
	return calcIndex(row, col)
}

func calcIndex(row, col int) int {
	n := row + col - 1
	return (n*n+n)>>1 - row
}

const (
	Start = 20151125
	Mul   = 252533
	Mod   = 33554393
)

func main() {
	n := getInput()
	s := sequence()
	fmt.Println("Part 1:", s[n%len(s)])
}
