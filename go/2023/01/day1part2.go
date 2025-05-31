package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createLists() (map[string]int, map[string]int) {
	n := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"9": 9,
		"one"  : 1,
		"two"  : 2,
		"three": 3,
		"four" : 4,
		"five" : 5,
		"six"  : 6,
		"seven": 7,
		"eight": 8,
		"nine" : 9,
	}

	r := map[string]int{}
	for k, v := range n {
		r[reverse(k)] = v
	}
	return n, r
}

func reverse(text string) string {
	chars := []rune(text)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

func firstDigit(haystack string, needles map[string]int) int {
	d, f := -1, 999999

	var i int
	for k, v := range needles {
		i = strings.Index(haystack, k)
		if i >= 0 && i < f {
			d, f = v, i
		}
	}

	return d
}

func extractNumber(text string) int {
	first := firstDigit(text, numbers)
	last  := firstDigit(reverse(text), reversed)
	return first * 10 + last
}

func sumFile(file string) int {
	f, err := os.Open(file)
	check(err)
	defer f.Close()

	s := bufio.NewScanner(f)

	var sum int
	for s.Scan() {
		sum += extractNumber(s.Text())
	}

	check(s.Err())
	return sum
}

var numbers, reversed = createLists()

func main() {
	sum := sumFile("input.txt")

	fmt.Println("Sum:", sum)
}


