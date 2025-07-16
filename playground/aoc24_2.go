package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func SolveDay2() {
	// Read the entire file into memory
	data, err := os.ReadFile("day3_input.txt")
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	// Split the file content into lines
	lines := strings.Split(string(data), "\n")

	fmt.Println("Lines:", len(lines))

	safeReportsCount := 0

	for _, line := range lines {
		levels := strings.Split(line, " ")

		if CheckLevels(levels) || slices.ContainsFunc(GetCombinations(levels), CheckLevels) {
			safeReportsCount++
		}
	}

	fmt.Println("Number of safe reports", safeReportsCount)
}

func GetCombinations(levels []string) [][]string {
	result := make([][]string, len(levels))

	for i := range levels {
		levelsCopy := make([]string, len(levels))
		copy(levelsCopy, levels)

		result[i] = slices.Delete(levelsCopy, i, i+1)
	}

	return result
}

func CheckLevels(levels []string) bool {
	increasingInitialised := false
	var increasing bool
	prevLevel := ParseInt(levels[0])

	for _, levelString := range levels[1:] {
		level := ParseInt(levelString)

		diff := level - prevLevel

		//compare level with previous level and see if its within 3
		if math.Abs(float64(diff)) > 3 || diff == 0 {
			return false
		}

		//initialise increasing bool
		if !increasingInitialised {
			if diff > 0 {
				increasing = true
			} else {
				increasing = false
			}

			increasingInitialised = true
		}

		//make sure that increasing never flips more than once
		if (increasing && diff < 0) || !increasing && diff > 0 {
			return false
		}

		prevLevel = level
	}

	return true
}

func ParseInt(s string) int {
	result, err := strconv.Atoi(strings.ReplaceAll(s, "\r", ""))
	if err != nil {
		panic(err)
	}

	return result
}
