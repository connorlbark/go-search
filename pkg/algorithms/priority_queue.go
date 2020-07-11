package algorithms

import (
	"github.com/porgull/go-search/pkg/environments"
)

// HueristicNodeQueue implements heap.Interface to allow
// for a priority queue based on the heuristic value
type PriorityNodeQueue struct {
	Frontier []environments.Node

	PriorityMap map[string]int
	NodeIndexes map[string]int
}

func (q *PriorityNodeQueue) Len() int {
	return len(q.Frontier)
}

func (q *PriorityNodeQueue) Less(i, j int) bool {
	return q.PriorityMap[q.Frontier[i].Name()] < q.PriorityMap[q.Frontier[j].Name()]
}

func (q *PriorityNodeQueue) Swap(i, j int) {

	q.Frontier[i], q.Frontier[j] = q.Frontier[j], q.Frontier[i]

	q.NodeIndexes[q.Frontier[i].Name()] = i
	q.NodeIndexes[q.Frontier[j].Name()] = j
}

// Push adds a new value
func (q *PriorityNodeQueue) Push(x interface{}) {
	node := x.(environments.Node)
	q.NodeIndexes[node.Name()] = len(q.Frontier)
	q.Frontier = append(q.Frontier, node)
}

// Pop returns the value with the lowest
func (q *PriorityNodeQueue) Pop() interface{} {
	old := q.Frontier
	n := len(old)
	poppedValue := old[n-1]
	q.Frontier = old[:n-1]
	delete(q.NodeIndexes, poppedValue.Name())
	return poppedValue
}
