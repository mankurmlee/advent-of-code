package main

import (
	"fmt"
	"bufio"
	"os"
	"unicode"
	"strconv"
)

func main() {
	f, e := os.Open("input.txt")
	if e != nil {
		fmt.Println("Error opening file:", e)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var sum int
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		num, err := extractDigits(line)
		if err != nil {
			fmt.Println("Error extracting digits:", err)
			continue
		}
		fmt.Println(num)
		sum += num
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	fmt.Println("Total:", sum)
}

func extractDigits(text string) (int, error) {
	// Initialize variables to store the first and last digits
	var firstDigit, lastDigit int

	// Iterate through the characters in the text
	for _, char := range text {
		// Check if the character is a digit
		if unicode.IsDigit(char) {
			// Convert the digit to an integer
			digit, err := strconv.Atoi(string(char))
			if err != nil {
				return 0, err
			}

			// If it's the first digit encountered, store it
			if firstDigit == 0 {
				firstDigit = digit
			}

			// Always update the last digit encountered
			lastDigit = digit
		}
	}

	// Combine the first and last digits into a single number
	result := firstDigit*10 + lastDigit

	return result, nil
}
