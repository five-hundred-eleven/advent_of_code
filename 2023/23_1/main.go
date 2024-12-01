package main

import (
	"log"
	"os"
	"strings"
)

type Coords struct {
	i, j   int
	parent *Coords
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
	grid := make([]string, 0)
	for _, row := range strings.Split(contents, "\n") {
		if len(row) == 0 {
			continue
		}
		grid = append(grid, row)
	}
	visited := make([][]int, len(grid))
	start := &Coords{}
	dest := &Coords{}
	for i, row := range grid {
		visited[i] = make([]int, len(row))
		for j, c := range row {
			if i == 0 && c == '.' {
				start.i = i
				start.j = j
			} else if i == len(grid)-1 && c == '.' {
				dest.i = i
				dest.j = j
			}
		}
	}
	frontier := make([]*Coords, 0)
	frontier = append(frontier, start)
	i := 1
	for len(frontier) > 0 {
		nextFrontier := make([]*Coords, 0)
		for _, c := range frontier {
			for _, e := range c.expand(grid) {
				if visited[e.i][e.j] >= i {
					continue
				}
				visited[e.i][e.j] = i
				nextFrontier = append(nextFrontier, e)
			}
		}
		frontier = nextFrontier
		i++
	}
	log.Printf("%d", visited[dest.i][dest.j])
}

func (c *Coords) expand(grid []string) (res []*Coords) {
	res = make([]*Coords, 0)
	if grid[c.i][c.j] == '>' {
		j := c.j + 1
		if j < len(grid[c.i]) && grid[c.i][j] != '#' {
			if c.searchParent(c.i, j) {
				res = append(res, &Coords{i: c.i, j: j, parent: c})
			}
		}
		return
	}
	if grid[c.i][c.j] == '<' {
		j := c.j - 1
		if j >= 0 && grid[c.i][j] != '#' {
			if c.searchParent(c.i, j) {
				res = append(res, &Coords{i: c.i, j: j, parent: c})
			}
		}
		return
	}
	if grid[c.i][c.j] == 'v' {
		i := c.i + 1
		if i < len(grid) && grid[i][c.j] != '#' {
			if c.searchParent(i, c.j) {
				res = append(res, &Coords{i: i, j: c.j, parent: c})
			}
		}
		return
	}
	i := c.i - 1
	j := c.j
	if i >= 0 && grid[i][j] != '#' && c.searchParent(i, j) {
		res = append(res, &Coords{i: i, j: j, parent: c})
	}
	i = c.i + 1
	if i < len(grid) && grid[i][j] != '#' && c.searchParent(i, j) {
		res = append(res, &Coords{i: i, j: j, parent: c})
	}
	i = c.i
	j = c.j - 1
	if j >= 0 && grid[i][j] != '#' && c.searchParent(i, j) {
		res = append(res, &Coords{i: i, j: j, parent: c})
	}
	j = c.j + 1
	if j < len(grid[i]) && grid[i][j] != '#' && c.searchParent(i, j) {
		res = append(res, &Coords{i: i, j: j, parent: c})
	}
	return
}

func (c *Coords) searchParent(i, j int) (ok bool) {
	p := c.parent
	for p != nil {
		if p.i == i && p.j == j {
			return false
		}
		p = p.parent
	}
	return true
}
