package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Board []int

type Puzzle struct {
	draworder   []int
	boards      []Board
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

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)

	// Draw numbers
	s.Scan()
	check(s.Err())
	for _, num := range strings.Split(s.Text(), ",") {
		p.draworder = append(p.draworder, atoi(num))
	}
	s.Scan()

	// Boards
	var board Board
	for s.Scan() {
		line := s.Text()
		if line == "" {
			p.boards = append(p.boards, board)
			board = Board{}
		}
		for _, num := range strings.Fields(line) {
			board = append(board, atoi(num))
		}
	}
	check(s.Err())
	p.boards = append(p.boards, board)

	return p
}

func (p Puzzle) WinnerScore() int {
	drawn := map[int]bool{}
	for i, num := range p.draworder {
		drawn[num] = true
		if i < 5 {
			continue
		}
		for _, board := range p.boards {
			if board.Bingo(drawn) {
				return board.Score(drawn) * num
			}
		}
	}
	return 0
}

func (b Board) Bingo(drawn map[int]bool) bool {
	var rowwins, colwins bool
	for i := 0; i < 5; i++ {
		rowwins = true
		colwins = true
		for j := 0; j < 5; j++ {
			rowwins = rowwins && drawn[b[i * 5 + j]]
			colwins = colwins && drawn[b[j * 5 + i]]
			if !rowwins && !colwins {
				break
			}
		}
		if rowwins || colwins {
			return true
		}
	}
	return false
}

func (p Puzzle) LoserScore() int {
	waiting := map[int]Board{}
	for i, board := range p.boards {
		waiting[i] = board
	}
	drawn := map[int]bool{}
	for i, num := range p.draworder {
		drawn[num] = true
		if i < 5 {
			continue
		}
		for i, board := range waiting {
			if board.Bingo(drawn) {
				if len(waiting) > 1 {
					delete(waiting, i)
				} else {
					return board.Score(drawn) * num
				}
			}
		}
	}
	return 0
}

func (b Board) Score(drawn map[int]bool) (score int) {
	for _, num := range b {
		if !drawn[num] {
			score += num
		}
	}
	return score
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	fmt.Println("Part 1:", p.WinnerScore())
	fmt.Println("Part 2:", p.LoserScore())
}
