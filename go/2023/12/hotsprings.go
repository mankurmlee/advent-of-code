package main

import (
	"flag"
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
)

type LUT map[string]int

type Record struct {
	pattern  string
	checksum []int
}

func (l *LUT) CountArrangements(r Record) int {
	var total int

	// Check cache
	count, ok := (*l)[fmt.Sprint(len(r.pattern), len(r.checksum))]
	if ok {
		return count
	}

	dof := len(r.pattern) - r.GetMinLen()

	broken := r.checksum[0]
	for i := 0; i < dof+1; i++ {
		if strings.Contains(r.pattern[0:i], "#") {
			continue
		}
		if strings.Contains(r.pattern[i:i+broken], ".") {
			continue
		}
		j := i + broken
		if len(r.checksum) > 1 {
			if r.pattern[j:j+1] == "#" {
				continue
			}
			j++
		}
		tailPatt := r.pattern[j:]
		tailChek := r.checksum[1:]

		if len(tailChek) == 0 {
			if !strings.Contains(tailPatt, "#") {
				total += 1
			}
			continue
		}

		if len(tailChek) == 1 {
			if  !strings.Contains(tailPatt, ".") &&
				!strings.Contains(tailPatt, "#") {
				total += len(tailPatt) - tailChek[0] + 1
				continue
			}
		}

		tail := Record{tailPatt, tailChek}

		//~ fmt.Println("   ", r, "i =", i, "=>", tail)
		total += l.CountArrangements(tail)
	}

	// Cache result
	s := fmt.Sprint(len(r.pattern), len(r.checksum))
	(*l)[s] = total

	return total
}

func (r Record) GetMinLen() (min int) {
	for _, c := range r.checksum {
		min += c + 1
	}
	return min - 1
}

func unfold(folded []Record) (unfolded []Record) {
	var pattern []string
	var checksum []int

	for _, r := range folded {
		pattern  = []string{}
		checksum = []int{}
		for i := 0; i < 5; i++ {
			pattern = append(pattern, r.pattern)
			checksum = append(checksum, r.checksum...)
		}
		unfolded = append(unfolded, Record{
			strings.Join(pattern, "?"),
			checksum,
		})
	}

	return unfolded
}

func parseRecord(line string) Record {
	var checksum []int
	s := strings.Fields(line)
	for _, str := range strings.Split(s[1], ",") {
		broken, err := strconv.Atoi(str)
		check(err)
		checksum = append(checksum, broken)
	}
	return Record{s[0], checksum}
}

func load(filename string) (records []Record) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		records = append(records, parseRecord(s.Text()))
	}
	check(s.Err())

	return records
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var total int
	var l LUT

	flag.Parse()
	records := load(flag.Arg(0))

	total = 0
	for _, r := range records {
		l = make(map[string]int)
		count := l.CountArrangements(r)
		total += count
		//~ fmt.Println(r, "-", count, "arrangements")
	}
	fmt.Println("Part 1 answer is", total)

	unfolded := unfold(records)

	total = 0
	for _, r := range unfolded {
		l = make(map[string]int)
		count := l.CountArrangements(r)
		total += count
		//~ fmt.Println("Pattern", i+1, ":", r, "-", count, "arrangements")
	}
	fmt.Println("Part 2 answer is", total)
}
