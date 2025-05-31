package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Operation int
const (
	TURN_OFF Operation = iota
	TURN_ON
	TOGGLE
)

type Vertex struct {x, y int}

type Command struct {
	op     Operation
	nw, se Vertex
}

type Puzzle []Command

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load(filename string) (p Puzzle) {
	var cmd Command
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		data := strings.Fields(s.Text())
		switch data[1] {
		case "on":
			cmd = Command{ TURN_ON,  parseVertex(data[2]), parseVertex(data[4]) }
		case "off":
			cmd = Command{ TURN_OFF, parseVertex(data[2]), parseVertex(data[4]) }
		default:
			cmd = Command{ TOGGLE,   parseVertex(data[1]), parseVertex(data[3]) }
		}
		p = append(p, cmd)
	}
	check(s.Err())
	return p
}

func parseVertex(s string) Vertex {
	data := strings.Split(s, ",")
	return Vertex{ Atoi(data[0]), Atoi(data[1]) }
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func (p Puzzle) CountLights() (num int) {
	lights := make([]bool, 1000000)
	for _, cmd := range p {
		switch cmd.op {
		case TURN_OFF:
			for y := cmd.nw.y; y <= cmd.se.y; y++ {
				for x := cmd.nw.x; x <= cmd.se.x; x++ {
					lights[y * 1000 + x] = false
				}
			}
		case TURN_ON:
			for y := cmd.nw.y; y <= cmd.se.y; y++ {
				for x := cmd.nw.x; x <= cmd.se.x; x++ {
					lights[y * 1000 + x] = true
				}
			}
		case TOGGLE:
			for y := cmd.nw.y; y <= cmd.se.y; y++ {
				for x := cmd.nw.x; x <= cmd.se.x; x++ {
					lights[y * 1000 + x] = !lights[y * 1000 + x]
				}
			}
		}
	}
	for _, l := range lights {
		if l {
			num++
		}
	}
	return num
}

func (p Puzzle) TotalBrightness() (num int) {
	lights := make([]int, 1000000)
	for _, cmd := range p {
		switch cmd.op {
		case TURN_OFF:
			for y := cmd.nw.y; y <= cmd.se.y; y++ {
				for x := cmd.nw.x; x <= cmd.se.x; x++ {
					if lights[y * 1000 + x] > 0 {
						lights[y * 1000 + x]--
					}
				}
			}
		case TURN_ON:
			for y := cmd.nw.y; y <= cmd.se.y; y++ {
				for x := cmd.nw.x; x <= cmd.se.x; x++ {
					lights[y * 1000 + x]++
				}
			}
		case TOGGLE:
			for y := cmd.nw.y; y <= cmd.se.y; y++ {
				for x := cmd.nw.x; x <= cmd.se.x; x++ {
					lights[y * 1000 + x] += 2
				}
			}
		}
	}
	for _, l := range lights {
		num += l
	}
	return num
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	//~ fmt.Println("Part 1:", p.CountLights())
	fmt.Println("Part 2:", p.TotalBrightness())
}
