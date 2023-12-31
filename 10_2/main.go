package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var connectors = map[byte][]int{
	'|': {-1, 0, 1, 0},
	'-': {0, -1, 0, 1},
	'L': {-1, 0, 0, 1},
	'J': {-1, 0, 0, -1},
	'7': {1, 0, 0, -1},
	'F': {1, 0, 0, 1},
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Need arg")
	}
	contentsBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("%s", err)
	}
	contents := string(contentsBytes)
	grid := strings.Split(contents, "\n")
	sy := -1
	sx := -1
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == 'S' {
				sy = y
				sx = x
				break
			}
		}
		if sy != -1 {
			break
		}
	}
	visited := make([][]int, 0)
	for y := range grid {
		newRow := make([]int, 0)
		for x := range grid[y] {
			if x == sx && y == sy {
				newRow = append(newRow, 0)
			} else {
				newRow = append(newRow, -1)
			}
		}
		visited = append(visited, newRow)
	}
	q := findConns(grid, sy, sx)
	for _, curr := range q {
		visited[curr[0]][curr[1]] = 1
	}
	i := 0
	turns := make([]int, len(q))
	for t := range turns {
		turns[t] = 1
	}
	dirs := make([]int, len(q))
	dirs[0] = -1
	dirs[1] = 1
	inside := make([][]int, 0)
	for y := range grid {
		newRow := make([]int, 0)
		for range grid[y] {
			newRow = append(newRow, 0)
		}
		log.Printf("%d", len(newRow))
		inside = append(inside, newRow)
	}
	inside[sy][sx] = -1
	for _, item := range q {
		inside[item[0]][item[1]] = -1
	}
	insideCounter1 := 1
	for i < len(q) {
		curr := q[i]
		turn := turns[i] + 1
		doneTurn, next := getConn(grid, curr[0], curr[1], visited)
		if doneTurn != -1 {
			break
		}
		if next != nil {
			if visited[next[0]][next[1]] == -1 {
				q = append(q, next)
				turns = append(turns, turn)
				visited[next[0]][next[1]] = turn
				inside[next[0]][next[1]] = -1
				if dirs[i] == 1 {
					dirs = append(dirs, 1)
					adj := getRight(next, curr)
					if adj[0] >= 0 && adj[0] < len(inside) && adj[1] >= 0 && adj[1] < len(inside[adj[0]]) && inside[adj[0]][adj[1]] == 0 {
						inside[adj[0]][adj[1]] = insideCounter1
					}
					/*
						adj = getLeft(next, curr)
						if adj[0] >= 0 && adj[0] < len(inside) && adj[1] >= 0 && adj[1] < len(inside[adj[0]]) && inside[adj[0]][adj[1]] == 0 {
							inside[adj[0]][adj[1]] = insideCounter2
						}
					*/
					adj = getLeft(curr, next)
					if adj[0] >= 0 && adj[0] < len(inside) && adj[1] >= 0 && adj[1] < len(inside[adj[0]]) && inside[adj[0]][adj[1]] == 0 {
						inside[adj[0]][adj[1]] = insideCounter1
					}
					/*
						adj = getRight(curr, next)
						if adj[0] >= 0 && adj[0] < len(inside) && adj[1] >= 0 && adj[1] < len(inside[adj[0]]) && inside[adj[0]][adj[1]] == 0 {
							inside[adj[0]][adj[1]] = insideCounter2
						}
					*/
				} else {
					dirs = append(dirs, -1)
					adj := getLeft(curr, next)
					/*
						if adj[0] >= 0 && adj[0] < len(inside) && adj[1] >= 0 && adj[1] < len(inside[adj[0]]) && inside[adj[0]][adj[1]] == 0 {
							inside[adj[0]][adj[1]] = insideCounter2
						}
					*/
					adj = getRight(curr, next)
					if adj[0] >= 0 && adj[0] < len(inside) && adj[1] >= 0 && adj[1] < len(inside[adj[0]]) && inside[adj[0]][adj[1]] == 0 {
						inside[adj[0]][adj[1]] = insideCounter1
					}
					/*
						adj = getRight(next, curr)
						if adj[0] >= 0 && adj[0] < len(inside) && adj[1] >= 0 && adj[1] < len(inside[adj[0]]) && inside[adj[0]][adj[1]] == 0 {
							inside[adj[0]][adj[1]] = insideCounter2
						}
					*/
					adj = getLeft(next, curr)
					if adj[0] >= 0 && adj[0] < len(inside) && adj[1] >= 0 && adj[1] < len(inside[adj[0]]) && inside[adj[0]][adj[1]] == 0 {
						inside[adj[0]][adj[1]] = insideCounter1
					}
				}
			}
		}
		i++
	}

	result := 0
	for y := range inside {
		for x := range inside[y] {
			if inside[y][x] == -1 {
				fmt.Printf("p")
			} else if rayTrace(inside, grid, y, x) {
				fmt.Printf("i")
				result++
			} else {
				fmt.Printf("u")
			}
		}
		fmt.Printf("\n")
	}
	log.Printf("%d", result)
	return
}

