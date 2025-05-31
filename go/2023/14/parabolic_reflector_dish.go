package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
)

type Platform struct {
	grid     [][]rune
	width    int
	height   int
}

func (p *Platform) Print() {
	for _, r := range p.grid {
		fmt.Println(string(r))
	}
}

func (p *Platform) Tilt() {
	var w int
	for x := 0; x < p.width; x++ {
		w = -1
		for y := 0; y < p.height; y++ {
			switch p.grid[y][x] {
			case '.':
				if w < 0 {
					w = y
				}
			case '#':
				if w >= 0 {
					w = -1
				}
			case 'O':
				if w >= 0 {
					p.grid[w][x], p.grid[y][x] = p.grid[y][x], p.grid[w][x]
					for w < y {
						if p.grid[w][x] == '.' {
							break
						}
						w++
					}
				}
			}
		}
	}
}

func (p *Platform) Cycle() {
	var w int

	width, height := p.width, p.height
	q := p.grid

	// Tilt North
	for x := 0; x < width; x++ {
		w = -1
		for y := 0; y < height; y++ {
			switch q[y][x] {
			case '.':
				if w < 0 {
					w = y
				}
			case '#':
				if w >= 0 {
					w = -1
				}
			case 'O':
				if w >= 0 {
					q[w][x], q[y][x] = q[y][x], q[w][x]
					for w < y {
						if q[w][x] == '.' {
							break
						}
						w++
					}
				}
			}
		}
	}

	// Tilt West
	for y := 0; y < height; y++ {
		w = -1
		for x := 0; x < width; x++ {
			switch q[y][x] {
			case '.':
				if w < 0 {
					w = x
				}
			case '#':
				if w >= 0 {
					w = -1
				}
			case 'O':
				if w >= 0 {
					q[y][w], q[y][x] = q[y][x], q[y][w]
					for w < x {
						if q[y][w] == '.' {
							break
						}
						w++
					}
				}
			}
		}
	}

	// Tilt South
	for x := 0; x < width; x++ {
		w = -1
		for y := height-1; y >= 0; y-- {
			switch q[y][x] {
			case '.':
				if w < 0 {
					w = y
				}
			case '#':
				if w >= 0 {
					w = -1
				}
			case 'O':
				if w >= 0 {
					q[w][x], q[y][x] = q[y][x], q[w][x]
					for w > y {
						if q[w][x] == '.' {
							break
						}
						w--
					}
				}
			}
		}
	}

	// Tilt East
	for y := 0; y < height; y++ {
		w = -1
		for x := width-1; x >= 0; x-- {
			switch q[y][x] {
			case '.':
				if w < 0 {
					w = x
				}
			case '#':
				if w >= 0 {
					w = -1
				}
			case 'O':
				if w >= 0 {
					q[y][w], q[y][x] = q[y][x], q[y][w]
					for w > x {
						if q[y][w] == '.' {
							break
						}
						w--
					}
				}
			}
		}
	}
}

func (p *Platform) GetLoad() int {
	var load int
	for y := 0; y < p.height; y++ {
		for x := 0; x < p.width; x++ {
			if p.grid[y][x] == 'O' {
				load += p.height - y
			}
		}
	}
	return load
}

func loadPlat(filename string) *Platform {
	var p Platform
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		p.grid = append(p.grid, []rune(s.Text()))
	}
	check(s.Err())

	p.width  = len(p.grid[0])
	p.height = len(p.grid)

	return &p
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func part2(history []int) int {
	var period int
	var matches bool
	num := len(history)
	needle := history[num - 1]
	for i := 1; i <= 100; i++ {
		matches = true
		for j := 1; j <= 3; j++ {
			if needle != history[num - i*j - 1] {
				matches = false
				break
			}
		}
		if matches {
			period = i
			break
		}
	}
	if period == 0 {
		panic("Answer out of range")
	}
	offset := 1000000000 % period
	for offset <= num {
		offset += period
	}
	return offset - num
}

func main() {
	var history []int

	flag.Parse()
	filename := flag.Arg(0)

	// Part 1
	p := loadPlat(filename)
	p.Tilt()
	fmt.Println("Part 1 answer:", p.GetLoad())

	// Part 2
	p = loadPlat(filename)
	for i := 0; i < 300; i++ {
		p.Cycle()
		history = append(history, p.GetLoad())
	}
	extra := part2(history)
	for i := 0; i < extra; i++ {
		p.Cycle()
	}
	p.Print()
	fmt.Println("Part 2 answer:", p.GetLoad())
}
