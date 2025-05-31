package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Vec struct{ x, y int }

func (v Vec) Add(o Vec) Vec {
	return Vec{v.x + o.x, v.y + o.y}
}

func (v Vec) Equals(o Vec) bool {
	return v.x == o.x && v.y == o.y
}

type State struct {
	pos  Vec
	cost int
}

type Puzzle struct {
	dest  Vec
	fav   int
	walls map[Vec]bool
}

var DIRS = []Vec{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

func (p Puzzle) Search() {
	var p2 int
	var s0 State
	s0.pos = Vec{1, 1}
	q := []State{s0}
	seen := make(map[Vec]struct{})
	seen[s0.pos] = struct{}{}
	for len(q) > 0 {
		s := q[0]
		q = q[1:]
		if s.pos.Equals(p.dest) {
			fmt.Println(s.cost)
			fmt.Println(p2)
			return
		}
		for _, d := range DIRS {
			s1 := State{
				s.pos.Add(d),
				s.cost + 1,
			}
			if s1.cost == 51 && p2 == 0 {
				p2 = len(seen)
			}
			if s1.pos.x < 0 || s1.pos.y < 0 {
				continue
			}
			if p.IsWall(s1.pos) {
				continue
			}
			if _, ok := seen[s1.pos]; ok {
				continue
			}
			seen[s1.pos] = struct{}{}
			q = append(q, s1)
		}
	}
	fmt.Println("Couldn't find solution!")
}

func (p Puzzle) IsWall(pos Vec) bool {
	w, ok := p.walls[pos]
	if ok {
		return w
	}
	n := (pos.x+pos.y)*(pos.x+pos.y+1) + pos.x + pos.x + p.fav
	w = countOnes(n)%2 == 1
	p.walls[pos] = w
	return w
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

func countOnes(n int) (count int) {
	for n != 0 {
		n &= n - 1
		count++
	}
	return count
}

func parseInts(s string) (nums []int) {
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

func load(filename string) Puzzle {
	d := readFile(filename)
	r := parseInts(d[0])
	return Puzzle{
		Vec{r[0], r[1]},
		atoi(d[1]),
		make(map[Vec]bool),
	}
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	p.Search()
}
