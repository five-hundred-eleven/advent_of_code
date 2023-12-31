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
	contents := string(contentsBytes)
	sum := 0
	for _, line := range strings.Split(contents, "\n") {
		firstDigit := -1
		lastDigit := -1
		for _, c := range line {
			if c < '0' || c > '9' {
				continue
			}
			n := int(c - '0')
			if firstDigit == -1 {
				firstDigit = n
			}
			lastDigit = n
		}
		if firstDigit == -1 || lastDigit == -1 {
			continue
		}
		sum += firstDigit*10 + lastDigit
	}
	log.Printf("%d\n", sum)
	return
}
