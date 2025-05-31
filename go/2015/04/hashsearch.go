package main

import (
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

type Puzzle struct {
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func calculateMD5(input string) string {
	// Create an MD5 hash object
	hasher := md5.New()

	// Write the input data to the hash object
	hasher.Write([]byte(input))

	// Get the final hash as a byte slice
	hashBytes := hasher.Sum(nil)

	// Convert the byte slice to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}

func main() {
	key := "bgvyzdsv"

	var n int
	for {
		hash := calculateMD5(key + strconv.Itoa(n))
		if hash[:5] == "00000" {
			break
		}
		n++
	}
	fmt.Println("Part 1:", n)

	for {
		hash := calculateMD5(key + strconv.Itoa(n))
		if hash[:6] == "000000" {
			break
		}
		n++
	}
	fmt.Println("Part 2:", n)
}
