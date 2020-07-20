package algorithms

import (
	"container/heap"
	"fmt"

	"github.com/porgull/go-search/pkg/environments"
	"github.com/porgull/go-search/pkg/search"
)

// BreadthFirst implements the breadth first
// search algorithm
type BreadthFirst struct {
	queue *PriorityNodeQueue

	depth map[string]int

	iterations int
}

// Run runs A* on the environment and returns the result
func (a BreadthFirst) Run(ctx search.Context, e environments.Environment) (search.Result, error) {
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

// initialize fields for this environment
func (a *BreadthFirst) setStart(start environments.Node) {
	a.depth = make(map[string]int, 512)
	a.depth[start.Name()] = 1

	a.iterations = 0

	a.queue = NewPriorityNodeQueue(start, a.depth, PriorityNodeQueueConfig{})
}

// find and return the goal node
func (a *BreadthFirst) findGoal(e environments.Environment) (environments.Node, error) {
	// if nothing in queue/frontier, then it is impossible
	// to find the goal node
	for a.queue.Len() > 0 {
		a.iterations++
		currentNode := heap.Pop(a.queue).(environments.Node)

		if e.IsGoalNode(currentNode) {
			return currentNode, nil
		}

		for _, child := range currentNode.Children() {
			prevDepth, seen := a.depth[child.Name()]
			currDepth := a.depth[currentNode.Name()] + 1

			if !seen {
				a.depth[child.Name()] = currDepth
				// new node, just add it
				heap.Push(a.queue, child)
			} else {
				if prevDepth < currDepth {
					continue
				}
				a.depth[child.Name()] = currDepth
				// we found a new route to the node. let's
				// update the depth and fix its placement in the
				// queue
				childIdx := a.queue.NodeIndexes[child.Name()]
				a.queue.Frontier[childIdx] = child
				heap.Fix(a.queue, childIdx)
			}
		}

	}
	return nil, fmt.Errorf("frontier is empty; searched entire space, but could not find goal state")
}
