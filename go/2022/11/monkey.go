package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"slices"
)

type WorryLevels []int

type Monkey struct {
	index       int
	nums        []int
	worries     []WorryLevels
	operator    string
	operand     string
	test        int
	istrue      int
	isfalse     int
	pops        int
}

type Notes []Monkey

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return int(i)
}

func load(filename string) (p Notes) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)

	for {
		m := parseMonkey(s)
		if m.test == 0 {
			break
		}
		m.index = len(p)
		p = append(p, m)
	}

	for i, m := range p {
		for _, item := range m.nums {
			m.worries = append(m.worries, p.GetWorry(item))
		}
		p[i] = m
	}

	return p
}

func parseMonkey(s *bufio.Scanner) (m Monkey) {
	if !s.Scan() {
		check(s.Err())
		return m
	}

	s.Scan()
	check(s.Err())
	r := strings.NewReplacer(",", " ")
	for _, itemData := range strings.Fields(r.Replace(s.Text()))[2:] {
		m.nums = append(m.nums, Atoi(itemData))
	}

	s.Scan()
	check(s.Err())
	operation := strings.Fields(s.Text())
	m.operator = operation[4]
	m.operand = operation[5]

	s.Scan()
	check(s.Err())
	testData := strings.Fields(s.Text())
	m.test = Atoi(testData[3])

	s.Scan()
	check(s.Err())
	trueData := strings.Fields(s.Text())
	m.istrue = Atoi(trueData[5])

	s.Scan()
	check(s.Err())
	falseData := strings.Fields(s.Text())
	m.isfalse = Atoi(falseData[5])

	s.Scan()
	check(s.Err())

	return m
}

func (n Notes) PrintNums() {
	for i, m := range n {
		fmt.Printf("Monkey %d: %v\n", i, m.nums)
	}
}

func (n Notes) ShowPops() {
	for i, m := range n {
		fmt.Printf("Monkey %d: %v\n", i, m.pops)
	}
}

func (n *Notes) DoNums() {
	for i, _ := range *n {
		n.UpdateNums(&(*n)[i])
	}
}

func (n *Notes) UpdateNums(m *Monkey) {
	old := *n
	for len(m.nums) > 0 {
		item := m.nums[0]
		m.nums = m.nums[1:]
		m.pops++
		switch m.operator {
		case "+":
			item += Atoi(m.operand)
		case "*":
			if m.operand == "old" {
				item *= item
			} else {
				item *= Atoi(m.operand)
			}
		}
		item /= 3
		next := m.isfalse
		if item % m.test == 0 {
			next = m.istrue
		}
		old[next].nums = append(old[next].nums, item)
	}
	*n = old
}

func (n Notes) GetMonkeyBusiness() int {
	var pops []int
	for _, m := range n {
		pops = append(pops, m.pops)
	}
	slices.Sort(pops)
	slices.Reverse(pops)
	return pops[0] * pops[1]
}

func (n Notes) GetWorry(x int) (w WorryLevels) {
	for _, m := range n {
		w = append(w, x % m.test)
	}
	return w
}

func (n *Notes) DoWorry() {
	for i, _ := range *n {
		n.UpdateWorry(&(*n)[i])
	}
}


func (n *Notes) UpdateWorry(m *Monkey) {
	old := *n
	for len(m.worries) > 0 {
		item := m.worries[0]
		m.worries = m.worries[1:]
		m.pops++

		var operand int
		if m.operand != "old" {
			operand = Atoi(m.operand)
		}
		switch m.operator {
		case "+":
			for i, _ := range item {
				item[i] += operand
				item[i] %= old[i].test
			}
		case "*":
			if m.operand == "old" {
				for i, w := range item {
					item[i] *= w
					item[i] %= old[i].test
				}
			} else {
				for i, _ := range item {
					item[i] *= operand
					item[i] %= old[i].test
				}
			}
		}
		next := m.isfalse
		if item[m.index] == 0 {
			next = m.istrue
		}
		old[next].worries = append(old[next].worries, item)
	}
	*n = old
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	for i := 0; i < 20; i++ {
		p.DoNums()
	}
	//~ p.PrintNums()
	fmt.Println("Part 1:", p.GetMonkeyBusiness())

	// Reset counters
	for i, _ := range p {
		p[i].pops = 0
	}

	for i := 0; i < 10000; i++ {
		p.DoWorry()
	}
	p.ShowPops()
	fmt.Println("Part 2:", p.GetMonkeyBusiness())
}
