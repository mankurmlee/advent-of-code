package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Pattern []string

func (p Pattern) GetReflect() int {
	nrows := len(p)
	for r := 0; r < nrows-1; r++ {
		if p.ReflectsAt(r) {
			return r + 1
		}
	}

	return 0
}

func (p Pattern) ReflectsAt(i int) bool {
	height := len(p)
	j := i + 1
	for i >= 0 && j < height {
		if p[i] != p[j] {
			return false
		}
		i--
		j++
	}
	return true
}

func (p Pattern) Inverse() Pattern {
	var inverse Pattern

	ncols := len(p[0])
	sbs := make([]strings.Builder, ncols)

	for _, row := range p {
		for i, ch := range row {
			sbs[i].WriteRune(ch)
		}
	}

	for _, sb := range sbs {
		inverse = append(inverse, sb.String())
	}

	return inverse
}

func (p Pattern) Print() {
	for _, r := range p {
		fmt.Println(r)
	}
}

func (p Pattern) GetAlmostReflect(ignore int) int {
	nrows := len(p)
	for r := 0; r < nrows-1; r++ {
		if (r + 1) * 100 == ignore {
			continue
		}
		if p.AlmostReflectsAt(r) {
			return (r + 1) * 100
		}
	}

	inv := p.Inverse()
	nrows = len(inv)
	for r := 0; r < nrows-1; r++ {
		if r + 1 == ignore {
			continue
		}
		if inv.AlmostReflectsAt(r) {
			return r + 1
		}
	}

	return 0
}

func (p Pattern) AlmostReflectsAt(i int) bool {
	var hasDiff bool
	var a, b []rune

	height := len(p)
	j := i + 1
	for i >= 0 && j < height {
		if p[i] != p[j] {
			if hasDiff {
				return false
			}
			a = []rune(p[i])
			b = []rune(p[j])
			hasDiff = true
		}
		i--
		j++
	}

	hasDiff = false
	width := len(a)
	for k := 0; k < width; k++ {
		if a[k] != b[k] {
			if hasDiff {
				return false
			}
			hasDiff = true
		}
	}

	return true
}

func load(filename string) []Pattern {
	var patterns []Pattern

	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	var p Pattern
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			patterns = append(patterns, p)
			p = Pattern{}
			continue
		}
		p = append(p, line)
	}
	check(s.Err())
	patterns = append(patterns, p)

	return patterns
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var part1, part2 int
	var ref []int

	flag.Parse()
	filename := flag.Args()[0]

	patterns := load(filename)

	for _, p := range patterns {
		sum := p.GetReflect() * 100 + p.Inverse().GetReflect()
		ref = append(ref, sum)
		part1 += sum
	}

	for i, p := range patterns {
		ref1 := p.GetAlmostReflect(ref[i])
		fmt.Println("Pattern", i+1, "reflects at", ref[i], "and almost reflects at", ref1)
		p.Print()
		part2 += ref1
	}

	fmt.Println("Part 1 result is", part1)
	fmt.Println("Part 2 result is", part2)
}
