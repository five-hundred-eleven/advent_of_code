package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Need name of input file as argument\n")
	}
	contentsBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("error reading file: %s\n", err)
	}
	contents := string(contentsBytes)

	contentsLines := strings.Split(contents, "\n")
	initParts := strings.Split(contentsLines[0], ":")
	seedsRaw := make([]int, 0)
	for _, seedString := range strings.Split(initParts[1], " ") {
		if len(seedString) == 0 {
			continue
		}
		seed, err := strconv.Atoi(seedString)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		seedsRaw = append(seedsRaw, seed)
	}
	seedsStart := make([]int, 0)
	seedsStop := make([]int, 0)
	for i := 0; i < len(seedsRaw); i += 2 {
		seedsStart = append(seedsStart, seedsRaw[i])
		seedsStop = append(seedsStop, seedsRaw[i]+seedsRaw[i+1])
	}
	seedsStartNext := make([]int, 0)
	seedsStopNext := make([]int, 0)
	i := 3
	for i < len(contentsLines) {
		line := contentsLines[i]
		if len(line) == 0 {
			for k := range seedsStart {
				if seedsStart[k] > seedsStop[k] {
					continue
				}
				seedsStartNext = append(seedsStartNext, seedsStart[k])
				seedsStopNext = append(seedsStopNext, seedsStop[k])
			}
			seedsStart = seedsStartNext
			seedsStop = seedsStopNext
			seedsStartNext = make([]int, 0)
			seedsStopNext = make([]int, 0)
			i += 2
			continue
		}
		lineParts := strings.Split(line, " ")
		destStart, _ := strconv.Atoi(lineParts[0])
		srcStart, _ := strconv.Atoi(lineParts[1])
		sz, _ := strconv.Atoi(lineParts[2])
		for k := 0; k < len(seedsStart); k++ {
			start := seedsStart[k]
			stop := seedsStop[k]
			if stop < start {
				continue
			}
			intersectStart := 0
			isStart := false
			isStop := false
			if start < srcStart {
				intersectStart = srcStart
				isStart = true
			} else {
				intersectStart = start
			}
			intersectStop := 0
			if stop < srcStart+sz-1 {
				intersectStop = stop
				isStop = true
			} else {
				intersectStop = srcStart + sz - 1
			}
			if intersectStart <= intersectStop {
				offset := destStart - srcStart
				seedsStartNext = append(seedsStartNext, intersectStart+offset)
				seedsStopNext = append(seedsStopNext, intersectStop+offset)
				if isStart && isStop {
					seedsStart = append(seedsStart, seedsStart[k])
					seedsStop = append(seedsStop, intersectStart-1)
					seedsStart = append(seedsStart, intersectStop+1)
					seedsStop = append(seedsStop, seedsStop[k])
					seedsStart[k] = 2
					seedsStop[k] = 1
				} else if isStart {
					seedsStop[k] = intersectStart - 1
				} else {
					seedsStart[k] = intersectStop + 1
				}
			}
		}
		i++
	}
	lowestSoFar := 999999999
	for _, n := range seedsStart {
		if n < lowestSoFar {
			lowestSoFar = n
		}
	}
	log.Printf("%d\n", lowestSoFar)
}
