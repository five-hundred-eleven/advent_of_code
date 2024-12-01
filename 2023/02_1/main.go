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
	contents := string(contentsBytes)
	result := 0
	maxCubes := map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}
	for _, line := range strings.Split(contents, "\n") {
		if len(line) == 0 {
			continue
		}
		lineParts := strings.Split(line, ":")
		gameIdParts := strings.Split(lineParts[0], " ")
		gameId, err := strconv.Atoi(gameIdParts[1])
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		isPossible := true
		for _, round := range strings.Split(lineParts[1], ";") {
			for _, item := range strings.Split(round, ",") {
				itemParts := strings.Split(item, " ")
				if len(itemParts) < 3 {
					continue
				}
				color := strings.TrimSpace(itemParts[2])
				countString := strings.TrimSpace(itemParts[1])
				count, err := strconv.Atoi(countString)
				if err != nil {
					log.Fatalf("%s\n", err)
				}
				if count > maxCubes[color] {
					isPossible = false
					break
				}
			}
			if !isPossible {
				break
			}
		}
		if isPossible {
			result += gameId
		}
	}
	log.Printf("%d\n", result)
}
