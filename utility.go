package main

import "log"

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func containsString(arr []string, value string) bool {
	for _, item := range arr {
		if item == value {
			return true
		}
	}

	return false
}

func clamp(value int, min int, max int) int {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}
