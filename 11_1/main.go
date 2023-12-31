package main

import (
	"log"
	"os"
	"strings"
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
	for _, line := range strings.Split(contents, "\n") {
		if len(line) == 0 {
			continue
		}
		isEmpty := true
		for _, col := range line {
			if col != '.' {
				isEmpty = false
			}
		}
		grid = append(grid, []byte(line))
		if isEmpty {
			grid = append(grid, []byte(line))
		}
	}
	for x := len(grid[0]) - 1; x >= 0; x-- {
		isEmpty := true
		for y := range grid {
			if grid[y][x] != '.' {
				isEmpty = false
			}
		}
		if !isEmpty {
			continue
		}
		for y := range grid {
			grid[y] = append(grid[y], '.')
			copy(grid[y][x+1:], grid[y][x:])
			grid[y][x] = '.'
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
			if g1[0] > g2[0] {
				result += g1[0] - g2[0]
			} else {
				result += g2[0] - g1[0]
			}
			if g1[1] > g2[1] {
				result += g1[1] - g2[1]
			} else {
				result += g2[1] - g1[1]
			}
		}
	}
	log.Printf("%d\n", result)
}
