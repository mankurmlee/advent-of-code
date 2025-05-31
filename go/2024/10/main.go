package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Vec struct {
	x, y int
}

type Map struct {
	w, h int
	data []int
}

var directions = []Vec{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

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

func load(filename string) (m Map) {
	d := readFile(filename)
	m.w = len(d[0])
	m.h = len(d)
	for _, l := range d {
		for _, r := range l {
			m.data = append(m.data, atoi(string(r)))
		}
	}
	return m
}

func main() {
	flag.Parse()
	m := load(flag.Arg(0))
	n := 0
	r := 0
	for i, v := range m.data {
		if v == 0 {
			n1, r1 := m.CountTrails(Vec{i % m.w, i / m.w})
			n += n1
			r += r1
		}
	}
	fmt.Println(n)
	fmt.Println(r)
}

func (m Map) Get(v Vec) int {
	if v.x < 0 || v.y < 0 || v.x >= m.w || v.y >= m.h {
		return -1
	}
	return m.data[v.y*m.w+v.x]
}

func (m Map) CountTrails(head Vec) (int, int) {
	s := make(map[Vec]struct{})
	r := 0
	q := []Vec{head}
	for len(q) > 0 {
		end := len(q) - 1
		v := q[end]
		q = q[:end]
		z := m.Get(v)
		if z == 9 {
			s[v] = struct{}{}
			r++
			continue
		}
		for _, d := range directions {
			v1 := Vec{v.x + d.x, v.y + d.y}
			z1 := m.Get(v1)
			if z1 == z+1 {
				q = append(q, v1)
			}
		}
	}
	return len(s), r
}
