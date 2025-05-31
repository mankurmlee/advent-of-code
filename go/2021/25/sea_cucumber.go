package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Vec [2]int

type State struct {
	cuc   [2][]Vec
	grid  []bool
	steps int
	size  Vec
}

func (s *State) Update() (changed bool) {
	var x, y, i int
	var v Vec
	willMove := make([]int, 0, max(len(s.cuc[0]), len(s.cuc[1])))
	w, h := s.size[0], s.size[1]

	for a := range 2 {
		willMove = willMove[:0]
		for i, v = range s.cuc[a] {
			x, y = v[0], v[1]
			switch a {
			case 0:
				x += 1
				if x >= w {
					x = 0
				}
			case 1:
				y += 1
				if y >= h {
					y = 0
				}
			}
			if !s.grid[x+y*w] {
				willMove = append(willMove, i)
			}
		}

		for _, i = range willMove {
			v = s.cuc[a][i]
			s.grid[v[0]+v[1]*w] = false
			switch a {
			case 0:
				v[0] += 1
				if v[0] >= w {
					v[0] = 0
				}
			case 1:
				v[1] += 1
				if v[1] >= h {
					v[1] = 0
				}
			}
			s.grid[v[0]+v[1]*w] = true
			s.cuc[a][i] = v
		}
		changed = changed || len(willMove) > 0
	}
	s.steps++
	return changed
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load(filename string) (state State) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	var w, h, y int
	for s.Scan() {
		data := []byte(s.Text())
		for x, b := range data {
			switch b {
			case '>':
				state.cuc[0] = append(state.cuc[0], Vec{x, y})
			case 'v':
				state.cuc[1] = append(state.cuc[1], Vec{x, y})
			}
		}
		if w == 0 {
			w = len(data) // store puzzle width
		}
		y++ // increase puzzle height
	}
	check(s.Err())
	h = y
	state.grid = make([]bool, w*h)
	for i := range 2 {
		for _, v := range state.cuc[i] {
			state.grid[v[0]+v[1]*w] = true
		}
	}
	state.size[0], state.size[1] = w, h
	return state
}

func main() {
	flag.Parse()
	s := load(flag.Arg(0))

	fmt.Println(s)
	changed := true
	for changed {
		changed = s.Update()
	}
	fmt.Println(s)

	fmt.Println("Part 1:", s.steps)
}

func (s State) String() string {
	w, h := s.size[0], s.size[1]
	l := w * h
	grid := make([]byte, l)
	for i := range grid {
		grid[i] = '.'
	}
	for _, v := range s.cuc[0] {
		grid[v[0]+v[1]*w] = '>'
	}
	for _, v := range s.cuc[1] {
		grid[v[0]+v[1]*w] = 'v'
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Step %d:\n", s.steps))
	for i := 0; i < l; i += w {
		sb.Write(grid[i : i+w])
		sb.WriteRune('\n')
	}
	return sb.String()
}
