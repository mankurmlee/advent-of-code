package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Circuit struct {
	expr map[string][]string
	eval map[string]uint16
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load(filename string) (p Circuit) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	p.expr = map[string][]string{}
	p.eval = map[string]uint16{}
	for s.Scan() {
		data := strings.Fields(s.Text())
		n := len(data)
		p.expr[data[n-1]] = data[:n-2]
	}
	check(s.Err())
	return p
}

func (p Circuit) Eval(s string) uint16 {
	var num uint16

	// Check cache
	if num, ok := p.eval[s]; ok {
		return num
	}

	expr, ok := p.expr[s]
	if !ok {
		i, err := strconv.ParseUint(s, 10, 16)
		check(err)
		return uint16(i)
	}

	switch len(expr) {
	case 1:
		num = p.Eval(expr[0])
	case 2:
		num = ^p.Eval(expr[1])
	case 3:
		a, b := p.Eval(expr[0]), p.Eval(expr[2])
		switch expr[1] {
		case "AND":
			num = a & b
		case "OR":
			num = a | b
		case "LSHIFT":
			num = a << b
		case "RSHIFT":
			num = a >> b
		}
	}

	p.eval[s] = num
	return num
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	for k := range p.expr {
		v := p.Eval(k)
		fmt.Println(k, "=", v)
	}

	if _, ok := p.expr["a"]; !ok {
		return
	}

	a := p.Eval("a")
	clear(p.eval)
	p.eval["b"] = a

	fmt.Println("Part 1:", a)
	fmt.Println("Part 2:", p.Eval("a"))
}
