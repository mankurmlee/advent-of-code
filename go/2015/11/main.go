package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func IsValid(pass []byte) bool {
	one := false
	two := true
	chain := 0
	last := byte(0)
	bad := map[byte]bool{'i': true, 'o': true, 'l': true}
	pairs := map[byte]bool{}
	for _, ch := range pass {
		if ch == last+1 {
			chain++
			if chain == 2 {
				one = true
			}
		} else {
			chain = 0
		}
		if bad[ch] {
			two = false
		}
		if ch == last {
			pairs[ch] = true
		}
		last = ch
	}
	return one && two && len(pairs) >= 2
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load() []byte {
	flag.Parse()
	f, err := os.Open(flag.Arg(0))
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Scan()
	check(s.Err())
	return []byte(s.Text())
}

func main() {
	s := load()
	valid := IsValid(s)
	fmt.Println(valid)
}
