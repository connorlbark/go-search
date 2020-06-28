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
	frontier          *NodeQueue
	cost              map[string]int
	costWithHeuristic map[string]int
	iterations        int
}

// Run runs A* on the environment and returns the result
func (a AStar) Run(sctx search.Context, e environments.Environment) (search.Result, error) {
	a.SetStart(e.Start())

	node, err := a.FindGoal(e)
	if err != nil {
		return search.Result{}, err
	}

	return search.Result{
		Node:        node,
		Iterations:  a.iterations,
		Environment: e,
	}, nil
}

func (a *AStar) SetStart(start environments.Node) {
	nodeQueue := NodeQueue(make([]environments.Node, 0, 512))
	a.frontier = &nodeQueue
	heap.Push(a.frontier, start)

	a.cost = make(map[string]int, 512)
	a.cost[start.Name()] = 0

	a.costWithHeuristic = make(map[string]int, 512)
	a.costWithHeuristic[start.Name()] = start.Heuristic()

	a.iterations = 0
}

func (a *AStar) FindGoal(e environments.Environment) (environments.Node, error) {
	for a.frontier.Len() > 0 {
		a.iterations++
		currentNode := heap.Pop(a.frontier).(environments.Node)

		if e.IsGoalNode(currentNode) {
			return currentNode, nil
		}

		currentNodeCost := a.cost[currentNode.Name()]
		for _, child := range currentNode.Children() {
			childCost := currentNodeCost + child.Cost()

			previousChildCost, seen := a.cost[child.Name()]

			// if we found a better route to this node OR this
			// node hasn't been seen yet
			if (!seen) || (seen && previousChildCost > childCost) {
				a.cost[child.Name()] = childCost
				a.costWithHeuristic[child.Name()] = childCost + child.Heuristic()
				heap.Push(a.frontier, child)
			}
		}
	}
	return nil, fmt.Errorf("frontier is empty; searched entire space, but could not find goal state")
}
