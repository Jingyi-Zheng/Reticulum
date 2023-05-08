package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generateadvepcoh(adv []int, num, psnum, csnum int) [][]int {
	// Create an empty slice to store the generated lists
	advlist := [][]int{}

	// Iterate over each ad ratio value in the adv slice
	for _, adratio := range adv {
		// Set the random seed based on the current time
		rand.Seed(time.Now().UnixNano())

		// Calculate the target number of elements for the current ad ratio
		target := int(float64(num) * float64(adratio) / float64(100))

		// Create an empty slice to store the selected numbers
		var selected []int

		// Continue selecting numbers until the desired target is reached
		for len(selected) < target {
			// Generate a random number between 0 and num (inclusive)
			randNum := rand.Intn(num + 1)

			// Check if the number is not already selected, not divisible by psnum, and not equal to 0
			if !contains(selected, randNum) && randNum%psnum != 1 && randNum != 0 {
				// Add the number to the selected slice
				selected = append(selected, randNum)
			}
		}

		// Add the selected slice to the advlist
		advlist = append(advlist, selected)
	}

	// Print the generated advlist
	fmt.Println(advlist)

	// Return the generated advlist
	return advlist
}

func contains(arr []int, num int) bool {
	// Iterate over each element in the arr slice
	for _, v := range arr {
		// Check if the current element is equal to the target num
		if v == num {
			// Return true if the element is found
			return true
		}
	}

	// Return false if the element is not found
	return false
}
