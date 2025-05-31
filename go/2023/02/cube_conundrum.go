package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Round struct {
	red, green, blue int
}

type Game struct {
	id     int
	rounds []Round
}

type Bag struct {
	red, green, blue int
}

func (b Bag) CouldPlay(g Game) bool {
	for _, r := range g.rounds {
		if r.red > b.red || r.green > b.green || r.blue > b.blue {
			return false
		}
	}
	return true
}

func (g Game) MinPower() int {
	var b Bag
	for _, r := range g.rounds {
		if r.red > b.red {
			b.red = r.red
		}
		if r.green > b.green {
			b.green = r.green
		}
		if r.blue > b.blue {
			b.blue = r.blue
		}
	}
	fmt.Println("Min is", b, "for game", g)
	return b.red * b.green * b.blue
}

func LoadGames(filename string) (games []Game) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		before, after, match := strings.Cut(line, ": ")
		if !match {
			panic("Bad format: " + line)
		}
		id, err := strconv.Atoi(before[5:])
		check(err)
		rounds := parseRounds(after)
		games = append(games, Game{id, rounds})
	}
	check(s.Err())
	return games
}

func parseRounds(text string) (rounds []Round) {
	for _, v := range strings.Split(text, "; ") {
		var r, g, b int
		for _, c := range strings.Split(v, ", ") {
			cnt, colour, match := strings.Cut(c, " ")
			if !match {
				panic("Bad colour format: " + c)
			}
			count, err := strconv.Atoi(cnt)
			check(err)
			switch colour {
			case "red":
				r = count
			case "green":
				g = count
			case "blue":
				b = count
			default:
				panic("Bad colour format: " + c)
			}
		}
		rounds = append(rounds, Round{r, g, b})
	}
	return rounds
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	bag := Bag{12, 13, 14}
	games := LoadGames(flag.Arg(0))

	var sum, power int
	for i, game := range games {
		if bag.CouldPlay(game) {
			sum += i + 1
		}
		power += game.MinPower()
	}

	fmt.Println("Part 1 answer:", sum)
	fmt.Println("Part 2 answer:", power)
}
