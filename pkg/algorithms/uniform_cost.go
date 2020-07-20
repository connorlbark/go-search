package algorithms

import (
	"container/heap"
	"fmt"

	"github.com/porgull/go-search/pkg/environments"
	"github.com/porgull/go-search/pkg/search"
)

// UniformCost implements the uniform cost
// search algorithm
type UniformCost struct {
	queue *PriorityNodeQueue

	cost map[string]int

	iterations int
}

// Run runs A* on the environment and returns the result
func (a UniformCost) Run(ctx search.Context, e environments.Environment) (search.Result, error) {
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
func (a *UniformCost) setStart(start environments.Node) {
	a.cost = make(map[string]int, 512)
	a.cost[start.Name()] = 0

	a.iterations = 0

	a.queue = NewPriorityNodeQueue(start, a.cost, PriorityNodeQueueConfig{})
}

// find and return the goal node
func (a *UniformCost) findGoal(e environments.Environment) (environments.Node, error) {
	// if nothing in queue/frontier, then it is impossible
	// to find the goal node
	for a.queue.Len() > 0 {
		a.iterations++
		currentNode := heap.Pop(a.queue).(environments.Node)

		if e.IsGoalNode(currentNode) {
			return currentNode, nil
		}

		for _, child := range currentNode.Children() {
			childIdx, inQueue := a.queue.NodeIndexes[child.Name()]
			prevCost, seenPrev := a.cost[child.Name()]
			currCost := a.cost[currentNode.Name()] + child.Cost()

			// if the cost of the node we just expanded
			// is higher than the pre-existing node in
			// the queue, skip this iteration
			if seenPrev && prevCost < currCost {
				continue
			}

			if !inQueue {
				a.cost[child.Name()] = a.cost[currentNode.Name()] + child.Cost()
				// new node, just add it
				heap.Push(a.queue, child)
			} else {

				a.cost[child.Name()] = a.cost[currentNode.Name()] + child.Cost()
				// we found a new route to the node. let's
				// update the cost and fix its placement in the
				// queue
				a.queue.Frontier[childIdx] = child
				heap.Fix(a.queue, childIdx)
			}
		}

	}
	return nil, fmt.Errorf("frontier is empty; searched entire space, but could not find goal state")
}
