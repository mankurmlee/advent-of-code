package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Processor struct {
	regs map[string]int
}

func (p Processor) Run(prog []string) {
	n := len(prog)
	for ip := p.regs["ip"]; ip >= 0 && ip < n; ip = p.regs["ip"] {
		p.Interp(prog[ip])
	}
	fmt.Println(p.regs["a"])
}

func (p Processor) Interp(stmt string) {
	w := parseWords(stmt)
	switch w[0] {
	case "cpy":
		p.regs[w[2]] = p.Operand(w[1])
	case "inc":
		p.regs[w[1]]++
	case "dec":
		p.regs[w[1]]--
	case "jnz":
		if p.Operand(w[1]) != 0 {
			p.regs["ip"] += atoi(w[2])
			return
		}
	default:
		panic("Unexpected command")
	}
	p.regs["ip"]++
}

func (p Processor) Operand(s string) int {
	if strings.Contains("abcd", s) {
		return p.regs[s]
	}
	return atoi(s)
}

func NewProcessor() (p Processor) {
	p.regs = make(map[string]int)
	return p
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

func parseWords(s string) (out []string) {
	re := regexp.MustCompile(`(\w|[-])+`)
	return re.FindAllString(s, -1)
}

func main() {
	flag.Parse()
	p := readFile(flag.Arg(0))
	z := NewProcessor()
	z.Run(p)
	z = NewProcessor()
	z.regs["c"] = 1
	z.Run(p)
}
