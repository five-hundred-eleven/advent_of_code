package main

import (
	"log"
	"os"
	"strings"
)

type Beam struct {
	i, j int
	dir  byte
}

type Offset struct {
	i, j int
}

var MirrorMap = map[byte]map[byte][]byte{
	'.': {
		'n': {'n'},
		'e': {'e'},
		's': {'s'},
		'w': {'w'},
	},
	'/': {
		'n': {'e'},
		'e': {'n'},
		's': {'w'},
		'w': {'s'},
	},
	'\\': {
		'n': {'w'},
		'e': {'s'},
		's': {'e'},
		'w': {'n'},
	},
	'|': {
		'n': {'n'},
		'e': {'n', 's'},
		's': {'s'},
		'w': {'n', 's'},
	},
	'-': {
		'n': {'e', 'w'},
		'e': {'e'},
		's': {'e', 'w'},
		'w': {'w'},
	},
}

var DirectionMap = map[byte]*Offset{
	'n': {i: -1, j: 0},
	'e': {i: 0, j: 1},
	's': {i: 1, j: 0},
	'w': {i: 0, j: -1},
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
	for _, line := range strings.Split(contents, "\n") {
		if len(line) == 0 {
			continue
		}
		grid = append(grid, line)
	}

	result := 0
	for i := range grid {
		b := &Beam{i: i, j: 0, dir: 'e'}
		r := countEnergy(grid, b)
		if r > result {
			result = r
		}
	}
	for j := range grid[0] {
		b := &Beam{i: 0, j: j, dir: 's'}
		r := countEnergy(grid, b)
		if r > result {
			result = r
		}
	}
	for i := range grid {
		b := &Beam{i: i, j: len(grid[0]) - 1, dir: 'w'}
		r := countEnergy(grid, b)
		if r > result {
			result = r
		}
	}
	for j := range grid[0] {
		b := &Beam{i: len(grid) - 1, j: j, dir: 'n'}
		r := countEnergy(grid, b)
		if r > result {
			result = r
		}
	}
	log.Printf("%d", result)
}

func countEnergy(grid []string, first *Beam) int {
	energized := make([][]bool, len(grid))
	directional := make(map[byte][][]bool)
	directional['n'] = make([][]bool, len(grid))
	directional['e'] = make([][]bool, len(grid))
	directional['s'] = make([][]bool, len(grid))
	directional['w'] = make([][]bool, len(grid))
	for i := range grid {
		energized[i] = make([]bool, len(grid[i]))
		directional['n'][i] = make([]bool, len(grid[i]))
		directional['e'][i] = make([]bool, len(grid[i]))
		directional['s'][i] = make([]bool, len(grid[i]))
		directional['w'][i] = make([]bool, len(grid[i]))
	}
	sz := &Offset{i: len(grid), j: len(grid[0])}
	q := make([]*Beam, 0)
	q = append(q, first)
	for qIndex := 0; qIndex < len(q); qIndex++ {
		curr := q[qIndex]
		energized[curr.i][curr.j] = true
		if directional[curr.dir][curr.i][curr.j] {
			continue
		}
		directional[curr.dir][curr.i][curr.j] = true
		terrain := grid[curr.i][curr.j]
		for _, newBeam := range curr.Travel(terrain, sz) {
			q = append(q, newBeam)
		}
	}
	result := 0
	for _, row := range energized {
		for _, e := range row {
			if !e {
				continue
			}
			result++
		}
	}
	return result
}

func (b *Beam) Travel(terrain byte, sz *Offset) (bs []*Beam) {
	newDirs := MirrorMap[terrain][b.dir]
	bs = make([]*Beam, 0)
	for _, d := range newDirs {
		o := DirectionMap[d]
		b := &Beam{dir: d, i: b.i + o.i, j: b.j + o.j}
		if b.i < 0 || b.i >= sz.i {
			continue
		}
		if b.j < 0 || b.j >= sz.j {
			continue
		}
		bs = append(bs, b)
	}
	return
}
