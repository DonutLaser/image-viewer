package main

import "math"

func clamp(value int32, min int32, max int32) int32 {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}

func max(value1 int32, value2 int32) int32 {
	if value1 > value2 {
		return value1
	}

	return value2
}

func min(value1 int32, value2 int32) int32 {
	if value1 < value2 {
		return value1
	}

	return value2
}

func floor(value float64) int32 {
	return int32(math.Floor(value))
}
