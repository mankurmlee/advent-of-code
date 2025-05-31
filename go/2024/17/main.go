package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

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

func findInts(s string) (nums []int) {
	r := regexp.MustCompile(`\-?\d+`)
	for _, v := range r.FindAllString(s, -1) {
		nums = append(nums, atoi(v))
	}
	return nums
}

func readChunkedFile(filename string) (chunks [][]string) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	var chunk []string
	for s.Scan() {
		l := s.Text()
		if l == "" {
			if len(chunk) > 0 {
				chunks = append(chunks, chunk)
				chunk = []string{}
			}
			continue
		}
		chunk = append(chunk, l)
	}
	check(s.Err())
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks
}

func main() {
	flag.Parse()
	var p Program
	p.Load(flag.Arg(0))
	p.Run(0, false)
	p.PartTwo()
}

type Registers struct {
	ip, a, b, c int
	out         []int
}

func (r *Registers) Print() {
	strs := make([]string, len(r.out))
	for i, num := range r.out {
		strs[i] = strconv.Itoa(num)
	}
	fmt.Println(strings.Join(strs, ","))
}

func (r *Registers) Execute(opcode int, operand int) {
	combo := r.ComboOperand(operand)
	switch opcode {
	case 0:
		r.a >>= combo
	case 1:
		r.b ^= operand
	case 2:
		r.b = combo % 8
	case 3:
		if r.a != 0 {
			r.ip = operand
		}
	case 4:
		r.b ^= r.c
	case 5:
		r.out = append(r.out, combo%8)
	case 6:
		r.b = r.a >> combo
	case 7:
		r.c = r.a >> combo
	default:
		panic("Unexpected opcode")
	}
}

func (r *Registers) ComboOperand(operand int) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return r.a
	case 5:
		return r.b
	case 6:
		return r.c
	default:
		panic("Unexpected operand")
	}
}

type SearchState struct {
	sum, length int
}

type Program struct {
	Registers
	data []int
}

func (p Program) PartTwo() {
	var results []int
	target := make([]int, len(p.data))
	tLength := len(target)
	copy(target, p.data)
	slices.Reverse(target)
	q := []SearchState{{}}
	for len(q) > 0 {
		i := len(q) - 1
		s := q[i]
		q = q[:i]
		if s.length == tLength {
			results = append(results, s.sum)
			continue
		}
		q = append(q, p.FindOctals(s, target)...)
	}
	if len(results) == 0 {
		fmt.Println("Program does not copy")
		return
	}
	slices.Sort(results)
	fmt.Println(results[0])
}

func (p Program) FindOctals(s SearchState, tar []int) (found []SearchState) {
	a := s.sum
	t := tar[s.length]
	for i := range 8 {
		g := (a << 3) + i
		h := p.Run(g, true)
		if len(h) != 1 {
			panic("Unexpected number of results")
		}
		if h[0] == t {
			found = append(found, SearchState{g, s.length + 1})
		}
	}
	return found
}

func (p Program) Run(init int, step bool) []int {
	copy := p.Registers
	r := &copy
	if init > 0 {
		r.a = init
	}
	end := len(p.data)
	for r.ip < end {
		oldIp := r.ip
		r.Execute(p.data[r.ip], p.data[r.ip+1])
		if oldIp == r.ip {
			r.ip += 2
		}
		if step && len(r.out) == 1 {
			return r.out
		}
	}
	r.Print()
	return r.out
}

func (p *Program) Load(filename string) {
	chunks := readChunkedFile(filename)
	p.a = findInts(chunks[0][0])[0]
	p.b = findInts(chunks[0][1])[0]
	p.c = findInts(chunks[0][2])[0]
	p.data = findInts(chunks[1][0])
}
