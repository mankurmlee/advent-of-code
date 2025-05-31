package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

const (
	X = 0
	Y = 1
)

type Vertex [2]int

type Paper map[Vertex]bool

type Matrix [6]int

type Fold struct {
	axis int
	n    int
}

type Puzzle struct {
	dots  Paper
	folds []Fold
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

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	r := strings.NewReplacer(",", " ", "=", " ")
	amap := map[string]int{ "x": X, "y": Y }
	p.dots = Paper{}
	for s.Scan() {
		d := strings.Fields(r.Replace(s.Text()))
		switch len(d) {
		case 2:
			p.dots[Vertex{ atoi(d[0]), atoi(d[1]) }] = true
		case 4:
			p.folds = append(p.folds, Fold{ amap[d[2]], atoi(d[3]) })
		}
	}
	check(s.Err())
	return p
}

func (f Fold) GetTransform() (m Matrix) {
	if f.axis == X {
		m[0] = -1
		m[4] = 1
		m[2] = 2 * f.n
	} else {
		m[0] = 1
		m[4] = -1
		m[5] = 2 * f.n
	}
	return m
}

func (p Paper) DoFold(f Fold) Paper {
	m := f.GetTransform()
	q := Paper{}
	for d := range p {
		if d[f.axis] > f.n {
			q[m.Transform(d)] = true
		} else {
			q[d] = true
		}
	}
	return q
}

func (m Matrix) Transform(in Vertex) (out Vertex) {
	out[0] = m[0] * in[0] + m[1] * in[1] + m[2]
	out[1] = m[3] * in[0] + m[4] * in[1] + m[5]
	return out
}

func (p Paper) String() string {
	var w, h int
	for v := range p {
		w, h = max(w, v[0]), max(h, v[1])
	}
	w++
	h++
	rows := make([]string, h)
	for y := 0; y < h; y++ {
		row := make([]string, w)
		for x := 0; x < w; x++ {
			if p[Vertex{x, y}] {
				row[x] = "*"
			} else {
				row[x] = " "
			}
		}
		rows[y] = strings.Join(row, "")
	}
	return strings.Join(rows, "\n")
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	f1 := p.dots.DoFold(p.folds[0])
	fmt.Println("Part 1:", len(f1))

	f2 := p.dots
	for _, f := range p.folds {
		f2 = f2.DoFold(f)
	}
	fmt.Println("Part 2:")
	fmt.Println(f2)
}
