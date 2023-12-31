package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Args struct {
	nIndex int
	sIndex int
}

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
	for i, line := range strings.Split(contents, "\n") {
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
		sOrig := s
		numsOrig := []int(nums)
		for i := 1; i < 5; i++ {
			s += "?" + sOrig
			nums = append(nums, numsOrig...)
		}
		if i >= 0 {
			singleResult := solve(s, nums)
			log.Printf("single result: %d", singleResult)
			result += singleResult
		}
	}
	log.Printf("%d", result)
}

func solve(s string, nums []int) (result int) {
	nums = append(nums, 0)
	copy(nums[1:], nums[0:])
	nums[0] = -1
	memo := make(map[string]int)
	sLen := len(s)
	stack := make([]Args, 0)
	limits := make([]int, len(nums))
	limit := sLen
	for i := len(nums) - 1; i >= 0; i-- {
		limit -= nums[i]
		limits[i] = limit
		limit--
	}
	base := Args{nIndex: 0, sIndex: 0}
	stack = append(stack, base)
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		//log.Printf("top: nIndex: %d sIndex: %d", top.nIndex, top.sIndex)
		currNum := nums[top.nIndex]
		isValid := true
		for i := 0; i < currNum; i++ {
			sIndex := top.sIndex + i
			if s[sIndex] == '.' {
				isValid = false
				break
			}
		}
		if !isValid {
			//log.Printf("not valid")
			key := top.String()
			memo[key] = 0
			stack = stack[:len(stack)-1]
			continue
		}
		if top.nIndex < len(nums)-1 {
			oIndex := top.nIndex + 1
			oLimit := limits[oIndex]
			//log.Printf("oIndex: %d oLimit: %d", oIndex, oLimit)
			r := 0
			isValid := true
			for i := top.sIndex + currNum; i < oLimit; i++ {
				if i >= 0 && s[i] == '#' {
					//log.Printf("hit # at %d", i)
					break
				}
				next := Args{nIndex: oIndex, sIndex: i + 1}
				key := next.String()
				ro, ok := memo[key]
				if !ok {
					//log.Printf("recursing with nIndex: %d sIndex: %d", oIndex, i+1)
					stack = append(stack, next)
					isValid = false
					break
				}
				//log.Printf("got memoized result: %d", ro)
				r += ro
			}
			if !isValid {
				continue
			}
			//log.Printf("finished: nIndex: %d sIndex: %d result: %d", top.nIndex, top.sIndex, result)
			key := top.String()
			memo[key] = r
			stack = stack[:len(stack)-1]
		} else { // at the rightmost num
			//log.Printf("At rightmost num")
			isValid := true
			for i := top.sIndex + currNum; i < sLen; i++ {
				if s[i] == '#' {
					isValid = false
					break
				}
			}
			key := top.String()
			if !isValid {
				//log.Printf("Not valid")
				memo[key] = 0
			} else {
				//log.Printf("Valid")
				memo[key] = 1
			}
			stack = stack[:len(stack)-1]
		}
	}
	result, _ = memo[base.String()]
	return
}

func resetIndices(sLen int, nums []int, indices []int, nIndex int) bool {
	numsLen := len(nums)
	for {
		offset := indices[nIndex] + nums[nIndex] + 1
		nIndex++
		if nIndex >= numsLen {
			return offset > sLen
		}
		indices[nIndex] = offset
	}
}

func incrementIndices(sLen int, nums []int, indices []int, nIndex int) (newIndex int, ok bool) {
	for nIndex >= 0 {
		indices[nIndex]++
		if resetIndices(sLen, nums, indices, nIndex) {
			newIndex = nIndex
			ok = true
			return
		}
		nIndex--
		continue
	}
	newIndex = -1
	ok = false
	return
}

func (a Args) String() string {
	return fmt.Sprintf("%d_%d", a.nIndex, a.sIndex)
}
