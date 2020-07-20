package algorithms

import (
	"container/heap"
	"fmt"

	"github.com/porgull/go-search/pkg/environments"
	"github.com/porgull/go-search/pkg/search"
)

// GreedyBestFirst implements the greedy best first
// search algorithm
type GreedyBestFirst struct {
	queue *PriorityNodeQueue

	heuristic map[string]int

	iterations int
}

// Run runs A* on the environment and returns the result
func (a GreedyBestFirst) Run(ctx search.Context, e environments.Environment) (search.Result, error) {
	a.setStart(e.Start())

	node, err := a.findGoal(e)
	if err != nil {
		return search.Result{}, err
	}

	return search.Result{
		Node:        node,
		Iterations:  a.iterations,
		Environment: e,
	}, nil
}

// initialize GreedyBestFirst's fields for this environment
func (a *GreedyBestFirst) setStart(start environments.Node) {

	a.heuristic = make(map[string]int, 512)
	a.heuristic[start.Name()] = start.Heuristic()

	a.iterations = 0

	a.queue = NewPriorityNodeQueue(start, a.heuristic, PriorityNodeQueueConfig{})
}

// find and return the goal node
func (a *GreedyBestFirst) findGoal(e environments.Environment) (environments.Node, error) {
	// if nothing in queue/frontier, then it is impossible
	// to find the goal node
	for a.queue.Len() > 0 {
		a.iterations++
		currentNode := heap.Pop(a.queue).(environments.Node)

		if e.IsGoalNode(currentNode) {
			return currentNode, nil
		}

		for _, child := range currentNode.Children() {
			_, seen := a.heuristic[child.Name()]

			// new node, add its heuristic
			if !seen {
				a.heuristic[child.Name()] = child.Heuristic()
				heap.Push(a.queue, child)
			}
		}

	}
	return nil, fmt.Errorf("frontier is empty; searched entire space, but could not find goal state")
}
