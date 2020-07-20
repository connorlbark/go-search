package algorithms

import (
	"container/heap"

	"github.com/porgull/go-search/pkg/environments"
)

// PriorityNodeQueueConfig defines optional
// arguments when creating a priority node queue
type PriorityNodeQueueConfig struct {
	HigherIsBetter bool
}

// NewPriorityNodeQueue initializes a priority queue with the provided
// start node and priority map to use
func NewPriorityNodeQueue(start environments.Node, priorityMap map[string]int, config PriorityNodeQueueConfig) *PriorityNodeQueue {
	frontier := make([]environments.Node, 1, 512)
	frontier[0] = start

	nodeIndexes := make(map[string]int, 512)
	nodeIndexes[start.Name()] = 0

	queue := &PriorityNodeQueue{
		NodeIndexes:             nodeIndexes,
		Frontier:                frontier,
		PriorityMap:             priorityMap,
		PriorityNodeQueueConfig: config,
	}

	heap.Init(queue)

	return queue

}

// PriorityNodeQueue implements heap.Interface to allow
// for a priority queue based on a custom priority map;
// see heap.Interface for usage
type PriorityNodeQueue struct {
	PriorityNodeQueueConfig

	// Frontier is the current
	// open set
	Frontier []environments.Node

	// PriorityMap is referenced
	// to get the priority value
	// for a node. Lowest priority
	// nodes are the ones that are
	// popped
	PriorityMap map[string]int

	// NodeIndexes keeps track
	// of the indexes of nodes
	// based upon their name
	NodeIndexes map[string]int
}

func (q *PriorityNodeQueue) Len() int {
	return len(q.Frontier)
}

func (q *PriorityNodeQueue) Less(i, j int) bool {
	if q.HigherIsBetter {
		return q.PriorityMap[q.Frontier[i].Name()] > q.PriorityMap[q.Frontier[j].Name()]
	}
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
