package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Instruction struct {
	cmd string
	val int
}

type Program []Instruction

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

func load(filename string) (p Program) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		data := strings.Fields(s.Text())
		cmd := data[0]
		ins := Instruction{ cmd: cmd }
		if cmd == "addx" {
			ins.val = Atoi(data[1])
		}
		p = append(p, ins)
	}
	check(s.Err())
	return p
}

func (p Program) SampleSignalStrength() (total int) {
	var ip int
	plen := len(p)
	cycle := 0
	x := 1
	for i := 20; i < 260; i += 40 {
		for cycle < i {
			if ip >= plen {
				cycle++
				continue
			}
			ins := p[ip]
			ip++
			if ins.cmd == "addx" {
				cycle += 2
			} else {
				cycle++
			}
			if cycle >= i {
				total += i * x
			}
			x += ins.val
		}
	}
	return total
}

func (p Program) Render() {
	var ip, start, cost int
	x := 1
	for c := 1; c <= 240; c++ {
		pxl := (c - 1) % 40
		if pxl < x - 1 || pxl > x + 1 {
			fmt.Print(" ")
		} else {
			fmt.Print("*")
		}

		if p[ip].cmd == "addx" {
			cost = 2
		} else {
			cost = 1
		}

		if c == start + cost {
			x += p[ip].val
			start = c
			ip++
		}

		if c % 40 == 0 {
			fmt.Println()
		}
	}
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	fmt.Println("Part 1:", p.SampleSignalStrength())

	fmt.Println("Part 2:")
	p.Render()
}
