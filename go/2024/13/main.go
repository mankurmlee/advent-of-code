package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Vec struct {
	x, y int
}

func (v Vec) Add(o Vec) Vec {
	return Vec{v.x + o.x, v.y + o.y}
}

func (v Vec) Mul(n int) Vec {
	return Vec{v.x * n, v.y * n}
}

type Machine struct {
	a, b, z Vec
}

func (m Machine) Solve() (int, int, bool) {
	i := m.z.y*m.b.x - m.z.x*m.b.y
	j := m.b.x*m.a.y - m.a.x*m.b.y
	if i%j != 0 {
		return 0, 0, false
	}
	a := i / j
	i = m.z.x - a*m.a.x
	j = m.b.x
	if i%j != 0 {
		return 0, 0, false
	}
	return a, i / j, true
}

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

func findInts(s string) (nums []int) {
	r := regexp.MustCompile(`\d+`)
	for _, v := range r.FindAllString(s, -1) {
		nums = append(nums, atoi(v))
	}
	return nums
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

func load(filename string) (p []Machine) {
	var y int
	var m Machine
	for _, l := range readFile(filename) {
		y++
		if l == "" {
			p = append(p, m)
			m = Machine{}
			y = 0
			continue
		}
		d := findInts(l)
		if y == 1 {
			m.a = Vec{d[0], d[1]}
		} else if y == 2 {
			m.b = Vec{d[0], d[1]}
		} else if y == 3 {
			m.z = Vec{d[0], d[1]}
		}
	}
	return append(p, m)
}

func main() {
	flag.Parse()
	puzzle := load(flag.Arg(0))
	GetCost(puzzle)
	GetCost(offset(puzzle))
}

func offset(ms []Machine) (out []Machine) {
	for _, m := range ms {
		z := m.z.Add(Vec{10000000000000, 10000000000000})
		out = append(out, Machine{m.a, m.b, z})
	}
	return out
}

func GetCost(ms []Machine) {
	var tot int
	for _, m := range ms {
		a, b, ok := m.Solve()
		if ok {
			tot += 3*a + b
		}
	}
	fmt.Println(tot)
}
