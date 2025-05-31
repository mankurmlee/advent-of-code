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

var REV = []int{1, 1, 6, 2, 7, 3, 0, 4}

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

func parseWords(s string) []string {
	re := regexp.MustCompile(`\w+`)
	return re.FindAllString(s, -1)
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

func load(s string) (p Puzzle) {
	ls := readFile(s)
	p.pass = ls[0]
	p.hash = ls[1]
	for _, l := range ls[2:] {
		p.stmts = append(p.stmts, parseWords(l))
	}
	return p
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	p.Scramble()
	p.Unscramble()
}

type Puzzle struct {
	pass  string
	hash  string
	stmts [][]string
}

func (p Puzzle) Scramble() {
	pass := NewBuffer(p.pass)
	for _, stmt := range p.stmts {
		pass.Exec(stmt, false)
	}
	pass.Println()
}

func (p Puzzle) Unscramble() {
	n := len(p.stmts)
	pass := NewBuffer(p.hash)
	for i := n - 1; i >= 0; i-- {
		pass.Exec(p.stmts[i], true)
	}
	pass.Println()
}

type Buffer struct {
	head int
	data []rune
}

func (b *Buffer) Exec(stmt []string, rollback bool) {
	switch stmt[0] {
	case "swap":
		if stmt[1] != "letter" {
			b.Swap(stmt)
		} else {
			b.SwapLetter(stmt)
		}
	case "rotate":
		if stmt[1] != "based" {
			b.Rotate(stmt, rollback)
		} else {
			b.RotateLetter(stmt, rollback)
		}
	case "reverse":
		b.Reverse(stmt)
	case "move":
		b.Move(stmt, rollback)
	}
}

func (b Buffer) Swap(stmt []string) {
	n := len(b.data)
	i := (b.head + atoi(stmt[2])) % n
	j := (b.head + atoi(stmt[5])) % n
	b.data[i], b.data[j] = b.data[j], b.data[i]
}

func (b Buffer) SwapLetter(stmt []string) {
	x := rune(stmt[2][0])
	y := rune(stmt[5][0])
	var i, j int
	for k, v := range b.data {
		switch v {
		case x:
			i = k
		case y:
			j = k
		}
	}
	b.data[i], b.data[j] = b.data[j], b.data[i]
}

func (b *Buffer) Rotate(stmt []string, rollback bool) {
	n := len(b.data)
	step := atoi(stmt[2])
	if stmt[1] == "right" {
		step = -step
	}
	if rollback {
		step = -step
	}
	b.head += step
	for b.head < 0 {
		b.head += n
	}
	b.head %= n
}

func (b *Buffer) RotateLetter(stmt []string, rollback bool) {
	n := len(b.data)
	x := rune(stmt[6][0])
	var steps int
	for k, v := range b.data {
		if v == x {
			steps = k
			break
		}
	}
	steps -= b.head
	if steps < 0 {
		steps += n
	}
	if !rollback {
		if steps >= 4 {
			steps += 2
		} else {
			steps++
		}
		b.head -= steps
		for b.head < 0 {
			b.head += n
		}
	} else {
		b.head += REV[steps]
	}
	b.head %= n
}

func (b Buffer) Reverse(stmt []string) {
	n := len(b.data)
	x := b.head + atoi(stmt[2])
	y := b.head + atoi(stmt[4])
	for i, j := x, y; i < j; i, j = i+1, j-1 {
		i1, j1 := i%n, j%n
		b.data[i1], b.data[j1] = b.data[j1], b.data[i1]
	}
}

func (b Buffer) Move(stmt []string, rollback bool) {
	n := len(b.data)
	x := b.head + atoi(stmt[2])
	y := b.head + atoi(stmt[5])
	if rollback {
		x, y = y, x
	}
	r := b.data[x%n]
	d := 1
	if x > y {
		d = -1
	}
	for i := x; i != y; i += d {
		b.data[(i+n)%n] = b.data[(i+d+n)%n]
	}
	b.data[y%n] = r
}

func (b Buffer) Println() {
	var sb strings.Builder
	n := len(b.data)
	last := b.head + n
	for i := b.head; i < last; i++ {
		sb.WriteRune(b.data[i%n])
	}
	fmt.Println(sb.String())
}

func NewBuffer(s string) Buffer {
	return Buffer{0, []rune(s)}
}
