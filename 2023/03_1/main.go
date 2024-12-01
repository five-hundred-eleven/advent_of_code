package main

import (
	"log"
	"strings"
	"os"
	"strconv"
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
		j := 0
		for j < len(row) {
			col := row[j]
			if col > '9' || col < '0' {
				j++
				continue
			}
			k := j
			for k < len(row) {
				col := row[k]
				if col > '9' || col < '0' {
					break
				}
				k++
			}
			n, err := strconv.Atoi(row[j:k])
			if err != nil {
				j++
				continue
			}
			if testGrid(grid, i, j, k-1) {
				result += n
			}
			j = k
		}
	}
	log.Printf("%d\n", result)
}

func testGrid(grid []string, i int, jLeft int, jRight int) bool {
	for iOffset := i-1; iOffset <= i+1; iOffset++ {
		if iOffset < 0 {
			continue
		}
		if iOffset >= len(grid) {
			continue
		}
		for jOffset := jLeft-1; jOffset <= jRight+1; jOffset++ {
			if jOffset < 0 {
				continue
			}
			if jOffset >= len(grid[iOffset]) {
				continue
			}
			if iOffset == i && jOffset >= jLeft && jOffset <= jRight {
				continue
			}
			log.Printf("%d %d\n", iOffset, jOffset)
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
