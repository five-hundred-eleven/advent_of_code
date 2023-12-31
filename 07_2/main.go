package main

import (
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Hand struct {
	Cards string
	Bid   int
}
type Hands []*Hand

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Need input file as argument\n")
	}
	contentBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	content := string(contentBytes)

	hands := make(Hands, 0)

	for _, line := range strings.Split(content, "\n") {
		if len(line) == 0 {
			continue
		}
		lineParts := strings.Split(line, " ")
		bid, err := strconv.Atoi(lineParts[1])
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		h := &Hand{Cards: lineParts[0], Bid: bid}
		hands = append(hands, h)
	}

	sort.Sort(hands)
	result := 0
	for i, h := range hands {
		result += h.Bid * (i + 1)
	}

	log.Printf("%d\n", result)

}

func rank(hand string) int {

	cardToCount := make(map[rune]int)
	jokerCount := 0
	for _, c := range hand {
		if c == 'J' {
			jokerCount++
			continue
		}
		count, isOk := cardToCount[c]
		if !isOk {
			cardToCount[c] = 1
		} else {
			cardToCount[c] = count + 1
		}
	}
	if jokerCount >= 4 {
		return 7
	}
	counts := make([]int, 0)
	num4s := 0
	for _, count := range cardToCount {
		if count == 5 || (count+jokerCount) >= 5 {
			return 7
		}
		if count == 4 || (count+jokerCount) >= 4 {
			num4s++
		}
		counts = append(counts, count)
	}
	if num4s > 0 {
		return 6
	}

	num2s := 0
	num3s := 0
	for _, c := range counts {
		if c == 2 {
			num2s++
		} else if c == 3 {
			num3s++
		}
	}

	if num3s >= 1 && num2s >= 1 || num2s >= 2 && jokerCount >= 1 || num2s >= 1 && jokerCount >= 2 || jokerCount >= 3 {
		return 5
	}

	if num3s >= 1 || num2s >= 1 && jokerCount >= 1 || jokerCount >= 2 {
		return 4
	}

	if num2s >= 2 || num2s >= 1 && jokerCount >= 1 {
		return 3
	}

	if num2s >= 1 || jokerCount >= 1 {
		return 2
	}

	return 1

}

func (h Hands) Less(i, j int) bool {
	ih := h[i].Cards
	jh := h[j].Cards
	iRank := rank(ih)
	jRank := rank(jh)
	if iRank < jRank {
		return true
	}
	if iRank > jRank {
		return false
	}

	strength := map[byte]int{
		'A': 14,
		'K': 13,
		'Q': 12,
		'J': 1,
		'T': 10,
		'9': 9,
		'8': 8,
		'7': 7,
		'6': 6,
		'5': 5,
		'4': 4,
		'3': 3,
		'2': 2,
	}

	for k := range ih {
		if strength[ih[k]] < strength[jh[k]] {
			return true
		}
		if strength[ih[k]] > strength[jh[k]] {
			return false
		}
	}

	return false

}

func (h Hands) Len() int {
	return len(h)
}

func (h Hands) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
