package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func SolveDay1() {
	// Read the entire file into memory
	data, err := os.ReadFile("day1_input.txt")
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	// Split the file content into lines
	lines := strings.Split(string(data), "\n")

	fmt.Println("Lines:", len(lines))

	var left []int
	var right []int

	for _, line := range lines {

		nums := strings.Split(strings.Replace(line, "\r", "", 1), "   ")

		// fmt.Println(line)
		// fmt.Println(len(nums), nums[1])

		leftInt, leftErr := strconv.Atoi(nums[0])
		if leftErr != nil {
			panic(leftErr)
		}

		left = append(left, leftInt)

		rightInt, rightErr := strconv.Atoi(nums[1])
		if rightErr != nil {
			panic(rightErr)
		}

		right = append(right, rightInt)
	}

	rightMap := make(map[int]int)

	for _, l := range left {
		_, ok := rightMap[l]

		if !ok {
			rightMap[l] = count(right, l)
		}
	}

	if len(left) != len(right) {
		panic("Arrays are not of equal size")
	}

	// fmt.Println(rightMap)

	var sum int

	for i := range left {
		l := left[i]
		sum += l * rightMap[l]
	}

	fmt.Println("The sum is:", sum)
}

func count(arr []int, target int) int {
	result := 0

	for _, num := range arr {
		if num == target {
			result++
		}
	}

	return result
}
