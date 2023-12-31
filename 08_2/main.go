package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Need arg")
	}
	contentBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	content := string(contentBytes)
	contentLines := strings.Split(content, "\n")
	directions := contentLines[0]
	nodes := make(map[string][]string)
	startNodes := make([]string, 0)
	for _, line := range contentLines[1:] {
		if len(line) == 0 {
			continue
		}
		lineParts := strings.Split(line, "=")
		if len(lineParts) != 2 {
			log.Printf("warning: line %s\n", line)
			continue
		}
		left := strings.TrimSpace(lineParts[0])
		if left[len(left)-1] == 'A' {
			startNodes = append(startNodes, left)
		}
		right := strings.TrimSpace(lineParts[1])
		right = strings.TrimLeft(right, "(")
		right = strings.TrimRight(right, ")")
		rightParts := strings.Split(right, ",")
		nodes[left] = make([]string, 0)
		for _, part := range rightParts {
			nodes[left] = append(nodes[left], strings.TrimSpace(part))
		}
	}

	steps := 0
	currNodes := startNodes
	numEnded := 0
	starts := make([]int, len(startNodes))
	intervals := make([][]int, len(startNodes))
	for i := range intervals {
		intervals[i] = make([]int, 0)
	}
	isDone := make([]bool, len(startNodes))
	for {
		for _, dir := range directions {
			steps++
			for i, curr := range currNodes {
				nextNodes, ok := nodes[curr]
				if !ok {
					log.Fatalf("%s\n", curr)
				}
				if dir == 'L' {
					currNodes[i] = nextNodes[0]
				} else if dir == 'R' {
					currNodes[i] = nextNodes[1]
				} else {
					log.Fatalf("bad dir: %c\n", dir)
				}
				if currNodes[i][len(currNodes[i])-1] == 'Z' && !isDone[i] {
					intervals[i] = append(intervals[i], steps)
					//log.Printf("intervals[%d]: adding %d, %d\n", i, steps, len(intervals[i]))
					x, checked := checkForCycles(intervals[i])
					if checked != nil {
						starts[i] = x
						intervals[i] = checked
						log.Printf("%d %d %d", i, steps, checked[0])
						isDone[i] = true
						numEnded++
					}
				}
			}
			if numEnded == len(currNodes) {
				break
			}
		}
		if numEnded == len(currNodes) {
			break
		}
	}

	for i, interval := range intervals {
		fmt.Printf("intervals[%d]: (%d) %d", i, starts[i], interval[0])
		if len(intervals) > 1 {
			for _, n := range interval[1:] {
				fmt.Printf(", %d", starts[i]+n)
			}
		}
		fmt.Printf("\n")
	}

}

func checkForCycles(nums []int) (int, []int) {
	for i := 0; i < len(nums)/2; i++ {
		if (len(nums)-i)%4 != 0 {
			continue
		}
		segLen := (len(nums) - i) / 4
		offsets := make([]int, 4)
		bases := make([]int, 4)
		for k := range offsets {
			offsets[k] = i + segLen*k
			fmt.Printf("%d: %d", k, nums[offsets[k]])
			for l := 1; l < segLen; l++ {
				fmt.Printf(", %d", nums[offsets[k]+l])
			}
			fmt.Printf("\n")
		}
		if i > 0 {
			bases[0] = nums[i-1]
		}
		for k := 1; k < len(bases); k++ {
			bases[k] = nums[offsets[k]-1]
		}
		isDone := true
		for j := 0; j < segLen; j++ {
			first := nums[offsets[0]+j] - bases[0]
			for k := 1; k < 4; k++ {
				later := nums[offsets[k]+j] - bases[k]
				//log.Printf("cmp: %d to %d\n", offsets[0]+j, offsets[k]+j)
				if first != later {
					log.Printf("difference: first: %d, %d, later: %d, %d\n", bases[0], first, bases[k], later)
					isDone = false
					break
				}
			}
			if !isDone {
				break
			}
		}
		if isDone {
			for l := 0; l < segLen; l++ {
				nums[i+l] -= bases[0]
			}
			return bases[0], nums[i : i+segLen]
		}
	}
	return -1, nil
}
