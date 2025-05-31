package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
)

type Buffer string

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

func (s Buffer) FirstPacketIdx(size int) int {
	charCount := make(map[byte]int)
	var start, end, uniq int

	n := len(s)
	for end < n {
		ch := s[end]

		if charCount[ch] == 0 {
			uniq++
		}
		charCount[ch]++

		for uniq == size {
			if end - start < size {
				return end
			}

			charStart := s[start]
			charCount[charStart]--

			if charCount[charStart] == 0 {
				uniq--
			}

			start++
		}
		end++
	}

	return end
}

func main() {
	flag.Parse()
	b := Buffer(load(flag.Arg(0)))
	fmt.Println(b)

	fmt.Println("Part 1:", b.FirstPacketIdx(4) + 1)
	fmt.Println("Part 2:", b.FirstPacketIdx(14) + 1)
}
