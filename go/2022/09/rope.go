package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Vec struct {
	x, y int
}

type Rope struct {
	knots  []Vec
	trail  map[Vec]struct{}
}

type Motion struct {
	dir   Vec
	steps int
}

type Instructions []Motion

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func Abs(n int) int {
	if n >= 0 {
		return n
	}
	return -n
}

func Sgn(n int) int {
	switch {
	case n > 0:
		return 1
	case n < 0:
		return -1
	}
	return 0
}

func load(filename string) (p Instructions) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		data := strings.Fields(s.Text())
		var dir Vec
		switch data[0] {
		case "U":
			dir = Vec{0, -1}
		case "D":
			dir = Vec{0, 1}
		case "L":
			dir = Vec{-1, 0}
		case "R":
			dir = Vec{1, 0}
		}
		p = append(p, Motion{ dir, Atoi(data[1]) })
	}
	check(s.Err())
	return p
}

func (r *Rope) Create(numKnots int) *Rope {
	r = &Rope{
		make([]Vec, numKnots),
		make(map[Vec]struct{}),
	}
	r.trail[r.knots[numKnots-1]] = struct{}{}
	return r
}

func (r *Rope) Follow(d Instructions) {
	for _, m := range d {
		r.Apply(m)
	}
}

func (r *Rope) Apply(m Motion) {
	n := len(r.knots)
	for i := 0; i < m.steps; i++ {
		r.knots[0].x += m.dir.x
		r.knots[0].y += m.dir.y
		for j := 1; j < n; j++ {
			dx := r.knots[j-1].x - r.knots[j].x
			dy := r.knots[j-1].y - r.knots[j].y
			if Abs(dx) <= 1 && Abs(dy) <= 1 {
				break
			}
			r.knots[j].x += Sgn(dx)
			r.knots[j].y += Sgn(dy)
		}
		r.trail[r.knots[n-1]] = struct{}{}
	}
	//~ fmt.Println(r.knots)
}

func main() {
	flag.Parse()
	d := load(flag.Arg(0))

	var r1 *Rope
	r1 = r1.Create(2)
	r1.Follow(d)
	fmt.Println("Part 1:", len(r1.trail))

	var r2 *Rope
	r2 = r2.Create(10)
	r2.Follow(d)
	fmt.Println("Part 2:", len(r2.trail))
}
