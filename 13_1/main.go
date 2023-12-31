package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("f")
	}
	contentsBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("g")
	}
	contents := string(contentsBytes)
	grid := make([]string, 0)
	result := 0
	for _, line := range strings.Split(contents, "\n") {
		if len(line) == 0 {
			result += solve(grid)
			grid = make([]string, 0)
			continue
		}
		grid = append(grid, line)
	}
	log.Printf("%d", result)
}

func solve(grid []string) int {
	for i := len(grid) - 1; i >= 1; i-- {
		numDiff := 0
		hi := i
		lo := i - 1
		for {
			for j := 0; j < len(grid[0]); j++ {
				if grid[lo][j] != grid[hi][j] {
					numDiff++
				}
				if numDiff > 1 {
					break
				}
			}
			if numDiff > 1 {
				break
			}
			hi++
			if hi >= len(grid) {
				break
			}
			lo--
			if lo < 0 {
				break
			}
		}
		if numDiff == 1 {
			return i * 100
		}
	}
	for i := len(grid[0]) - 1; i >= 1; i-- {
		numDiff := 0
		hi := i
		lo := i - 1
		for {
			for j := 0; j < len(grid); j++ {
				if grid[j][lo] != grid[j][hi] {
					numDiff++
				}
				if numDiff > 1 {
					break
				}
			}
			if numDiff > 1 {
				break
			}
			hi++
			if hi >= len(grid[0]) {
				break
			}
			lo--
			if lo < 0 {
				break
			}
		}
		if numDiff == 1 {
			return i
		}
	}
	return 0
}
