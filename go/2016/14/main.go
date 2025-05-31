package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type Triplet struct {
	index int
	char  rune
	isKey bool
}

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

func main() {
	flag.Parse()
	salt := readFile(flag.Arg(0))[0]
	getNthKey(salt, md5hash)
	getNthKey(salt, stretchedHash)
}

func getNthKey(salt string, hashFunc func(string) string) {
	keysFound := 0
	var q []Triplet
	for index := 0; true; index++ {
		// Prune queue
		for len(q) > 0 && q[0].index+1000 < index {
			if q[0].isKey {
				keysFound++
				if keysFound == 64 {
					fmt.Println(q[0].index)
					return
				}
			}
			q = q[1:]
		}
		// Parse string for triplets and quintuplets
		s := fmt.Sprint(salt, index)
		h := hashFunc(s)
		ok, tChar, qs := parse(h)
		if !ok {
			continue
		}
		// Mark triplets with quintuplets as keys
		for i, t := range q {
			if t.isKey {
				continue
			}
			if slices.Contains(qs, t.char) {
				q[i].isKey = true
			}
		}
		q = append(q, Triplet{index, tChar, false})
	}
}

func parse(s string) (ok bool, t rune, qs []rune) {
	var last rune
	var reps int
	for _, c := range s {
		if c != last {
			reps = 0
			last = c
			continue
		}
		reps++
		if !ok && reps == 2 {
			ok = true
			t = last
		}
		if reps == 4 && !slices.Contains(qs, c) {
			qs = append(qs, c)
		}
	}
	return ok, t, qs
}

func md5hash(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func stretchedHash(s string) string {
	for range 2017 {
		s = md5hash(s)
	}
	return s
}
