package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	filename := os.Args[1]
	contentsBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	sToN := map[string]int{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"9": 9,
		"zero": 0,
		"one": 1,
		"two": 2,
		"three": 3,
		"four": 4,
		"five": 5, 
		"six": 6,
		"seven": 7,
		"eight": 8,
		"nine": 9,
	}
	contents := string(contentsBytes)
	sum := 0
	for _, line := range strings.Split(contents, "\n") {
		firstDigit := -1
		lastDigit := -1
		for i := 0; i < len(line); i++ {
			remaining := len(line) - i
			for s, n := range sToN {
				if len(s) > remaining {
					continue
				}
				sLen := len(s)
				segment := line[i:i+sLen]
				if segment == s {
					if firstDigit == -1 {
						firstDigit = n
					}
					lastDigit = n
				}
			}
		}
		if firstDigit == -1 || lastDigit == -1 {
			continue
		}
		sum += firstDigit*10 + lastDigit
	}
	log.Printf("%d\n", sum)
	return
}
