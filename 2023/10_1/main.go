package main

import (
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
	log.Printf("%d %d", sy, sx)
	q := findConns(grid, sy, sx)
	for _, curr := range q {
		visited[curr[0]][curr[1]] = 1
	}
	log.Printf("q: %d", len(q))
	i := 0
	turns := make([]int, len(q))
	for t := range turns {
		turns[t] = 1
	}
	for i < len(q) {
		curr := q[i]
		turn := turns[i] + 1
		doneTurn, next := getConn(grid, curr[0], curr[1], visited)
		if doneTurn != -1 {
			log.Printf("%d", doneTurn)
			break
		}
		if next != nil && visited[next[0]][next[1]] != -1 {
			log.Printf("already visited: %d %d, %s", next[0], next[1], grid[next[0]])
		}
		if next != nil {
			log.Printf("(%d, %d) -> (%d, %d) turn %d", curr[0], curr[1], next[0], next[1], turn)
			if visited[next[0]][next[1]] == -1 {
				q = append(q, next)
				turns = append(turns, turn)
				visited[next[0]][next[1]] = turn
			}
		}
		i++
	}
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
			log.Printf("(%d, %d): %d %d %d %d", y, x, conns[0], conns[1], conns[2], conns[3])
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
	log.Printf("getConn: %d %d", y1, x1)
	c, ok := connectors[grid[y1][x1]]
	if !ok {
		return 0, nil
	}
	y2 := y1 + c[0]
	x2 := x1 + c[1]
	log.Printf("getConn: %d %d", y2, x2)
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
	} else {
		log.Printf("(%d, %d) %c and (%d, %d) %c not connected", y1, x1, grid[y1][x1], y2, x2, grid[y2][x2])
	}
	y2 = y1 + c[2]
	x2 = x1 + c[3]
	log.Printf("getConn: %d %d", y2, x2)
	is2 := false
	if isConnected(grid, y1, x1, y2, x2) {
		res2 = make([]int, 2)
		res2[0] = y2
		res2[1] = x2
		turn2 = visited[y2][x2]
		is2 = true
	} else {
		log.Printf("(%d, %d) %c and (%d, %d) %c not connected", y1, x1, grid[y1][x1], y2, x2, grid[y2][x2])
	}
	log.Printf("getConn: turn1: %d turn2: %d", turn1, turn2)
	if turn1 != -1 {
		log.Printf("%d %d %c", res2[0], res2[1], grid[res1[0]][res1[1]])
	}
	if turn2 != -1 {
		log.Printf("%d %d %c", res2[0], res2[1], grid[res2[0]][res2[1]])
	}
	if is1 && is2 && turn1 != -1 && turn2 != -1 {
		if turn1 > turn2 {
			return turn1, res1
		} else {
			return turn2, res2
		}
	}
	if is1 && turn1 == -1 {
		log.Printf("getConn: returning: (%d, %d)", res1[0], res1[1])
		return -1, res1
	}
	if is2 && turn2 == -1 {
		log.Printf("getConn: returning: (%d, %d)", res2[0], res2[1])
		return -1, res2
	}
	return -1, nil
}
