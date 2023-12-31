package main

import (
	"log"
	"os"
	"strings"
)

const (
	MULT = 1_000_000
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Need arg")
	}
	contentsBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("%s", err)
	}
	contents := string(contentsBytes)
	grid := make([][]byte, 0)
	emptyRows := make([]bool, 0)
	for _, line := range strings.Split(contents, "\n") {
		if len(line) == 0 {
			continue
		}
		isEmpty := true
		for _, col := range line {
			if col != '.' {
				isEmpty = false
				break
			}
		}
		grid = append(grid, []byte(line))
		if isEmpty {
			emptyRows = append(emptyRows, true)
		} else {
			emptyRows = append(emptyRows, false)
		}
	}
	emptyCols := make([]bool, 0)
	for x := range grid[0] {
		isEmpty := true
		for y := range grid {
			if grid[y][x] != '.' {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			emptyCols = append(emptyCols, true)
		} else {
			emptyCols = append(emptyCols, false)
		}
	}
	galaxies := make([][]int, 0)
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == '#' {
				galaxy := make([]int, 2)
				galaxy[0] = x
				galaxy[1] = y
				galaxies = append(galaxies, galaxy)
			}
		}
	}
	result := 0
	for i := 0; i < len(galaxies); i++ {
		g1 := galaxies[i]
		for j := i + 1; j < len(galaxies); j++ {
			g2 := galaxies[j]
			for k := min(g1[0], g2[0]); k < max(g1[0], g2[0]); k++ {
				if emptyCols[k] {
					result += MULT
				} else {
					result++
				}
			}
			for k := min(g1[1], g2[1]); k < max(g1[1], g2[1]); k++ {
				if emptyRows[k] {
					result += MULT
				} else {
					result++
				}
			}
		}
	}
	log.Printf("%d\n", result)
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
