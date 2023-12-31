package main

import (
	"log"
	"os"
)

const (
	FNV_OFFSET = 2166136261
	FNV_PRIME  = 16777619
	STEPS      = 100
)

type Coords struct {
	i, j int
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need arg")
	}
	contentsBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("%s", err)
	}
	grid := make([][]byte, 0)
	row := make([]byte, 0)
	for _, c := range contentsBytes {
		if c == '\n' {
			if len(row) > 0 {
				grid = append(grid, row)
				row = make([]byte, 0)
			}
		} else {
			row = append(row, c)
		}
	}

	start := &Coords{}
	for i, row := range grid {
		for j, c := range row {
			if c == 'S' {
				start.i = i
				start.j = j
			}
		}
	}

	frontier := make([]*Coords, 0)
	frontier = append(frontier, start)

	cycles := make(map[int][]int)
	chashToCycles := make(map[int][]int)
	cycleStart := make(map[int]int)
	sectorsDone := make(map[int]bool)
	gridHashes := make(map[int]int)

	for i := 0; i < STEPS; i++ {
		nextFrontier := make([]*Coords, 0)
		globalVisited := make(map[int]bool)
		for _, c := range frontier {
			for _, adj := range expand(grid, c) {
				//log.Printf("grid[%d][%d] = %c", adj.i, adj.j, grid[adj.i][adj.j])
				gridId := adj.gridhash(grid)
				_, ok := sectorsDone[gridId]
				if ok {
					continue
				}
				gridHash, ok := gridHashes[gridId]
				if !ok {
					gridHash = 0
				}
				gridHashes[gridId] = adj.hash(grid, gridHash)
				h := adj.globalhash()
				_, ok = globalVisited[h]
				if ok {
					continue
				}
				globalVisited[h] = true
				nextFrontier = append(nextFrontier, adj)
			}
		}
		frontier = nextFrontier
		for k, v := range gridHashes {
			c, ok := cycles[k]
			if !ok {
				c = make([]int, 0)
			}
			c = append(c, v)
			if i%16 == 0 {
				normalized, ok := detectCycles(c)
				if !ok {
					continue
				}
				ch := chash(normalized)
				existingNormalized, ok := chashToCycles[ch]
				if !ok {
					chashToCycles[ch] = normalized
					cycles[k] = normalized
				} else {
					chashToCycles[ch] = existingNormalized
					cycles[k] = existingNormalized
				}
				cycleStart[k] = i + 1
				sectorsDone[k] = true
			}
		}
		//log.Printf("iter: %d, frontier size: %d", i, len(frontier))
	}

	result := len(frontier)
	for k, v := range cycles {
		ok, _ := sectorsDone[k]
		if !ok {
			continue
		}
		finalStep := STEPS - cycleStart[k]
		m := finalStep % len(v)
		result += v[m]
	}
	/*
		for i := range grid {
			log.Printf("%s - %s", string(grid[i]), string(visited[i]))
		}
	*/

	log.Printf("result: %d", result)

}

func expand(grid [][]byte, c *Coords) (res []*Coords) {
	res = make([]*Coords, 0, 4)
	isz := len(grid)
	jsz := len(grid[0])
	iBase := c.i
	for iBase < 0 {
		iBase += isz
	}
	iBase %= isz
	jBase := c.j
	for jBase < 0 {
		jBase += jsz
	}
	jBase %= jsz
	i := (iBase + isz - 1) % isz
	if grid[i][jBase] != '#' {
		res = append(res, &Coords{i: c.i - 1, j: c.j})
	}
	i = (iBase + 1) % isz
	if grid[i][jBase] != '#' {
		res = append(res, &Coords{i: c.i + 1, j: c.j})
	}
	j := (jBase + jsz - 1) % jsz
	if grid[iBase][j] != '#' {
		res = append(res, &Coords{i: c.i, j: c.j - 1})
	}
	j = (jBase + 1) % jsz
	if grid[iBase][j] != '#' {
		res = append(res, &Coords{i: c.i, j: c.j + 1})
	}
	return
}

func detectCycles(c []int) (n []int, ok bool) {
	sz := len(c)
	for i := 3; i < len(c)/2; i++ {
		isDone := true
		for j := 0; j < i; j++ {
			c1 := c[sz-j-1]
			c2 := c[sz-i-j]
			if c1 != c2 {
				isDone = false
				break
			}
		}
		if isDone {
			n = make([]int, i)
			for j := sz - i; j < sz; j++ {
				n = append(n, c[j])
			}
			ok = true
			return
		}
	}
	ok = false
	return
}

func (c *Coords) globalhash() (h int) {
	h = FNV_OFFSET
	h ^= c.i
	h *= FNV_PRIME
	h ^= c.j
	h *= FNV_PRIME
	return
}

func (c *Coords) hash(grid [][]byte, h int) (res int) {
	if h == 0 {
		res = FNV_OFFSET
	} else {
		res = h
	}
	isz := len(grid)
	jsz := len(grid[0])
	i := c.i
	for i < 0 {
		i += isz
	}
	i %= isz
	j := c.j
	for j < 0 {
		j += jsz
	}
	j %= jsz
	res ^= i
	res *= FNV_PRIME
	res ^= j
	res *= FNV_PRIME
	return
}

func chash(ns []int) (res int) {
	res = FNV_OFFSET
	for _, n := range ns {
		res ^= n
		res *= FNV_PRIME
	}
	return
}

func (c *Coords) gridhash(grid [][]byte) (h int) {
	h = FNV_OFFSET
	h ^= c.i / len(grid)
	h *= FNV_PRIME
	h ^= c.j / len(grid[0])
	h *= FNV_PRIME
	return
}
