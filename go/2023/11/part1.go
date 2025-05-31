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

type Universe [][]rune

func (u Universe) Print() {
	for _, r := range u {
		fmt.Println(string(r))
	}
}

func (u Universe) Expand() Universe {
	wider := u.ExpandHorizontally()
	return wider.ExpandVertically()
}

func (u Universe) ExpandHorizontally() Universe {
	var expanded Universe
	var row []rune
	var emptyCols []int
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
			emptyCols = append(emptyCols, i)
		}
	}

	for _, r := range u {
		row = []rune{}
		for i, v := range r {
			if slices.Contains(emptyCols, i) {
				row = append(row, v)
			}
			row = append(row, v)
		}
		expanded = append(expanded, row)
	}

	return expanded
}

func (u Universe) ExpandVertically() Universe {
	var new Universe

	for _, r := range u {
		if !slices.Contains(r, '#') {
			new = append(new, r)
		}
		new = append(new, r)
	}

	return new
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

func sumDistances(gs []Galaxy, pairs []Pair) int {
	var sum int

	for _, p := range pairs {
		sum += getDistance(gs, p)
	}

	return sum
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getDistance(gs []Galaxy, p Pair) int {
	dist := Abs(gs[p[1]].x - gs[p[0]].x) + Abs(gs[p[1]].y - gs[p[0]].y)
	//fmt.Printf("Distance between galaxy %d and %d: %d\n", p[0]+1, p[1]+1, dist)
	return dist
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	filename := flag.Args()[0]

	fmt.Println("Part 1")

	u := loadUniverse(filename)
	fmt.Println("")
	fmt.Println("The Universe:")
	u.Print()

	exp := u.Expand()
	fmt.Println("")
	fmt.Println("Expanded Universe:")
	exp.Print()

	gs := exp.GetGalaxies()
	fmt.Println("")
	fmt.Println(len(gs), "galaxies found")

	pairs := genPairs(len(gs))
	sum := sumDistances(gs, pairs)
	fmt.Println("The sum of the shortest path between all", len(pairs),
		"pairs of galaxies is", sum)
}
