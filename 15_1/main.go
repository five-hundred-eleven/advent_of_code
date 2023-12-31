package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need arg")
	}
	contentsBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("%s", err)
	}
	contents := string(contentsBytes)
	steps := make([]string, 0)
	curr := make([]byte, 0)
	for _, c := range contents {
		if c == ',' {
			steps = append(steps, string(curr))
			curr = make([]byte, 0)
		} else if c == '\n' {
			continue
		} else {
			curr = append(curr, byte(c))
		}
	}
	steps = append(steps, string(curr))
	result := 0
	for _, s := range steps {
		result += HASH(s)
	}
	log.Printf("%d", result)
}

func HASH(s string) int {
	curr := 0
	for _, c := range s {
		curr += int(c)
		curr *= 17
		curr %= 256
	}
	log.Printf("%s: %d", s, curr)
	return curr
}
