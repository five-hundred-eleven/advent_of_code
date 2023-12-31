package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"fmt"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Need input file as arg")
	}
	contentsBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	contents := string(contentsBytes)
	contentsParts := strings.Split(contents, "\n")
	timeString := contentsParts[0]
	timeString = strings.Split(timeString, ":")[1]
	times := make([]int, 0)
	for _, timeNumString := range strings.Split(timeString, " ") {
		if len(timeNumString) == 0 {
			continue
		}
		timeNum, err := strconv.Atoi(timeNumString)
		if err != nil {
			log.Fatalf("from atoi: %s\n", err)
		}
		times = append(times, timeNum)
	}
	distanceString := contentsParts[1]
	distanceString = strings.Split(distanceString, ":")[1]
	distances := make([]int, 0)
	for _, distanceNumString := range strings.Split(distanceString, " ") {
		if len(distanceNumString) == 0 {
			continue
		}
		distanceNum, err := strconv.Atoi(distanceNumString)
		if err != nil {
			log.Fatalf("from atoi: %s\n", err)
		}
		distances = append(distances, distanceNum)
	}

	timeString = ""
	for _, t := range times {
		timeString += fmt.Sprintf("%d", t)
	}
	time, err := strconv.Atoi(timeString)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	distanceString = ""
	for _, d := range distances {
		distanceString += fmt.Sprintf("%d", d)
	}
	distance, err := strconv.Atoi(distanceString)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	searchLo := make([]int, 0)
	searchLo = append(searchLo, 0)
	searchHi := make([]int, 0)
	searchHi = append(searchHi, time)
	leftBounds := -1
	rightBounds := -1
	for j := 0; j < len(searchLo); j++ {
		lo := searchLo[j]
		hi := searchHi[j]
		if lo >= hi {
			continue
		}
		mid := (lo + hi) / 2
		loSuccess := isSuccess(lo, time, distance)
		hiSuccess := isSuccess(hi, time, distance)
		if !loSuccess && !hiSuccess {
			searchLo = append(searchLo, lo)
			searchHi = append(searchHi, mid)
			searchLo = append(searchLo, mid)
			searchHi = append(searchHi, hi)
		} else if loSuccess && !hiSuccess {
			if lo + 1 >= hi {
				rightBounds = hi
				if leftBounds != -1 {
					break
				}
			} else {
				if isSuccess(mid, time, distance) {
					searchLo = append(searchLo, mid)
					searchHi = append(searchHi, hi)
				} else {
					searchLo = append(searchLo, lo)
					searchHi = append(searchHi, mid)
				}
			}
		} else if !loSuccess && hiSuccess {
			if lo + 1 >= hi {
				leftBounds = hi
				if rightBounds != -1 {
					break
				}
			} else {
				if isSuccess(mid, time, distance) {
					searchLo = append(searchLo, lo)
					searchHi = append(searchHi, mid)
				} else {
					searchLo = append(searchLo, mid)
					searchHi = append(searchHi, hi)
				}
			}
		}
	}
	raceResult := rightBounds - leftBounds
	log.Printf("result for race: %d (%d, %d)\n", raceResult, leftBounds, rightBounds)
}

func isSuccess(buttonpress int, time int, distance int) bool {
	if buttonpress >= time {
		return false
	}
	timeLeft := time - buttonpress
	distanceCovered := buttonpress * timeLeft
	result := distanceCovered > distance
	/*
	if result {
		log.Printf("buttonpress: %d, time: %d, distance: %d, success\n", buttonpress, time, distance)
	} else {
		log.Printf("buttonpress: %d, time: %d, distance: %d, not success\n", buttonpress, time, distance)
	}
	*/
	return result
}
