package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	PMIN = 200000000000000
	PMAX = 400000000000000
	//PMIN = 7
	//PMAX = 27
)

type Vector struct {
	x, y, z float64
}

type Hailstone struct {
	line               string
	position, velocity Vector
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
	hailstones := make([]*Hailstone, 0)
	for _, line := range strings.Split(contents, "\n") {
		if len(line) == 0 {
			continue
		}
		lineParts := strings.Split(line, "@")
		posParts := strings.Split(lineParts[0], ",")
		px, err := strconv.Atoi(strings.TrimSpace(posParts[0]))
		if err != nil {
			log.Fatalf("%s", err)
		}
		py, err := strconv.Atoi(strings.TrimSpace(posParts[1]))
		if err != nil {
			log.Fatalf("%s", err)
		}
		pz, err := strconv.Atoi(strings.TrimSpace(posParts[2]))
		if err != nil {
			log.Fatalf("%s", err)
		}
		velParts := strings.Split(lineParts[1], ", ")
		vx, err := strconv.Atoi(strings.TrimSpace(velParts[0]))
		if err != nil {
			log.Fatalf("%s", err)
		}
		vx *= PMAX
		vy, err := strconv.Atoi(strings.TrimSpace(velParts[1]))
		if err != nil {
			log.Fatalf("%s", err)
		}
		vy *= PMAX
		vz, err := strconv.Atoi(strings.TrimSpace(velParts[2]))
		if err != nil {
			log.Fatalf("%s", err)
		}
		vz *= PMAX
		hailstones = append(hailstones, &Hailstone{
			line:     line,
			position: Vector{x: float64(px), y: float64(py), z: float64(pz)},
			velocity: Vector{x: float64(vx), y: float64(vy), z: float64(vz)},
		})
	}

	result := 0
	for i := 0; i < len(hailstones); i++ {
		for j := i + 1; j < len(hailstones); j++ {
			if intersection(hailstones[i], hailstones[j]) {
				result++
			}
		}
	}
	log.Printf("%d", result)
}

func intersection(h1, h2 *Hailstone) (ok bool) {
	//log.Printf("A: %s", h1.line)
	//log.Printf("B: %s", h2.line)
	hp1 := Vector{
		x: h1.position.x + h1.velocity.x,
		y: h1.position.y + h1.velocity.y,
		z: h1.position.z + h1.velocity.z,
	}
	hp2 := Vector{
		x: h2.position.x + h2.velocity.x,
		y: h2.position.y + h2.velocity.y,
		z: h2.position.z + h2.velocity.z,
	}
	den := (-h1.velocity.x * -h2.velocity.y) - (-h1.velocity.y * -h2.velocity.x)
	if den < 0.0000001 && den > -0.0000001 {
		//log.Printf("Discarding hailstones because denominator")
		ok = false
		return
	}
	c1 := h1.position.x*hp1.y - h1.position.y*hp1.x
	c2 := h2.position.x*hp2.y - h2.position.y*hp2.x
	resultX := (c1*-h2.velocity.x - (-h1.velocity.x * c2)) / den
	if resultX < PMIN || resultX > PMAX {
		//log.Printf("Got X result: %f", resultX)
		ok = false
		return
	}
	resultY := (c1*-h2.velocity.y - (-h1.velocity.y * c2)) / den
	if resultY < PMIN || resultY > PMAX {
		//log.Printf("Got Y result: %f", resultY)
		ok = false
		return
	}
	//log.Printf("Intersection: X: %f, Y: %f", resultX, resultY)
	t := (h1.position.x-h2.position.x)*(-h2.velocity.y) - (h1.position.y-h2.position.y)*(-h2.velocity.x)
	t /= den
	if t < 0 || t > 1.0 {
		//log.Printf("Got den %f, T result: %f", den, t)
		ok = false
		return
	}
	u := (h1.position.x-h2.position.x)*(-h1.velocity.y) - (h1.position.y-h2.position.y)*(-h1.velocity.x)
	u /= den
	if u < 0 || u > 1.0 {
		//log.Printf("Got den %f, U result: %f", den, u)
		ok = false
		return
	}
	ok = true
	return
}
