package main

import (
	"log"
	"os"
	"strings"
)

type Coords struct {
	i, j             int
	parent           *Coords
	steps            int
	numValidChildren int
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
	nodes := make([][]*Coords, len(grid))
	start := &Coords{}
	dest := &Coords{}
	for i, row := range grid {
		visited[i] = make([]int, len(row))
		nodes[i] = make([]*Coords, len(row))
		for j, c := range row {
			if i == 0 && c == '.' {
				start.i = i
				start.j = j
			} else if i == len(grid)-1 && c == '.' {
				dest.i = i
				dest.j = j
			}
			visited[i][j] = 0
		}
	}
	frontier := make([]*Coords, 0)
	frontier = append(frontier, start)
	for len(frontier) > 0 {
		c := frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]
		if c.i == dest.i && c.j == dest.j {
			if nodes[c.i][c.j] == nil || nodes[c.i][c.j].steps < c.steps {
				nodes[c.i][c.j] = c
			}
			continue
		}
		if nodes[c.i][c.j] != nil {
			if nodes[c.i][c.j].numValidChildren < 0 {
				log.Printf("num valid children lt 0: %d", nodes[c.i][c.j].numValidChildren)
			}
			if nodes[c.i][c.j].numValidChildren == 0 {
				continue
			}
		}
		if nodes[c.i][c.j] == nil || nodes[c.i][c.j].steps < c.steps {
			nodes[c.i][c.j] = c
		}
		children := c.expand(grid)
		if len(children) == 0 {
			c.backtrace()
			continue
		}
		for _, child := range children {
			frontier = append(frontier, child)
		}
	}
	result := nodes[dest.i][dest.j].steps
	/*
		for i := range visited {
			for j := range visited[i] {
				visited[i][j] += 1000
			}
			log.Printf("%v", visited[i])
		}
	*/
	gridBytes := make([][]byte, len(grid))
	for i := range grid {
		gridBytes[i] = make([]byte, len(grid[i]))
		for j := range grid[i] {
			gridBytes[i][j] = grid[i][j]
		}
	}
	p := nodes[dest.i][dest.j]
	for p != nil {
		gridBytes[p.i][p.j] = 'O'
		p = p.parent
	}
	for i := range gridBytes {
		log.Printf("%s", string(gridBytes[i]))
	}
	log.Printf("%d", result)
}

func (c *Coords) expand(grid []string) (res []*Coords) {
	res = make([]*Coords, 0)
	/*
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
	*/
	steps := c.steps + 1
	i := c.i - 1
	j := c.j
	if i >= 0 && grid[i][j] != '#' && c.searchParent(i, j) {
		res = append(res, &Coords{i: i, j: j, parent: c, steps: steps})
	}
	i = c.i + 1
	if i < len(grid) && grid[i][j] != '#' && c.searchParent(i, j) {
		res = append(res, &Coords{i: i, j: j, parent: c, steps: steps})
	}
	i = c.i
	j = c.j - 1
	if j >= 0 && grid[i][j] != '#' && c.searchParent(i, j) {
		res = append(res, &Coords{i: i, j: j, parent: c, steps: steps})
	}
	j = c.j + 1
	if j < len(grid[i]) && grid[i][j] != '#' && c.searchParent(i, j) {
		res = append(res, &Coords{i: i, j: j, parent: c, steps: steps})
	}
	c.numValidChildren = len(res)
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

func (c *Coords) backtrace() {
	p := c.parent
	for p != nil {
		if p.numValidChildren > 1 {
			return
		}
		p.numValidChildren--
		p = p.parent
	}
}
