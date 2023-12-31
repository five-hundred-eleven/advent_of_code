package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

const (
	limit = 1_000_000_000
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Need an arg")
	}
	contentsBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("%s", err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(contentsBytes))
	grid := make([][]byte, 0)
	origSanity := 0
	for scanner.Scan() {
		row := []byte(scanner.Bytes())
		if len(row) == 0 {
			continue
		}
		gridRow := make([]byte, 0)
		for _, c := range row {
			if c == 'O' {
				origSanity++
			}
			gridRow = append(gridRow, c)
		}
		grid = append(grid, gridRow)
	}
	results := make([]int, 0, 100)
	finalResult := 0
	isDone := false
	for i := 0; i < limit; i++ {
		load, sanity := cycle(grid)
		if sanity != origSanity {
			log.Fatalf("sanity changed: %d -> %d", origSanity, sanity)
		}
		results = append(results, load)
		if len(results) > 100 {
			results = results[1:]
			isDone, finalResult = checkResults(i, results)
			if isDone {
				break
			}
		}
	}
	//log.Printf("%v", results)
	printGrid(grid)
	log.Printf("%d", finalResult)
}

func checkResults(it int, results []int) (done bool, finalResult int) {
	resultsLen := len(results)
	foundSomething := true
	cycleLen := 0
	for i := 5; i < resultsLen/2; i++ {
		foundSomething = true
		for j := 0; j < i; j++ {
			if results[resultsLen-j-1] != results[resultsLen-i-j-1] {
				foundSomething = false
				break
			}
		}
		if foundSomething {
			cycleLen = i
			break
		}
	}
	if !foundSomething {
		done = false
		finalResult = 0
		return
	}
	log.Printf("cycle: %v", results[resultsLen-(cycleLen*2):])
	index := resultsLen - 1
	for i := it + 1; i < limit; i++ {
		index++
		if index >= resultsLen {
			index = resultsLen - cycleLen
		}
	}
	done = true
	finalResult = results[index]
	return
}

func cycle(grid [][]byte) (load int, sanity int) {
	north(grid)
	west(grid)
	south(grid)
	east(grid)
	vertLen := len(grid)
	load = 0
	sanity = 0
	for i, row := range grid {
		for _, c := range row {
			if c == 'O' {
				load += vertLen - i
				sanity++
			}
		}
	}
	return
}

func north(grid [][]byte) {
	vertLen := len(grid)
	for j := 0; j < len(grid[0]); j++ {
		jp := j
		rounds := make([]int, 0, vertLen)
		blanks := make([]int, 0, vertLen)
		numRounds := 0
		numBlanks := 0
		for _, row := range grid {
			if row[j] == '#' {
				//log.Printf("row[%d] is square: %c", j, row[j])
				rounds = append(rounds, numRounds)
				blanks = append(blanks, numBlanks)
				numRounds = 0
				numBlanks = 0
			} else if row[j] == 'O' {
				//log.Printf("row[%d] is round: %c", j, row[j])
				numRounds++
			} else if row[j] == '.' {
				//log.Printf("row[%d] is blank: %c", j, row[j])
				numBlanks++
			} else {
				log.Fatalf("unrecognized symbol: %c", row[j])
			}
		}
		if numRounds > 0 || numBlanks > 0 {
			rounds = append(rounds, numRounds)
			blanks = append(blanks, numBlanks)
		}
		//log.Printf("%v", rounds)
		//log.Printf("%v", blanks)
		rowIndex := 0
		for k := 0; k < len(rounds); k++ {
			for l := 0; l < rounds[k]; l++ {
				/*
					if rowIndex > 0 && grid[rowIndex-1][j] == '.' {
						log.Fatalf("bad byte: %d %d", rowIndex, j)
					}
				*/
				grid[rowIndex][j] = 'O'
				/*
					if !check(grid, rowIndex, j) {
						log.Fatalf("O (%d, %d) fatal", rowIndex, j)
					}
				*/
				rowIndex++
			}
			for l := 0; l < blanks[k]; l++ {
				/*
					beforeChange := make([]byte, 0, vertLen)
					for m := 0; m < vertLen; m++ {
						beforeChange = append(beforeChange, grid[m][j])
					}
				*/
				grid[rowIndex][j] = '.'
				/*
					if !check(grid, rowIndex, j) {
						afterChange := make([]byte, 0, vertLen)
						for m := 0; m < vertLen; m++ {
							afterChange = append(afterChange, grid[m][j])
						}
						bc := string(beforeChange)
						ac := string(afterChange)
						log.Printf("before: %s", bc)
						log.Printf("after : %s", ac)
						log.Fatalf(". (%d, %d) fatal", rowIndex, j)
					}
				*/
				rowIndex++
			}
			if rowIndex < vertLen {
				grid[rowIndex][j] = '#'
				rowIndex++
			}
		}
		if rowIndex != vertLen {
			log.Fatalf("Bad row index: %d != %d", rowIndex, vertLen)
		}
		if jp != j {
			log.Fatalf("j shifted: %d -> %d", jp, j)
		}
	}
	/*
		fmt.Printf("after:\n")
		printGrid(grid)
	*/
}

