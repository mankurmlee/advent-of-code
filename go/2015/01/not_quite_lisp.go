package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load(filename string) string {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Scan()
	check(s.Err())
	return s.Text()
}

func getFloor(s string) (floor int) {
	for _, r := range s {
		if r == '(' {
			floor++
		} else {
			floor--
		}
	}
	return floor
}

func getBasement(s string) int {
	var floor int
	for i, r := range s {
		if r == '(' {
			floor++
		} else {
			floor--
		}
		if floor == -1 {
			return i + 1
		}
	}
	return 0
}

func main() {
	flag.Parse()
	s := load(flag.Arg(0))

	fmt.Println("Part 1:", getFloor(s))
	fmt.Println("Part 2:", getBasement(s))
}
