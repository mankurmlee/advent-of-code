package main

import (
	"flag"
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
)

type Race struct {
	time   int
	record int
}

func (r Race) WaysToWin() int {
	i := 1
	for i * (r.time - i) <= r.record {
		i++
	}
	return r.time - 1 - 2 * (i - 1)
}

type Puzzle struct {
	races []Race
}

func (p Puzzle) WinProduct() (product int) {
	product = 1
	for _, r := range p.races {
		product *= r.WaysToWin()
	}
	return product
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	s := bufio.NewScanner(f)
	s.Scan()
	timeData := strings.Fields(s.Text())
	check(s.Err())
	s.Scan()
	recordData := strings.Fields(s.Text())
	check(s.Err())
	for i := 1; i < len(timeData); i++ {
		p.races = append(p.races, Race{
			Atoi(timeData[i]),
			Atoi(recordData[i]),
		})
	}
	return p
}

func load2(filename string) (p Puzzle) {
	r := strings.NewReplacer(" ", "")
	f, err := os.Open(filename)
	check(err)
	s := bufio.NewScanner(f)
	s.Scan()
	timeData := strings.Split(r.Replace(s.Text()), ":")
	check(s.Err())
	s.Scan()
	recordData := strings.Split(r.Replace(s.Text()), ":")
	check(s.Err())
	p.races = append(p.races, Race{
		Atoi(timeData[1]),
		Atoi(recordData[1]),
	})
	return p
}

func Atoi(A string) (i int) {
	i, err := strconv.Atoi(A)
	check(err)
	return i
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	filename := flag.Arg(0)
	p := load(filename)
	prod := p.WinProduct()
	fmt.Println("Answer to part 1 is", prod)

	q := load2(filename)
	prod2 := q.WinProduct()
	fmt.Println("Answer to part 2 is", prod2)
}
