package algorithms

import (
	"github.com/porgull/go-search/pkg/environments"
)

// HueristicNodeQueue implements heap.Interface to allow
// for a priority queue based on the heuristic value
type HueristicNodeQueue struct {
	Frontier []environments.Node

	NodeCosts             map[string]int
	NodeCostWithHeuristic map[string]int
	NodeIndexes           map[string]int
}

func (q *HueristicNodeQueue) Len() int {
	return len(q.Frontier)
}

func (q *HueristicNodeQueue) Less(i, j int) bool {
	return q.NodeCostWithHeuristic[q.Frontier[i].Name()] < q.NodeCostWithHeuristic[q.Frontier[j].Name()]
}

func (q *HueristicNodeQueue) Swap(i, j int) {

	q.Frontier[i], q.Frontier[j] = q.Frontier[j], q.Frontier[i]

	q.NodeIndexes[q.Frontier[i].Name()] = i
	q.NodeIndexes[q.Frontier[j].Name()] = j
}

// Push adds a new value
func (q *HueristicNodeQueue) Push(x interface{}) {
	node := x.(environments.Node)
	q.NodeIndexes[node.Name()] = len(q.Frontier)
	q.Frontier = append(q.Frontier, node)
}

// Pop returns the value with the lowest
func (q *HueristicNodeQueue) Pop() interface{} {
	old := q.Frontier
	n := len(old)
	poppedValue := old[n-1]
	q.Frontier = old[:n-1]
	delete(q.NodeIndexes, poppedValue.Name())
	return poppedValue
}
