package main

import (
	"log"
	"os"
	"strconv"
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
	hm := make([]map[string]int, 256)
	lru := make([][]string, 256)
	for i := 0; i < len(hm); i++ {
		hm[i] = make(map[string]int)
		lru[i] = make([]string, 0)
	}
	for _, s := range steps {
		doStep(hm, lru, s)
	}
	result := 0
	for i := 0; i < len(hm); i++ {
		for j := 0; j < len(lru[i]); j++ {
			label := lru[i][j]
			focalLength, _ := hm[i][label]
			prod := (i + 1) * (j + 1) * focalLength
			log.Printf("%s -> %d: %d * %d * %d = %d", label, focalLength, i+1, j+1, focalLength, prod)
			result += prod
		}
	}
	log.Printf("%d", result)
}

func doStep(hm []map[string]int, lru [][]string, s string) {
	slen := len(s)
	if s[slen-1] == '-' {
		label := s[:slen-1]
		h := HASH(label)
		box := hm[h]
		_, ok := box[label]
		if !ok {
			return
		}
		delete(box, label)
		for i := 0; i < len(lru[h]); i++ {
			if lru[h][i] == label {
				copy(lru[h][i:], lru[h][i+1:])
				lru[h] = lru[h][:len(lru[h])-1]
			}
		}
	} else {
		sParts := strings.Split(s, "=")
		label := sParts[0]
		value, _ := strconv.Atoi(sParts[1])
		h := HASH(label)
		box := hm[h]
		_, ok := box[label]
		box[label] = value
		if !ok {
			lru[h] = append(lru[h], label)
		}
	}
}

func HASH(s string) int {
	curr := 0
	for _, c := range s {
		curr += int(c)
		curr *= 17
		curr %= 256
	}
	return curr
}
