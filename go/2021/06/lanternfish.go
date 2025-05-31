package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
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

func load(filename string) (state []int) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Scan()
	check(s.Err())
	for _, num := range strings.Split(s.Text(), ",") {
		state = append(state, atoi(num))
	}
	return state
}

func getCount(initial []int, days int) (count int) {
	timer := make([]int, 9)
	for _, t := range initial {
		timer[t]++
	}
	for i := 0; i < days; i++ {
		timer = nextState(timer)
	}
	for _, c := range timer {
		count += c
	}
	return count
}

func nextState(now []int) []int {
	next := now[1:]
	next[6] += now[0]
	return append(next, now[0])
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	fmt.Println("Part 1:", getCount(p, 80))
	fmt.Println("Part 2:", getCount(p, 256))
}