func findConns(grid []string, sy int, sx int) (res [][]int) {
	res = make([][]int, 0)
	for y := sy - 1; y <= sy+1; y++ {
		if y < 0 || y >= len(grid) {
			continue
		}
		for x := sx - 1; x <= sx+1; x++ {
			if x < 0 || x >= len(grid[y]) || (x == sx && y == sy) {
				continue
			}
			conns, ok := connectors[grid[y][x]]
			if !ok {
				continue
			}
			if y+conns[0] == sy && x+conns[1] == sx {
				res = append(res, []int{y, x})
			} else if y+conns[2] == sy && x+conns[3] == sx {
				res = append(res, []int{y, x})
			}
		}
	}
	return
}

func isConnected(grid []string, y1, x1, y2, x2 int) bool {
	c1, ok := connectors[grid[y1][x1]]
	if !ok {
		return false
	}
	if !((y1+c1[0] == y2 && x1+c1[1] == x2) || (y1+c1[2] == y2 && x1+c1[3] == x2)) {
		return false
	}
	c2, ok := connectors[grid[y2][x2]]
	if !ok {
		return false
	}
	if !((y2+c2[0] == y1 && x2+c2[1] == x1) || (y2+c2[2] == y1 && x2+c2[3] == x1)) {
		return false
	}
	return true
}

func getConn(grid []string, y1, x1 int, visited [][]int) (int, []int) {
	c, ok := connectors[grid[y1][x1]]
	if !ok {
		return 0, nil
	}
	y2 := y1 + c[0]
	x2 := x1 + c[1]
	var res1, res2 []int
	turn1 := -1
	turn2 := -1
	is1 := false
	if isConnected(grid, y1, x1, y2, x2) {
		res1 = make([]int, 2)
		res1[0] = y2
		res1[1] = x2
		turn1 = visited[y2][x2]
		is1 = true
	}
	y2 = y1 + c[2]
	x2 = x1 + c[3]
	is2 := false
	if isConnected(grid, y1, x1, y2, x2) {
		res2 = make([]int, 2)
		res2[0] = y2
		res2[1] = x2
		turn2 = visited[y2][x2]
		is2 = true
	}
	if is1 && is2 && turn1 != -1 && turn2 != -1 {
		if turn1 > turn2 {
			return turn1, res1
		} else {
			return turn2, res2
		}
	}
	if is1 && turn1 == -1 {
		return -1, res1
	}
	if is2 && turn2 == -1 {
		return -1, res2
	}
	return -1, nil
}

func getLeft(curr []int, next []int) []int {
	res := make([]int, 2)
	if curr[0] < next[0] {
		res[0] = curr[0]
		res[1] = curr[1] - 1
	}
	if curr[0] > next[0] {
		res[0] = curr[0]
		res[1] = curr[1] + 1
	}
	if curr[1] < next[1] {
		res[0] = curr[0] + 1
		res[1] = curr[1]
	}
	if curr[1] > next[1] {
		res[0] = curr[0] - 1
		res[1] = curr[1]
	}
	return res
}

func getRight(curr []int, next []int) []int {
	res := make([]int, 2)
	if curr[0] < next[0] {
		res[0] = curr[0]
		res[1] = curr[1] + 1
	}
	if curr[0] > next[0] {
		res[0] = curr[0]
		res[1] = curr[1] - 1
	}
	if curr[1] < next[1] {
		res[0] = curr[0] - 1
		res[1] = curr[1]
	}
	if curr[1] > next[1] {
		res[0] = curr[0] + 1
		res[1] = curr[1]
	}
	return res
}

func countInside(inside [][]int, side int) int {
	//log.Printf("searching side: %d", side)
	frontier := make([][]int, 0)
	for y := range inside {
		for x := range inside[y] {
			if inside[y][x] == side {
				//log.Printf("found side: %d %d", y, x)
				curr := make([]int, 2)
				curr[0] = y
				curr[1] = x
				frontier = append(frontier, curr)
			}
		}
	}
	i := 0
	for i < len(frontier) {
		curr := frontier[i]
		for py := curr[0] - 1; py <= curr[0]+1; py++ {
			if py < 0 || py >= len(inside) {
				return -1
			}
			for px := curr[1] - 1; px <= curr[1]+1; px++ {
				if px < 0 || px >= len(inside[py]) {
					return -1
				}
				if py == curr[0] && px == curr[1] {
					continue
				}
				if inside[py][px] == 0 {
					inside[py][px] = side
					next := make([]int, 2)
					next[0] = py
					next[1] = px
					frontier = append(frontier, next)
				}
			}
		}
		i++
	}
	return len(frontier)
}

func rayTrace(inside [][]int, grid []string, sy, sx int) bool {
	if inside[sy][sx] == -1 {
		return false
	}
	counter := 0
	for y := sy + 1; y < len(inside[sx]); y++ {
		if inside[y][sx] == -1 && grid[y][sx] != 'S' && grid[y][sx] != '|' && grid[y][sx] != 'L' && grid[y][sx] != 'F' {
			counter++
		}
	}
	if counter%2 == 0 {
		return false
	}
	return true
}