func south(grid [][]byte) {
	vertLen := len(grid)
	for j := 0; j < len(grid[0]); j++ {
		rounds := make([]int, 0, vertLen)
		blanks := make([]int, 0, vertLen)
		numRounds := 0
		numBlanks := 0
		for i := len(grid) - 1; i >= 0; i-- {
			row := grid[i]
			if row[j] == '#' {
				rounds = append(rounds, numRounds)
				blanks = append(blanks, numBlanks)
				numRounds = 0
				numBlanks = 0
			} else if row[j] == 'O' {
				numRounds++
			} else if row[j] == '.' {
				numBlanks++
			} else {
				log.Fatalf("unrecognized symbol: %c", row[j])
			}
		}
		if numRounds > 0 || numBlanks > 0 {
			rounds = append(rounds, numRounds)
			blanks = append(blanks, numBlanks)
		}
		rowIndex := vertLen - 1
		for k := 0; k < len(rounds); k++ {
			for l := 0; l < rounds[k]; l++ {
				grid[rowIndex][j] = 'O'
				rowIndex--
			}
			for l := 0; l < blanks[k]; l++ {
				grid[rowIndex][j] = '.'
				rowIndex--
			}
			if rowIndex >= 0 {
				grid[rowIndex][j] = '#'
				rowIndex--
			}
		}
	}
}

func west(grid [][]byte) {
	horiLen := len(grid[0])
	vertLen := len(grid)
	for i := 0; i < vertLen; i++ {
		rounds := make([]int, 0, vertLen)
		blanks := make([]int, 0, vertLen)
		numRounds := 0
		numBlanks := 0
		for j := 0; j < horiLen; j++ {
			if grid[i][j] == '#' {
				rounds = append(rounds, numRounds)
				blanks = append(blanks, numBlanks)
				numRounds = 0
				numBlanks = 0
			} else if grid[i][j] == 'O' {
				numRounds++
			} else if grid[i][j] == '.' {
				numBlanks++
			} else {
				log.Fatalf("unrecognized symbol: %c", grid[i][j])
			}
		}
		if numRounds > 0 || numBlanks > 0 {
			rounds = append(rounds, numRounds)
			blanks = append(blanks, numBlanks)
		}
		colIndex := 0
		for k := 0; k < len(rounds); k++ {
			for m := 0; m < rounds[k]; m++ {
				grid[i][colIndex] = 'O'
				colIndex++
			}
			for m := 0; m < blanks[k]; m++ {
				grid[i][colIndex] = '.'
				colIndex++
			}
			if colIndex < horiLen {
				grid[i][colIndex] = '#'
				colIndex++
			}
		}
	}
}

func east(grid [][]byte) {
	horiLen := len(grid[0])
	vertLen := len(grid)
	for i := 0; i < vertLen; i++ {
		rounds := make([]int, 0, vertLen)
		blanks := make([]int, 0, vertLen)
		numRounds := 0
		numBlanks := 0
		for j := horiLen - 1; j >= 0; j-- {
			if grid[i][j] == '#' {
				rounds = append(rounds, numRounds)
				blanks = append(blanks, numBlanks)
				numRounds = 0
				numBlanks = 0
			} else if grid[i][j] == 'O' {
				numRounds++
			} else if grid[i][j] == '.' {
				numBlanks++
			} else {
				log.Fatalf("unrecognized symbol: %c", grid[i][j])
			}
		}
		if numRounds > 0 || numBlanks > 0 {
			rounds = append(rounds, numRounds)
			blanks = append(blanks, numBlanks)
		}
		colIndex := horiLen - 1
		for k := 0; k < len(rounds); k++ {
			for m := 0; m < rounds[k]; m++ {
				grid[i][colIndex] = 'O'
				colIndex--
			}
			for m := 0; m < blanks[k]; m++ {
				grid[i][colIndex] = '.'
				colIndex--
			}
			if colIndex >= 0 {
				grid[i][colIndex] = '#'
				colIndex--
			}
		}
	}
}

func printGrid(grid [][]byte) {
	fmt.Printf("    ")
	for i := 0; i < len(grid[0]); i++ {
		printD(i, 1, 1)
	}
	fmt.Printf("\n")
	for i, row := range grid {
		printD(i, 0, 1)
		for _, c := range row {
			fmt.Printf("   %c ", c)
		}
		fmt.Printf("\n")
	}
}

func printD(d int, lpad int, rpad int) {
	for i := 0; i < lpad; i++ {
		fmt.Printf(" ")
	}
	if d < 10 {
		fmt.Printf("  %d", d)
	} else if d < 100 {
		fmt.Printf(" %d", d)
	} else {
		fmt.Printf("%d", d)
	}
	for i := 0; i < rpad; i++ {
		fmt.Printf(" ")
	}
}

func check(grid [][]byte, rowLimit int, colLimit int) bool {
	for j := 0; j <= colLimit; j++ {
		for i := 0; i <= rowLimit; i++ {
			if i > 0 && grid[i-1][j] == '.' && grid[i][j] == 'O' {
				for k := 1; k <= i; k++ {
					log.Printf("(%d, %d): %c (prev %c)", k, j, grid[k][j], grid[k-1][j])
				}
				log.Printf("bad byte at (%d %d) row %d col %d", i, j, rowLimit, colLimit)
				return false
			} else if i > 0 {
				log.Printf("(%d, %d) ok: %c (prev %c)", i, j, grid[i][j], grid[i-1][j])
			}
		}
	}
	return true
}
