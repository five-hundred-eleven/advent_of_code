package main

import (
	"container/heap"
	"log"
	"os"
	"strings"
)

type Node struct {
	i, j           int
	loss, priority int
	momentum       int
	dir            int
	children       []*Node
	index          int
}

type NodeHeap struct {
	heap       []*Node
	throughput int
}

type Offset struct {
	i, j int
}

var dirMap = []*Offset{
	{i: -1, j: 0},
	{i: 0, j: 1},
	{i: 1, j: 0},
	{i: 0, j: -1},
}

var oppMap = map[int]int{
	0: 2,
	1: 3,
	2: 0,
	3: 1,
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
	grid := make([][]int, 0)
	visited := make([][][][]int, 0)
	for i, line := range strings.Split(contents, "\n") {
		if len(line) == 0 {
			continue
		}
		grid = append(grid, make([]int, len(line)))
		visited = append(visited, make([][][]int, len(line)))
		for j, c := range line {
			grid[i][j] = int(c) - '0'
			visited[i][j] = make([][]int, 4)
			for k := 0; k < 4; k++ {
				visited[i][j][k] = make([]int, 3)
				for l := 0; l < 3; l++ {
					visited[i][j][k][l] = 1 << 31
				}
			}
		}
	}
	sz := &Offset{i: len(grid), j: len(grid[0])}
	dest := &Offset{i: sz.i - 1, j: sz.j - 1}
	destLoss := 1 << 31
	nh := NewNodeHeap()
	heap.Init(nh)
	n0 := &Node{
		i:        0,
		j:        0,
		loss:     0,
		priority: dest.i + dest.j,
		momentum: 0,
		dir:      1,
	}
	heap.Push(nh, n0)
	n1 := &Node{
		i:        0,
		j:        0,
		loss:     0,
		priority: dest.i + dest.j,
		momentum: 0,
		dir:      2,
	}
	heap.Push(nh, n1)
	numPrunedByPriority := 0
	//numPrunedByLoss := 0
	for nh.Len() > 0 {
		curr := heap.Pop(nh).(*Node)
		if curr.i == dest.i && curr.j == dest.j {
			if curr.priority < destLoss {
				destLoss = curr.priority
			}
			continue
		}
		if curr.priority > destLoss {
			numPrunedByPriority++
			continue
		}
		/*
			// in my tests, the following caught literally 0 items
			if curr.loss > visited[curr.i][curr.j][curr.dir][curr.momentum] {
				numPrunedByLoss++
				continue
			}
		*/
		//log.Printf("got node with priority %d", curr.priority)
		exp := curr.explore(sz)
		for _, e := range exp {
			eLoss := curr.loss + grid[e.i][e.j]
			isValid := true
			for fm, f := range visited[e.i][e.j][e.dir] {
				if fm > e.momentum {
					break
				}
				if f <= eLoss {
					isValid = false
					break
				}
			}
			if !isValid {
				continue
			}
			ePriority := aStar(eLoss, dest.i+dest.j-e.i-e.j)
			e.loss = eLoss
			e.priority = ePriority
			heap.Push(nh, e)
			visited[e.i][e.j][e.dir][e.momentum] = eLoss
		}
	}
	/*
		for _, line := range ascii {
			log.Printf("%s", string(line))
		}
	*/
	log.Printf("pruned by priority: %d", numPrunedByPriority)
	//log.Printf("pruned by loss: %d", numPrunedByLoss)
	log.Printf("%d", destLoss)
}

func (n *Node) explore(sz *Offset) (res []*Node) {
	opp, ok := oppMap[n.dir]
	if !ok {
		log.Fatalf("bad dir %c", n.dir)
	}
	isMaxMomentum := n.momentum == 2
	res = make([]*Node, 0)
	for k, v := range dirMap {
		if k == opp {
			continue
		}
		momentum := 0
		if k == n.dir {
			if isMaxMomentum {
				continue
			}
			momentum = n.momentum + 1
		}
		i := n.i + v.i
		if i < 0 || i >= sz.i {
			continue
		}
		j := n.j + v.j
		if j < 0 || j >= sz.j {
			continue
		}
		n2 := &Node{
			i:        i,
			j:        j,
			momentum: momentum,
			dir:      k,
		}
		res = append(res, n2)
	}
	return res
}

/*
func (n *Node) claimChildren(nh *NodeHeap, o *Node) {
	if o.loss <= n.loss {
		log.Printf("claimChildren called incorrectly")
		return
	}
	if o.children == nil || len(o.children) == 0 {
		return
	}
	if o.isDefunct {
		return
	}
	adjust := n.loss - o.loss
	o.isDefunct = true
	n.children = make([]*Node, len(o.children))
	copy(n.children, o.children)
	descendents := make([]*Node, len(o.children))
	copy(descendents, o.children)
	for i := 0; i < len(descendents); i++ {
		curr := descendents[i]
		curr.loss -= adjust
		curr.priority -= adjust
		if curr.isQueued {
			heap.Fix(nh, curr.index)
		}
		if curr.children == nil || len(curr.children) == 0 {
			continue
		}
			if curr.isDefunct {
				continue
			}
		descendents = append(descendents, curr.children...)
	}
}
*/

func NewNodeHeap() (nh *NodeHeap) {
	nh = &NodeHeap{
		heap:       make([]*Node, 0),
		throughput: 0,
	}
	return
}

func (nh NodeHeap) Less(i, j int) bool {
	// Note that this is intentionally reversed
	// 2 reasons:
	//  1. pop from right hand end of stack
	//  2. node.index will remain valid after pop
	return nh.heap[i].priority > nh.heap[j].priority
}

func (nh NodeHeap) Swap(i, j int) {
	nh.heap[i], nh.heap[j] = nh.heap[j], nh.heap[i]
	nh.heap[i].index = i
	nh.heap[j].index = j
}

func (nh NodeHeap) Len() int {
	return len(nh.heap)
}

func (nh *NodeHeap) Push(a any) {
	n := a.(*Node)
	n.index = nh.Len()
	nh.heap = append(nh.heap, n)
}

func (nh *NodeHeap) Pop() (top any) {
	topIndex := nh.Len() - 1
	top = nh.heap[topIndex]
	nh.heap = nh.heap[:topIndex]
	nh.throughput++
	if nh.throughput%100_000 == 0 {
		log.Printf("throughput: %d, len: %d", nh.throughput, topIndex)
	}
	return
}

func aStar(loss int, heuristic int) int {
	return loss + 1*heuristic
}
