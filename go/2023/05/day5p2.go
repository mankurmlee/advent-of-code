package main


import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	"flag"
	"runtime"
)

type Range struct {
	start int
	end   int
}

type Rule struct {
	first  int
	last   int
	offset int
}

type Almanac [][]Rule

func (a Almanac) GetLocation(seed int) int {
	s := seed
	for _, v := range a {
		for _, r := range v {
			if s >= r.first && s <= r.last {
				s += r.offset
				break
			}
		}
	}
	return s
}


// Parses the input file
func parseInput(filename string) ([]Range, Almanac) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	s := bufio.NewScanner(f)

	var seeds   []Range
	var alm     Almanac
	var row     []Rule

	for s.Scan() {
		l := s.Text()

		// Empty line
		if l == "" {
			continue
		}

		// Seeds
		if strings.HasPrefix(l, "seeds:") {
			parts := strings.Split(l[7:], " ")
			num := len(parts)
			for i := 0; i < num; i += 2 {
				a, err := strconv.Atoi(parts[i])
				check(err)
				b, err := strconv.Atoi(parts[i+1])
				check(err)
				r := Range{
					start: a,
					end:   a+b,
				}
				seeds = append(seeds, r)
			}
			continue
		}

		// Maps
		if strings.HasSuffix(l, "map:") {
			alm = append(alm, row)
			row = []Rule{}
			continue
		}

		// Ranges
		rule := parseRule(l)
		row = append(row, rule)
	}

	// Make sure to add the last row
	alm = append(alm, row)

	check(s.Err())

	return seeds, alm[1:]
}


// Parses and returns a rule
func parseRule(text string) Rule {
	parts := strings.Split(text, " ")

	dest, err := strconv.Atoi(parts[0])
	check(err)

	src, err := strconv.Atoi(parts[1])
	check(err)

	length, err := strconv.Atoi(parts[2])
	check(err)

	rule := Rule{
		first:  src,
		last:   src+length-1,
		offset: dest-src,
	}

	return rule
}


// Lazy error checking
func check(err error) {
	if err != nil {
		panic(err)
	}
}


func findMin(id int, a Almanac, seeds []int, ch chan int) {
	var loc int

	min := a.GetLocation(seeds[0])

	for _, s := range seeds {
		loc = a.GetLocation(s)
		if loc < min {
			min = loc
		}
	}

	ch <- min
}


func main() {
	flag.Parse()
	filename := flag.Args()[0]

	// Get the number of CPUs
	nCPU := runtime.NumCPU()
	fmt.Println("Num CPUs", nCPU)

	ch := make(chan int, nCPU)
	defer close(ch)

	// Parse input file
	rSeeds, alm := parseInput(filename)

	// Count number of seeds
	var nSeeds int
	for _, r := range rSeeds {
		nSeeds += r.end - r.start
	}
	fmt.Println("Number of seeds:", nSeeds)

	// Populate the seeds
	seeds := make([]int, nSeeds)
	var idx int
	for _, r := range rSeeds {
		for s := r.start; s < r.end; s++ {
			seeds[idx] = s
			idx++
		}
	}

	for i := 0; i < nCPU; i++ {
		end := (i+1)*(nSeeds/nCPU)
		if i == nCPU - 1 {
			end = nSeeds
		}
		go findMin(i, alm, seeds[i*(nSeeds/nCPU):end], ch)
	}

	var min int
	for i := 0; i < nCPU; i++ {
		loc := <-ch
		if i == 0 || loc < min {
			min = loc
		}
	}

	fmt.Println("Lowest Location:", min)
}
