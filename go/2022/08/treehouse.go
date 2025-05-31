package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strconv"
)

type Vec struct {
	x, y int
}

type Grid struct {
	width, height int
	trees []int
}

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

func load(filename string) (g Grid) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)

	for s.Scan() {
		data := s.Text()

		if g.width == 0 {
			g.width = len(data)
		}
		g.height++

		for _, ch := range data {
			g.trees = append(g.trees, Atoi(string(ch)))
		}
	}

	check(s.Err())
	return g
}

func (g Grid) CountVisible() (count int) {
	w, h := g.width, g.height
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if  y == 0 || y == h - 1 ||
				x == 0 || x == w - 1 {
				count++
				continue
			}
			if g.TreeVisibleAt(Vec{x, y}) {
				count++
			}
		}
	}
	return count
}

func (v1 *Vec) Add(v2 Vec) {
	v1.x += v2.x
	v1.y += v2.y
}

func (g Grid) TreeVisibleAt(pos Vec) bool {
	w, h := g.width, g.height
	treeheight := g.trees[pos.y * w + pos.x]

	for _, v := range []Vec{
		Vec{-1, 0},
		Vec{1, 0},
		Vec{0, -1},
		Vec{0, 1},
	} {
		curr := pos
		for {
			curr.Add(v)
			if  curr.x < 0 || curr.x > w - 1 ||
				curr.y < 0 || curr.y > h - 1 {
					return true
			}
			if g.trees[curr.y * w + curr.x] >= treeheight {
				break
			}
		}
	}
	return false
}

func (g Grid) ScenicScore() (best int) {
	w, h := g.width, g.height
	for y := 1; y < h - 1; y++ {
		for x := 1; x < w - 1; x++ {
			score := g.TreeScoreAt(Vec{x, y})
			if score > best {
				best = score
			}
		}
	}
	return best
}

func (g Grid) TreeScoreAt(pos Vec) (score int) {
	w, h := g.width, g.height
	treeheight := g.trees[pos.y * w + pos.x]

	score = 1
	for _, v := range []Vec{
		Vec{-1, 0},
		Vec{1, 0},
		Vec{0, -1},
		Vec{0, 1},
	} {
		curr := pos
		viewdist := 0
		for {
			curr.Add(v)
			if  curr.x < 0 || curr.x > w - 1 ||
				curr.y < 0 || curr.y > h - 1 {
					break
			}
			viewdist++
			if g.trees[curr.y * w + curr.x] >= treeheight {
				break
			}
		}
		score *= viewdist
	}
	return score
}

func main() {
	flag.Parse()
	g := load(flag.Arg(0))

	fmt.Println("Part 1:", g.CountVisible())
	fmt.Println("Part 2:", g.ScenicScore())
}
