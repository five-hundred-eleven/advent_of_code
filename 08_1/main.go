package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Need arg")
	}
	contentBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	content := string(contentBytes)
	contentLines := strings.Split(content, "\n")
	directions := contentLines[0]
	nodes := make(map[string][]string)
	startNodes := make([]string, 0)
	for _, line := range contentLines[1:] {
		if len(line) == 0 {
			continue
		}
		lineParts := strings.Split(line, "=")
		if len(lineParts) != 2 {
			log.Printf("warning: line %s\n", line)
			continue
		}
		left := strings.TrimSpace(lineParts[0])
		if left[len(left)-1] == 'A' {
			startNodes = append(startNodes, left)
		}
		right := strings.TrimSpace(lineParts[1])
		right = strings.TrimLeft(right, "(")
		right = strings.TrimRight(right, ")")
		rightParts := strings.Split(right, ",")
		nodes[left] = make([]string, 0)
		for _, part := range rightParts {
			nodes[left] = append(nodes[left], strings.TrimSpace(part))
		}
	}

	steps := 0
	currNodes := startNodes
	for {
		for _, dir := range directions {
			steps++
			numEnded := 0
			for i, curr := range currNodes {
				nextNodes, ok := nodes[curr]
				if !ok {
					log.Fatalf("%s\n", curr)
				}
				if dir == 'L' {
					currNodes[i] = nextNodes[0]
				} else if dir == 'R' {
					currNodes[i] = nextNodes[1]
				} else {
					log.Fatalf("bad dir: %c\n", dir)
				}
				if currNodes[i][len(currNodes[i])-1] == 'Z' {
					numEnded++
				}
			}
			if numEnded == len(currNodes) {
				log.Printf("%d\n", steps)
				return
			}
		}
	}
}
