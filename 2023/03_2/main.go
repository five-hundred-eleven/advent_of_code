package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	filename := os.Args[1]
	contentsBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	result := 0
	contents := string(contentsBytes)
	grid := strings.Split(contents, "\n")
	for i, row := range grid {
		for j, col := range row {
			if col != '*' {
				continue
			}
			n := testNumbers(grid, i, j)
			if n == -1 {
				continue
			}
			log.Printf("%d\n", n)
			result += n
		}
	}
	log.Printf("%d\n", result)
}

func testNumbers(grid []string, iBase int, jBase int) int {
	numbers := make([]int, 0)
	for i := iBase - 1; i <= iBase+1; i++ {
		j := jBase - 1
		for j <= jBase+1 {
			result, _, jRight := exploreNumber(grid, i, j)
			if result == -1 {
				j++
				continue
			}
			numbers = append(numbers, result)
			j = jRight
		}
	}
	if len(numbers) == 2 {
		return numbers[0] * numbers[1]
	}
	return -1
}

func exploreNumber(grid []string, iBase int, jBase int) (result int, jLeft int, jRight int) {
	result = -1
	jLeft = jBase
	jRight = jBase
	if iBase < 0 || iBase >= len(grid) || jBase < 0 || jBase >= len(grid[iBase]) {
		return
	}
	col := grid[iBase][jBase]
	if col < '0' || col > '9' {
		return
	}
	for jLeft >= 0 && grid[iBase][jLeft] >= '0' && grid[iBase][jLeft] <= '9' {
		jLeft--
	}
	jLeft++
	for jRight < len(grid[iBase]) && grid[iBase][jRight] >= '0' && grid[iBase][jRight] <= '9' {
		jRight++
	}
	result, err := strconv.Atoi(grid[iBase][jLeft:jRight])
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	return
}

func testGrid(grid []string, i int, jLeft int, jRight int) bool {
	for iOffset := i - 1; iOffset <= i+1; iOffset++ {
		if iOffset < 0 {
			continue
		}
		if iOffset >= len(grid) {
			continue
		}
		for jOffset := jLeft - 1; jOffset <= jRight+1; jOffset++ {
			if jOffset < 0 {
				continue
			}
			if jOffset >= len(grid[iOffset]) {
				continue
			}
			if iOffset == i && jOffset >= jLeft && jOffset <= jRight {
				continue
			}
			//log.Printf("%d %d\n", iOffset, jOffset)
			g := grid[iOffset][jOffset]
			if g <= '9' && g >= '0' {
				continue
			}
			if g == '.' {
				continue
			}
			return true
		}
	}
	return false
}
