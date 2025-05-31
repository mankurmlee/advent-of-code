package main

import (
	"fmt"
	"os"
	"bufio"
	"flag"
	"slices"
)

type Pair [2]int

type Galaxy struct {
	x, y int
}

type Voids struct {
	scale int
	cols, rows []int
}

func (v Voids) getDistance(gs []Galaxy, p Pair) int {
	i := gs[p[0]]
	j := gs[p[1]]

	x1, x2 := Order(i.x, j.x)
	y1, y2 := Order(i.y, j.y)

	dx := x2 - x1
	dy := y2 - y1

	for _, c := range v.cols {
		if x1 < c && c < x2 {
			dx += v.scale - 1
		}
	}

	for _, r := range v.rows {
		if y1 < r && r < y2 {
			dy += v.scale - 1
		}
	}

	dist := dx + dy
	return dist
}

func (v Voids) sumDistances(gs []Galaxy, pairs []Pair) int {
	var sum int

	for _, p := range pairs {
		sum += v.getDistance(gs, p)
	}

	return sum
}

type Universe [][]rune

func (u Universe) GetVoids(scale int) Voids {
	return Voids {
		scale,
		u.GetHorizontalVoids(),
		u.GetVerticalVoids(),
	}
}

func (u Universe) GetHorizontalVoids() []int {
	var cols []int
	var isEmpty bool

	for i := 0; i < len(u[0]); i++ {
		isEmpty = true
		for _, r := range u {
			if r[i] == '#' {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			cols = append(cols, i)
		}
	}

	return cols
}

func (u Universe) GetVerticalVoids() []int {
	var rows []int

	for i, r := range u {
		if !slices.Contains(r, '#') {
			rows = append(rows, i)
		}
	}

	return rows
}

func (u Universe) GetGalaxies() []Galaxy {
	var gs []Galaxy

	for y, r := range u {
		for x, v := range r {
			if v == '#' {
				gs = append(gs, Galaxy{x, y})
			}
		}
	}

	return gs
}

func genPairs(n int) []Pair {
	var pairs []Pair

	for i := 0; i < n-1; i++ {
		for j := i+1; j < n; j++ {
			pairs = append(pairs, Pair{i, j})
		}
	}

	return pairs
}

func loadUniverse(filename string) Universe {
	var u Universe
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		u = append(u, []rune(s.Text()))
	}

	check(s.Err())
	return u
}

func Order(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var scale int
	flag.IntVar(&scale, "scale", 2, "Cosmic Expansion Scale")
	flag.Parse()
	filename := flag.Args()[0]

	uni := loadUniverse(filename)

	voids := uni.GetVoids(scale)
	gs := uni.GetGalaxies()
	pairs := genPairs(len(gs))

	sum := voids.sumDistances(gs, pairs)

	fmt.Println(len(gs), "galaxies found")
	fmt.Println("The sum of the shortest path between all", len(pairs),
		"pairs of galaxies is", sum)
}
