package main

import (
	"log"
	"os"
	"strings"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Need a filename as argument\n")
	}
	contentsBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("os.ReadFile: %s\n", err)
	}
	result := 0
	contents := string(contentsBytes)
	for _, row := range strings.Split(contents, "\n") {
		rowParts := strings.Split(row, ":")
		if len(rowParts) != 2 {
			continue
		}
		rowData := rowParts[1]
		rowDataParts := strings.Split(rowData, "|")
		winningString := rowDataParts[0]
		winning := make(map[int]bool)
		for _, nString := range strings.Split(winningString, " ") {
			if len(nString) < 1 {
				continue
			}
			n, err := strconv.Atoi(nString)
			if err != nil {
				log.Fatalf("%s\n", err)
			}
			winning[n] = true
		}
		ownString := rowDataParts[1]
		cardTotal := 0
		for _, nString := range strings.Split(ownString, " ") {
			if len(nString) < 1 {
				continue
			}
			n, err := strconv.Atoi(nString)
			if err != nil {
				log.Fatalf("%s\n", err)
			}
			_, isOk := winning[n]
			if !isOk {
				continue
			}
			if cardTotal == 0 {
				cardTotal++
			} else {
				cardTotal *= 2
			}
		}
		result += cardTotal
	}
	log.Printf("%d\n", result)
}
