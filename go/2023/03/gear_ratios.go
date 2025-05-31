package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type num struct {
	val  string
	l, r, y int
}

type symbol struct {
	val  rune
	x, y int
}

func (s symbol) GearRatio(nums []num) int {
	var adjacent []num
	var left, right, top, bottom int
	for _, n := range nums {
		left   = n.l - 1
		right  = n.r + 1
		top    = n.y - 1
		bottom = n.y + 1
		if s.x >= left && s.x <= right && s.y >= top && s.y <= bottom {
			adjacent = append(adjacent, n)
		}
	}
	if len(adjacent) != 2 {
		return 0
	}
	a, err := strconv.Atoi(adjacent[0].val)
	check(err)
	b, err := strconv.Atoi(adjacent[1].val)
	check(err)
	return a*b
}

func getParts(nums []num, syms []symbol) (parts []int) {
	var left, right, top, bottom int
	for _, n := range nums {
		left   = n.l - 1
		right  = n.r + 1
		top    = n.y - 1
		bottom = n.y + 1
		for _, s := range syms {
			if s.x >= left && s.x <= right && s.y >= top && s.y <= bottom {
				part, err := strconv.Atoi(n.val)
				check(err)
				parts = append(parts, part)
				break
			}
		}
	}
	return parts
}

func load(filename string) (nums []num, syms []symbol) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for y := 0; s.Scan(); y++ {
		a, b := parseLine(s.Text(), y)
		nums = append(nums, a...)
		syms = append(syms, b...)
	}
	check(s.Err())
	return nums, syms
}

func parseLine(text string, y int) (nums []num, syms []symbol) {
	left := -1
	for x, ch := range text {
		if strings.ContainsRune("1234567890", ch) {
			if left < 0 {
				left = x
			}
		} else {
			if left >= 0 {
				nums = append(nums, num{text[left:x], left, x-1, y})
				left = -1
			}
			if ch != '.' {
				syms = append(syms, symbol{ch, x, y})
			}
		}
	}
	if left >= 0 {
		x := len(text)
		nums = append(nums, num{text[left:], left, x-1, y})
	}
	return nums, syms
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	nums, syms := load(flag.Arg(0))

	var sum int
	for _, p := range getParts(nums, syms) {
		sum += p
	}
	fmt.Println("Part 1 answer is", sum)

	var total int
	for _, s := range syms {
		if s.val != '*' {
			continue
		}
		total += s.GearRatio(nums)
	}
	fmt.Println("Part 1 answer is", total)
}
