package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"slices"
	"strconv"
)

type Entry struct { patts, pin []string }

type Puzzle []Entry

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func parseCodes(s string) (out []string) {
	for _, code := range strings.Fields(s) {
		bytes := []byte(code)
		slices.Sort(bytes)
		out = append(out, string(bytes))
	}
	return out
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		data := strings.Split(s.Text(), "|")
		p = append(p, Entry{ parseCodes(data[0]), parseCodes(data[1]) })
	}
	check(s.Err())
	return p
}

func (p Puzzle) PartOne() (count int) {
	valid := map[int]bool{ 2: true, 3: true, 4: true, 7: true}
	for _, entry := range p {
		for _, digit := range entry.pin {
			if valid[len(digit)] {
				count++
			}
		}
	}
	return count
}

func (p Puzzle) PartTwo() (sum int) {
	for _, entry := range p {
		sum += entry.GetCode()
	}
	return sum
}

func (e Entry) GetCode() int {
	decode := e.GetMapping()
	var code strings.Builder
	for _, digit := range e.pin {
		code.WriteRune(decode[digit])
	}
	i, err := strconv.Atoi(code.String())
	check(err)
	return i
}

func mod(s1, s2 string) (count int) {
	bad := map[rune]bool{}
	for _, ch := range s2 {
		bad[ch] = true
	}
	for _, ch := range s1 {
		if !bad[ch] {
			count++
		}
	}
	return count
}

func (e Entry) GetMapping() map[string]rune {
	encode := map[rune]string{}
	decode := map[string]rune{}
	todo := map[string]struct{}{}
	for _, patt := range e.patts {
		if len(patt) == 2 {
			encode['1'] = patt
			decode[patt] = '1'
			continue
		}
		if len(patt) == 3 {
			encode['7'] = patt
			decode[patt] = '7'
			continue
		}
		if len(patt) == 4 {
			decode[patt] = '4'
			continue
		}
		if len(patt) == 7 {
			decode[patt] = '8'
			continue
		}
		todo[patt] = struct{}{}
	}
	one := encode['1']
	seven := encode['7']
	for patt := range todo {
		if mod(patt, one) == 3 {
			encode['3'] = patt
			decode[patt] = '3'
			delete(todo, patt)
		} else if mod(patt, seven) == 4 {
			encode['6'] = patt
			decode[patt] = '6'
			delete(todo, patt)
		}
	}
	three := encode['3']
	for patt := range todo {
		if mod(patt, three) == 2 {
			decode[patt] = '0'
			delete(todo, patt)
			break
		}
	}
	six := encode['6']
	for patt := range todo {
		if len(patt) == 6 {
			decode[patt] = '9'
		} else if mod(six, patt) == 1 {
			decode[patt] = '5'
		} else {
			decode[patt] = '2'
		}
		delete(todo, patt)
	}

	return decode
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	fmt.Println("Part 1:", p.PartOne())
	fmt.Println("Part 2:", p.PartTwo())
}
