package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func SolveDay3() {
	// Read the entire file into memory
	data, err := os.ReadFile("day3_input.txt")
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	// Compile the regex pattern
	regex := regexp.MustCompile(`(mul\(\d{1,3},\d{1,3}\))|(don't\(\))|(do\(\))`)

	// Find all matches, -1 means no limit
	matches := regex.FindAllString(string(data), -1)

	result := 0
	enabled := true

	// Print all matches
	for _, match := range matches {
		// fmt.Println(match)

		if match == "do()" {
			enabled = true
		} else if match == "don't()" {
			enabled = false
		} else if enabled {
			match = strings.ReplaceAll(match[:len(match)-1], "mul(", "")
			// fmt.Println(match)
			numsArr := strings.Split(match, ",")

			result += ParseInt(numsArr[0]) * ParseInt(numsArr[1])
		}
	}

	fmt.Println("The total result is:", result)
}
