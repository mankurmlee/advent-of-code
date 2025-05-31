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

type Beam struct {
	pos Vec
	vel Vec
}

type Tile struct {
	val rune
	top, bottom, left, right bool
}

type Grid [][]Tile

func (g *Grid) GetEnergy() (count int) {
	for _, row := range *g {
		for _, t := range row {
			if t.left || t.right || t.top || t.bottom {
				count++
			}
		}
	}
	return count
}

func (g *Grid) Cast(start Beam) {
	beams := []Beam{start}
	for len(beams) > 0 {
		for i := 0; i < len(beams); i++ {
			ok, split := g.CastBeam(&beams[i])
			if split != nil {
				beams = append(beams, *split)
			}
			if !ok {
				beams = append(beams[:i], beams[i+1:]...)
			}
		}
	}
}

func (g *Grid) CastBeam(b *Beam) (ok bool, split *Beam) {
	b.pos = b.pos.Add(b.vel)
	if  b.pos.x < 0 ||
		b.pos.y < 0 ||
		b.pos.x >= len((*g)[0]) ||
		b.pos.y >= len(*g) {
		return false, nil
	}
	tile := &(*g)[b.pos.y][b.pos.x]
	if  b.vel.x < 0 && tile.right ||
		b.vel.x > 0 && tile.left ||
		b.vel.y < 0 && tile.bottom ||
		b.vel.y > 0 && tile.top {
		return false, nil
	}

	// Update the tile status
	if b.vel.x < 0 {
		tile.right = true
	} else if b.vel.x > 0 {
		tile.left = true
	} else if b.vel.y < 0 {
		tile.bottom = true
	} else if b.vel.y > 0 {
		tile.top = true
	}

	// Update the beam status
	switch tile.val {
	case '/':
		if b.vel.x < 0 {
			b.vel = Vec{0, 1}
		} else if b.vel.x > 0 {
			b.vel = Vec{0, -1}
		} else if b.vel.y < 0 {
			b.vel = Vec{1, 0}
		} else if b.vel.y > 0 {
			b.vel = Vec{-1, 0}
		}
	case '\\':
		if b.vel.x < 0 {
			b.vel = Vec{0, -1}
		} else if b.vel.x > 0 {
			b.vel = Vec{0, 1}
		} else if b.vel.y < 0 {
			b.vel = Vec{-1, 0}
		} else if b.vel.y > 0 {
			b.vel = Vec{1, 0}
		}
	case '|':
		if b.vel.y == 0 {
			b.vel = Vec{0, -1}
			split = &Beam{ b.pos, Vec{0, 1} }
		}
	case '-':
		if b.vel.x == 0 {
			b.vel = Vec{-1, 0}
			split = &Beam{ b.pos, Vec{1, 0} }
		}
	}

	return true, split
}

func (g *Grid) Print() {
	var s strings.Builder
	for _, row := range *g {
		s.Reset()
		for _, t := range row {
			r := string(t.val)
			if r == "." {
				var energised int
				if t.right {
					energised++
					r = "<"
				}
				if t.left {
					energised++
					r = ">"
				}
				if t.top {
					energised++
					r = "v"
				}
				if t.bottom {
					energised++
					r = "^"
				}
				if energised > 1 {
					r = strconv.Itoa(energised)
				}
			}
			s.WriteString(r)
		}
		fmt.Println(s.String())
	}
}

func (v Vec) Add(v2 Vec) Vec {
	return Vec {
		v.x + v2.x,
		v.y + v2.y,
	}
}

func load(filename string) *Grid {
	var grid Grid

	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		var row []Tile
		for _, ch := range s.Text() {
			row = append(row, Tile{val: ch})
		}
		grid = append(grid, row)
	}
	check(s.Err())
	return &grid
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	g := load(flag.Arg(0))
	g.Cast(Beam{
		Vec{-1, 0},
		Vec{ 1, 0},
	})
	g.Print()
	fmt.Println("Part 1 answer is", g.GetEnergy())

	var startingBeams []Beam
	width, height := len((*g)[0]), len(*g)
	for x := 0; x < width; x++ {
		startingBeams = append(startingBeams, Beam{
			Vec{x, -1},
			Vec{0,  1},
		})
		startingBeams = append(startingBeams, Beam{
			Vec{x, height},
			Vec{0, -1},
		})
	}
	for y := 0; y < height; y++ {
		startingBeams = append(startingBeams, Beam{
			Vec{-1, y},
			Vec{1, 0},
		})
		startingBeams = append(startingBeams, Beam{
			Vec{width, y},
			Vec{-1, 0},
		})
	}
	var max int
	for _, b := range startingBeams {
		g := load(flag.Arg(0))
		g.Cast(b)
		c := g.GetEnergy()
		if max < c {
			max = c
		}
	}
	fmt.Println("Part 2 answer is", max)
}
