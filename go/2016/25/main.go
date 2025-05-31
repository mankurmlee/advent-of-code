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
	prog []string
}

func (p Processor) Run() {
	n := len(p.prog)
	for ip := p.regs["ip"]; ip >= 0 && ip < n; ip = p.regs["ip"] {
		p.Interp()
	}
	fmt.Println(p.regs["a"])
}

func (p Processor) Interp() {
	stmt := p.prog[p.regs["ip"]]
	w := parseWords(stmt)
	switch w[0] {
	case "cpy":
		if strings.Contains("abcd", w[2]) {
			p.regs[w[2]] = p.Operand(w[1])
		}
	case "inc":
		p.regs[w[1]]++
	case "dec":
		p.regs[w[1]]--
	case "jnz":
		if p.Operand(w[1]) != 0 {
			p.regs["ip"] += p.Operand(w[2])
			// fmt.Println(stmt, p.regs)
			return
		}
	case "out":
		fmt.Print(p.Operand(w[1]))
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

func NewProcessor(prog []string) (p Processor) {
	p.regs = make(map[string]int)
	p.prog = make([]string, len(prog))
	copy(p.prog, prog)
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
	z := NewProcessor(p)
	z.regs["a"] = 196
	z.Run()
}
