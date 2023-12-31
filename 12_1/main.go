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
		lineParts := strings.Split(line, " ")
		s := lineParts[0]
		nums := make([]int, 0)
		for _, numString := range strings.Split(lineParts[1], ",") {
			num, err := strconv.Atoi(numString)
			if err != nil {
				log.Fatalf("%s", err)
			}
			nums = append(nums, num)
		}
		result += parse(s, nums)
	}
	log.Printf("%d", result)
}

func parse(s string, nums []int) int {
	result := 0
	numUnknown := 0
	for _, b := range s {
		if b == '?' {
			numUnknown++
		}
	}
	stop := 1 << numUnknown
	for perm := 0; perm < stop; perm++ {
		bsTemp := []byte(s)
		counter := 0
		for i := 0; i < len(bsTemp); i++ {
			if s[i] == '?' {
				if (1<<counter)&perm == 0 {
					bsTemp[i] = '.'
				} else {
					bsTemp[i] = '#'
				}
				counter++
			} else {
				bsTemp[i] = s[i]
			}
		}
		if isValid(string(bsTemp), nums) {
			result++
		}
	}
	return result
}

func isValid(s string, nums []int) bool {
	i := 0
	counter := 0
	for _, b := range s {
		if b == '.' {
			if counter == 0 {
				continue
			}
			if i >= len(nums) || nums[i] != counter {
				return false
			}
			i++
			counter = 0
		} else if b == '#' {
			counter++
		} else {
			log.Printf("bad byte: %c", b)
			return false
		}
	}
	if counter > 0 {
		if i != len(nums)-1 || nums[i] != counter {
			return false
		}
		//log.Printf("valid: %s", s)
		return true
	}
	if i < len(nums) {
		return false
	}
	//log.Printf("valid: %s", s)
	return true
}
