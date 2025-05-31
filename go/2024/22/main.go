package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
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

func load(filename string) (out []int) {
	for _, l := range readFile(filename) {
		out = append(out, atoi(l))
	}
	return out
}

func main() {
	flag.Parse()
	nums := load(flag.Arg(0))

	var sum int
	tree := make(map[[4]int]int)
	for _, i := range nums {
		n, b := generate(i, 2000)
		sum += n
		for k, v := range b {
			tree[k] += v
		}
	}
	fmt.Println(sum)
	var max int
	for _, v := range tree {
		if v > max {
			max = v
		}
	}
	fmt.Println(max)
}

func generate(i, n int) (o int, bananas map[[4]int]int) {
	var buf Loop
	bananas = make(map[[4]int]int)
	for range n {
		o = next(i)
		a := i % 10
		b := o % 10
		i = o
		buf.Push(b - a)
		if buf.index < 4 {
			continue
		}
		k := buf.Key()
		if _, ok := bananas[k]; !ok {
			bananas[k] = b
		}
	}
	return o, bananas
}

func next(n int) int {
	n = mixnprune(n, n<<6)
	n = mixnprune(n, n>>5)
	n = mixnprune(n, n<<11)
	return n
}

func mixnprune(n, i int) int {
	n ^= i
	n %= 16777216
	return n
}

type Loop struct {
	data  [4]int
	index int
}

func (l *Loop) Push(x int) {
	l.data[l.index%4] = x
	l.index++
}

func (l *Loop) Key() [4]int {
	if l.index < 4 {
		panic("Not enough entries")
	}
	i := l.index
	return [4]int{
		l.data[i%4],
		l.data[(i+1)%4],
		l.data[(i+2)%4],
		l.data[(i+3)%4],
	}
}
