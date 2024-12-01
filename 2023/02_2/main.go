package main

import (
	"os"
	"log"
	"strings"
	"strconv"
)

func main() {
	filename := os.Args[1]
	contentsBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	contents := string(contentsBytes)
	result := int64(0)
	for _, line := range strings.Split(contents, "\n") {
		if len(line) == 0 {
			continue
		}
		lineParts := strings.Split(line, ":")
		/*
		gameIdParts := strings.Split(lineParts[0], " ")
		gameId, err := strconv.Atoi(gameIdParts[1])
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		*/
		gameMaxCubes := map[string]int{
			"red": 0,
			"green": 0,
			"blue": 0,
		}
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
				if count > gameMaxCubes[color] {
					gameMaxCubes[color] = count
				}
			}
		}
		result += int64(gameMaxCubes["red"] * gameMaxCubes["green"] * gameMaxCubes["blue"])
	}
	log.Printf("%d\n", result)
}
