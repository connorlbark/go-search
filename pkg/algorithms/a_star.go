package algorithms

import (
	"container/heap"
	"fmt"

	"github.com/porgull/go-search/pkg/environments"
	"github.com/porgull/go-search/pkg/search"
)

// AStar implements the A* search algorithm. See
// https://en.wikipedia.org/wiki/A*_search_algorithm
// for a quick overview.
type AStar struct {
	queue      *HueristicNodeQueue
	iterations int
}

// Run runs A* on the environment and returns the result
func (a AStar) Run(e environments.Environment) (search.Result, error) {
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

// initialize AStar's fields for this environment
func (a *AStar) setStart(start environments.Node) {
	frontier := make([]environments.Node, 1, 512)
	frontier[0] = start

	cost := make(map[string]int, 512)
	cost[start.Name()] = 0

	costWithHeuristic := make(map[string]int, 512)
	costWithHeuristic[start.Name()] = start.Heuristic()

	nodeIndexes := make(map[string]int, 512)
	nodeIndexes[start.Name()] = 0

	a.iterations = 0

	a.queue = &HueristicNodeQueue{
		Frontier:              frontier,
		NodeCosts:             cost,
		NodeCostWithHeuristic: costWithHeuristic,
		NodeIndexes:           nodeIndexes,
	}

	heap.Init(a.queue)
}

// find and return the goal node
func (a *AStar) findGoal(e environments.Environment) (environments.Node, error) {
	// if nothing in queue/frontier, then it is impossible
	// to find the goal node
	for a.queue.Len() > 0 {
		a.iterations++
		currentNode := heap.Pop(a.queue).(environments.Node)

		if e.IsGoalNode(currentNode) {
			return currentNode, nil
		}

		currentNodeCost := a.queue.NodeCosts[currentNode.Name()]
		for _, child := range currentNode.Children() {
			childCost := currentNodeCost + child.Cost()

			previousChildCost, seen := a.queue.NodeCosts[child.Name()]

			// if we found a better route to this node OR this
			// node hasn't been seen yet
			if (!seen) || (seen && previousChildCost > childCost) {
				a.queue.NodeCosts[child.Name()] = childCost
				a.queue.NodeCostWithHeuristic[child.Name()] = childCost + child.Heuristic()
				if currIdx, inQueue := a.queue.NodeIndexes[child.Name()]; inQueue {
					// if the child is already in the frontier, replace it
					// with this node b/c this node has a lower cost
					a.queue.Frontier[currIdx] = child
					// now that the cost has changed,
					// fix the placement of that node
					heap.Fix(a.queue, currIdx)
				} else {
					heap.Push(a.queue, child)
				}
			}
		}

	}
	return nil, fmt.Errorf("frontier is empty; searched entire space, but could not find goal state")
}
