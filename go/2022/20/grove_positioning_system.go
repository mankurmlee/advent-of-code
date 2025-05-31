package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strconv"
)

type Node struct {
	value int
	next  *Node
	prev  *Node
}

type List struct {
	nodes []*Node
	zero  *Node
}

func (l *List) Append(x int) {
	n := len(l.nodes)
	o := Node{ value : x }

	if x == 0 {
		l.zero = &o
	}

	if n == 0 {
		o.prev = &o
		o.next = &o
	} else {
		next := l.nodes[0]
		prev := l.nodes[n-1]

		o.next = next
		o.prev = prev

		next.prev = &o
		prev.next = &o
	}

	l.nodes = append(l.nodes, &o)
}

func (l List) GetCoordinate(n int) int {
	n %= len(l.nodes)
	o := l.zero
	for i := 0; i < n; i++ {
		o = o.next
	}
	return o.value
}

func (l List) MixNode(o *Node) {
	if o.value == 0 {
		return
	}
	prev := o.prev
	next := o.next
	prev.next = next
	next.prev = prev
	if o.value > 0 {
		n := o.value
		n %= len(l.nodes) - 1
		for i := 0; i < n; i++ {
			prev = next
			next = next.next
		}
	} else {
		n := -o.value
		n %= len(l.nodes) - 1
		for i := 0; i < n; i++ {
			next = prev
			prev = prev.prev
		}
	}
	o.prev = prev
	o.next = next
	prev.next = o
	next.prev = o
}

func (l List) Mix() {
	for _, o := range l.nodes {
		l.MixNode(o)
	}
}

func (l *List) Decrypt(key int) {
	for _, o := range l.nodes {
		o.value *= key
	}
}

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

func load(filename string) (l List) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		l.Append(atoi(s.Text()))
	}
	check(s.Err())
	return l
}

func main() {
	flag.Parse()
	l := load(flag.Arg(0))

	l.Mix()

	sum := l.GetCoordinate(1000)
	sum += l.GetCoordinate(2000)
	sum += l.GetCoordinate(3000)

	l = load(flag.Arg(0))
	l.Decrypt(811589153)
	for i := 0; i < 10; i++ {
		l.Mix()
	}

	decrypted := l.GetCoordinate(1000)
	decrypted += l.GetCoordinate(2000)
	decrypted += l.GetCoordinate(3000)

	fmt.Println("Part 1:", sum)
	fmt.Println("Part 2:", decrypted)
}
