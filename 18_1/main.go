package main

import (
	"log"
	"os"
	"strings"
)

const (
	INC = 1024
)

type Coords struct {
	x, y int
}

type Line struct {
	start, stop Coords
}

type LineArr struct {
	arr []*Line
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need arg")
	}
	contentsBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("%s", err)
	}
	contents := string(contentsBytes)
	directions := make([]byte, 0)
	distances := make([]int, 0)
	for _, line := range strings.Split(contents, "\n") {
		if len(line) == 0 {
			continue
		}
		lineParts := strings.Split(line, " ")
		hex := lineParts[2]
		/*
			directions = append(directions, lineParts[0][0])
			dist, _ := strconv.Atoi(lineParts[1])
			distances = append(distances, dist)
		*/
		dir := byte(' ')
		if hex[7] == '0' {
			dir = 'R'
		} else if hex[7] == '1' {
			dir = 'D'
		} else if hex[7] == '2' {
			dir = 'L'
		} else if hex[7] == '3' {
			dir = 'U'
		} else {
			log.Fatalf("bad dir: %s, %c", hex, hex[7])
		}
		directions = append(directions, dir)
		dist := 0
		mult := 1
		for i := 6; i > 1; i-- {
			val := 0
			if hex[i] >= '0' && hex[i] <= '9' {
				val = int(hex[i] - '0')
			} else if hex[i] >= 'a' && hex[i] <= 'f' {
				val = int(hex[i] - 'a' + 10)
			} else {
				log.Fatalf("bad hex: %s, %c", hex, hex[i])
			}
			dist += val * mult
			mult *= 16
		}
		distances = append(distances, dist)
		log.Printf("hex: %s, dist: %d, dir: %c", hex, dist, directions[len(directions)-1])
	}

	verts := make([]*Coords, 0)
	perimeter := 0
	curr := Coords{x: 0, y: 0}
	verts = append(verts, &Coords{x: 0, y: 0})
	for i := 0; i < len(directions); i++ {
		dir := directions[i]
		dist := distances[i]
		perimeter += dist
		if dir == 'U' {
			curr.y -= dist
			point := &Coords{x: curr.x, y: curr.y}
			verts = append(verts, point)
		} else if dir == 'D' {
			curr.y += dist
			point := &Coords{x: curr.x, y: curr.y}
			verts = append(verts, point)
		} else if dir == 'L' {
			curr.x -= dist
			point := &Coords{x: curr.x, y: curr.y}
			verts = append(verts, point)
		} else if dir == 'R' {
			curr.x += dist
			point := &Coords{x: curr.x, y: curr.y}
			verts = append(verts, point)
		} else {
			log.Fatal("bad dir")
		}
	}
	result := shoelace(verts) + perimeter/2 + 1
	log.Printf("%d", result)
}

func shoelace(verts []*Coords) (result int) {
	for i := 0; i < len(verts)-1; i++ {
		result += verts[i].y*verts[i+1].x - verts[i+1].y*verts[i].x
	}
	if result < 0 {
		result = -result
	}
	result /= 2
	return
}
