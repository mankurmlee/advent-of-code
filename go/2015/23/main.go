package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func fileinput() (out []string) {
	flag.Parse()
	f, err := os.Open(flag.Arg(0))
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		out = append(out, s.Text())
	}
	check(s.Err())
	return out
}

type Processor struct {
	Prog [][]string
	IP   int
	Regs map[string]int
}

func Load(a int) Processor {
	var prog [][]string
	re := strings.NewReplacer(",", " ")
	for _, l := range fileinput() {
		prog = append(prog, strings.Fields(re.Replace(l)))
	}
	regs := map[string]int{"a": a, "b": 0}
	return Processor{prog, 0, regs}
}

func (p *Processor) Step() {
	arg := p.Prog[p.IP]
	r := p.Regs
	switch arg[0] {
	case "hlf":
		r[arg[1]] >>= 1
	case "tpl":
		r[arg[1]] *= 3
	case "inc":
		r[arg[1]]++
	case "jmp":
		p.IP += atoi(arg[1]) - 1
	case "jie":
		if r[arg[1]]%2 == 0 {
			p.IP += atoi(arg[2]) - 1
		}
	case "jio":
		if r[arg[1]] == 1 {
			p.IP += atoi(arg[2]) - 1
		}
	}
	p.IP++
}

func (p Processor) Halted() bool {
	return p.IP >= len(p.Prog)
}

func (p *Processor) Run() int {
	for !p.Halted() {
		p.Step()
	}
	return p.Regs["b"]
}

func main() {
	p := Load(0)
	fmt.Println("Part 1:", p.Run())
	p = Load(1)
	fmt.Println("Part 2:", p.Run())
}
