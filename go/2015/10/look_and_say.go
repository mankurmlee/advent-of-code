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

func load(filename string) (num string) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Scan()
	num = s.Text()
	check(s.Err())
	return num
}

func lookandsay(s string) string {
	var (
		out     strings.Builder
		last    rune
		repeats int
	)
	if len(s) == 0 {
		return ""
	}
	for _, r := range s {
		if last == r {
			repeats++
			continue
		}
		if last != 0 {
			num := strconv.Itoa(repeats)
			out.WriteString(num)
			out.WriteRune(last)
		}
		last = r
		repeats = 1
	}
	num := strconv.Itoa(repeats)
	out.WriteString(num)
	out.WriteRune(last)

	return out.String()
}

func main() {
	flag.Parse()
	num := load(flag.Arg(0))
	fmt.Printf("Num %d: %s\n", 0, num)

	for i := 0; i < 40; i++ {
		num = lookandsay(num)
	}
	fmt.Println("Part 1:", len(num))

	for i := 0; i < 10; i++ {
		num = lookandsay(num)
	}
	fmt.Println("Part 2:", len(num))
}
