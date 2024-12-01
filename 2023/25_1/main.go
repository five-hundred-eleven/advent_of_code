package main

import (
	"log"
	"os"
	"strings"
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
	components := make(map[string][]string)
	for _, line := range strings.Split(contents, "\n") {
		if len(line) == 0 {
			continue
		}
		lineParts := strings.Split(line, ":")
		left := strings.TrimSpace(lineParts[0])
		if len(left) == 0 {
			log.Printf("line with no left: %s", line)
		}
		conns, ok := components[left]
		if !ok {
			conns = make([]string, 0)
		}
		for _, otherBase := range strings.Split(lineParts[1], " ") {
			other := strings.TrimSpace(otherBase)
			if len(other) == 0 {
				continue
			}
			conns = append(conns, other)
			otherConns, ok := components[other]
			if !ok {
				otherConns = make([]string, 0)
			}
			otherConns = append(otherConns, left)
			components[other] = otherConns
		}
		components[left] = conns
	}
	combos := make([][]string, 0)
	for k := range components {
		item := make([]string, 1)
		item[0] = k
		combos = append(combos, item)
	}
	for k := range components {
		for l := range components {
			if k == l {
				continue
			}
			item := make([]string, 2)
			item[0] = k
			item[1] = l
			combos = append(combos, item)
		}
	}
	for k := range components {
		for l := range components {
			if k == l {
				continue
			}
			for m := range components {
				if m == k {
					continue
				}
				if m == l {
					continue
				}
				item := make([]string, 3)
				item[0] = k
				item[1] = l
				item[2] = m
				combos = append(combos, item)
			}
		}
	}
	for _, item := range combos {
		log.Printf("components: %v", item)
		visited := make(map[string]bool)
		for _, x := range item {
			visited[x] = true
		}
		v := item
		for len(v) > 0 {
			nextV := make([]string, 0)
			for _, c := range v {
				_, ok := visited[c]
				if ok {
					continue
				}
				visited[c] = true
				nextV = append(nextV, c)
			}
			if len(nextV) == 3 {
				g1 := len(visited)
				g2 := len(components) - g1
				log.Printf("%d vs %d: %d", g1, g2, g1*g2)
				return
			}
			v = nextV
		}
	}
}
