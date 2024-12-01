package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Need a filename as argument\n")
	}
	contentsBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("os.ReadFile: %s\n", err)
	}
	contents := string(contentsBytes)
	cardToResult := make(map[int]int)
	rowIds := make([]int, 0)
	for _, row := range strings.Split(contents, "\n") {
		rowParts := strings.Split(row, ":")
		if len(rowParts) != 2 {
			continue
		}
		rowIdStr := rowParts[0]
		rowIdParts := strings.Split(rowIdStr, " ")
		rowId, err := strconv.Atoi(rowIdParts[len(rowIdParts)-1])
		if err != nil {
			log.Fatalf("%s\n", err)
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
			cardTotal++
		}
		rowIds = append(rowIds, rowId)
		cardToResult[rowId] = cardTotal
	}
	rowIndex := 0
	for rowIndex < len(rowIds) {
		currRowId := rowIds[rowIndex]
		cardResult := cardToResult[currRowId]
		for i := 1; i <= cardResult; i++ {
			rowIds = append(rowIds, currRowId+i)
		}
		rowIndex++
	}
	log.Printf("%d\n", len(rowIds))
}
