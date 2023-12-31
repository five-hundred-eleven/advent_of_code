package main

import (
	"log"
	"os"
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

	visited := make([][]byte, len(grid))
	start := &Coords{}
	for i, row := range grid {
		visited[i] = make([]byte, len(row))
		for j, c := range row {
			if c == 'S' {
				start.i = i
				start.j = j
			}
			visited[i][j] = '.'
		}
	}

	frontier := make([]*Coords, 0)
	frontier = append(frontier, start)

	for i := 0; i < 64; i++ {
		nextFrontier := make([]*Coords, 0)
		for i := range visited {
			//log.Printf("%s - %s", string(visited[i]), string(grid[i]))
			for j := range visited[i] {
				visited[i][j] = '.'
			}
		}
		for _, c := range frontier {
			for _, adj := range expand(grid, c) {
				//log.Printf("grid[%d][%d] = %c", adj.i, adj.j, grid[adj.i][adj.j])
				if visited[adj.i][adj.j] == 'O' {
					continue
				}
				visited[adj.i][adj.j] = 'O'
				nextFrontier = append(nextFrontier, adj)
			}
		}
		frontier = nextFrontier
		//log.Printf("iter: %d, frontier size: %d", i, len(frontier))
	}

	/*
		for i := range grid {
			log.Printf("%s - %s", string(grid[i]), string(visited[i]))
		}
	*/

	log.Printf("result: %d", len(frontier))

}

func expand(grid [][]byte, c *Coords) (res []*Coords) {
	res = make([]*Coords, 0, 4)
	i := c.i - 1
	j := c.j
	if i >= 0 && grid[i][j] != '#' {
		res = append(res, &Coords{i: i, j: j})
	}
	i = c.i + 1
	if i < len(grid) && grid[i][j] != '#' {
		res = append(res, &Coords{i: i, j: j})
	}
	i = c.i
	j = c.j - 1
	if j >= 0 && grid[i][j] != '#' {
		res = append(res, &Coords{i: i, j: j})
	}
	j = c.j + 1
	if j < len(grid) && grid[i][j] != '#' {
		res = append(res, &Coords{i: i, j: j})
	}
	return
}
