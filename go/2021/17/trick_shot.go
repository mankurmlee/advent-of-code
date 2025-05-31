package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"slices"
)

type Vertex [2]int
type Rect struct { lo, hi Vertex }

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func load(filename string) (t Rect) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Scan()
	check(s.Err())
	r := strings.NewReplacer("=", " ", "..", " ", ",", " ")
	data := strings.Fields(r.Replace(s.Text()))
	t.lo[0] = atoi(data[3])
	t.lo[1] = atoi(data[6])
	t.hi[0] = atoi(data[4])
	t.hi[1] = atoi(data[7])
	return t
}

func (t Rect) ValidX() (valid []int) {
	for start := 1; start <= t.hi[0]; start++ {
		x := 0
		for xs := start; xs >= 0; xs-- {
			x += xs
			if x >= t.lo[0] && x <= t.hi[0] {
				valid = append(valid, start)
				break
			}
		}
	}
	return valid
}

func (t Rect) ValidY() (valid []int) {
	for start := t.lo[1]; start <= -t.lo[1]; start++ {
		ys := start
		for y := 0; y > t.lo[1]; ys-- {
			y += ys
			if y >= t.lo[1] && y <= t.hi[1] {
				valid = append(valid, start)
				break
			}
		}
	}
	return valid
}

func (t Rect) Simulate() (best, num int) {
	xList := t.ValidX()
	yList := t.ValidY()
	slices.Reverse(yList)
	for _, ys := range yList {
		for _, xs := range xList {
			if t.HitBy(xs, ys) {
				if best == 0 {
					best = (ys * ys + ys) / 2
				}
				num++
			}
		}
	}
	return best, num
}

func (t Rect) HitBy(xs, ys int) bool {
	var x, y int
	for x < t.hi[0] && y > t.lo[1] {
		x += xs
		y += ys
		if x >= t.lo[0] && x <= t.hi[0] && y >= t.lo[1] && y <= t.hi[1] {
			return true
		}
		if xs > 0 {
			xs--
		}
		ys--
	}
	return false
}

func main() {
	flag.Parse()
	target := load(flag.Arg(0))
	fmt.Println(target)
	p1, p2 := target.Simulate()
	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)
}
