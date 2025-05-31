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

type Vec struct {
	x, y int
}

func (v *Vec) AddEquals(o Vec) {
	v.x += o.x
	v.y += o.y
}

type Robot struct {
	p, v Vec
}

func (r *Robot) Move() {
	r.p.AddEquals(r.v)
}

type Puzzle struct {
	room   Vec
	robots []Robot
}

func (p Puzzle) Animate() {
	w, h := p.room.x, p.room.y
	if w < 20 {
		return
	}
	rs := make([]Robot, len(p.robots))
	copy(rs, p.robots)
	var i int
	for ; !p.Draw(rs); i++ {
		for j := range rs {
			rs[j].Move()
			rs[j].p.x %= w
			rs[j].p.y %= h
			if rs[j].p.x < 0 {
				rs[j].p.x += w
			}
			if rs[j].p.y < 0 {
				rs[j].p.y += h
			}
		}
	}
	fmt.Println(i)
}

func (p Puzzle) Draw(rs []Robot) (found bool) {
	m := make(map[Vec]bool)
	for _, r := range rs {
		m[r.p] = true
	}
	var pic []string
	for y := range p.room.y {
		var sb strings.Builder
		for x := range p.room.x {
			if m[Vec{x, y}] {
				sb.WriteString("*")
			} else {
				sb.WriteString(" ")
			}
		}
		s := sb.String()
		if strings.Contains(s, "********") {
			found = true
		}
		pic = append(pic, s)
	}

	if found {
		for _, l := range pic {
			fmt.Println(l)
		}
	}

	return found
}

func (p Puzzle) SafetyFactor(n int) {
	q := make(map[Vec]int)
	for _, r := range p.Simulate(n) {
		w := (p.room.x + 1) >> 1
		h := (p.room.y + 1) >> 1
		if r.p.x == w-1 || r.p.y == h-1 {
			continue
		}
		q[Vec{r.p.x / w, r.p.y / h}]++
	}
	prod := 1
	for _, v := range q {
		prod *= v
	}
	fmt.Println(prod)
}

func (p Puzzle) Simulate(n int) []Robot {
	w, h := p.room.x, p.room.y
	rs := make([]Robot, len(p.robots))
	copy(rs, p.robots)
	for range n {
		for i := range rs {
			rs[i].Move()
		}
	}
	for i := range rs {
		rs[i].p.x %= w
		rs[i].p.y %= h
		if rs[i].p.x < 0 {
			rs[i].p.x += w
		}
		if rs[i].p.y < 0 {
			rs[i].p.y += h
		}
	}
	return rs
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
	r := regexp.MustCompile(`\-?\d+`)
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

func load(filename string) (p Puzzle) {
	lines := readFile(filename)
	d := findInts(lines[0])
	p.room = Vec{d[0], d[1]}
	for _, l := range lines[1:] {
		d = findInts(l)
		p.robots = append(p.robots, Robot{Vec{d[0], d[1]}, Vec{d[2], d[3]}})
	}
	return p
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	p.SafetyFactor(100)
	p.Animate()
}
