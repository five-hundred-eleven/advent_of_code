package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need arg")
	}
	contentsBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("%s", err)
	}
	contents := string(contentsBytes)
	result := 0
	for _, line := range strings.Split(contents, "\n") {
		if len(line) == 0 {
			continue
		}
		nums := make([]int, 0)
		for _, numStr := range strings.Split(line, " ") {
			num, _ := strconv.Atoi(numStr)
			nums = append(nums, num)
		}
		numsLen := len(nums)
		numsReversed := make([]int, numsLen)
		for i := range nums {
			numsReversed[i] = nums[numsLen-i-1]
		}
		x := predict(numsReversed)
		log.Printf("%d", x)
		result += x
	}
	log.Printf("%d", result)
}

func predict(nums []int) int {
	seq := make([][]int, 1)
	seq[0] = nums
	for {
		prev := seq[len(seq)-1]
		next := make([]int, len(prev)-1)
		isDone := true
		for i := 0; i < len(next); i++ {
			next[i] = prev[i+1] - prev[i]
			if next[i] != 0 {
				isDone = false
			}
		}
		seq = append(seq, next)
		if isDone {
			break
		}
	}
	seq[len(seq)-1] = append(seq[len(seq)-1], 0)
	for i := len(seq) - 2; i >= 0; i-- {
		curr := seq[i]
		currLen := len(curr)
		prev := seq[i+1]
		prevLen := len(prev)
		seq[i] = append(curr, curr[currLen-1]+prev[prevLen-1])
	}
	return seq[0][len(seq[0])-1]
}
