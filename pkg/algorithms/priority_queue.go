package algorithms

import "github.com/porgull/go-search/pkg/environments"

// A priorityQueue implements heap.Interface
type NodeQueue []environments.Node

func (q NodeQueue) Len() int {
	return len(q)
}

func (q NodeQueue) Less(i, j int) bool {
	return q[i].Heuristic() > q[j].Heuristic()
}

func (q NodeQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *NodeQueue) Push(x interface{}) {
	val := x.(environments.Node)
	*q = append(*q, val)
}

func (q *NodeQueue) Pop() interface{} {
	old := *q
	n := len(old)
	poppedValue := old[n-1]
	*q = old[:n-1]
	return poppedValue
}
