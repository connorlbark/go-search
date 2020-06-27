package algorithms

import (
	"fmt"

	"github.com/porgull/go-search/pkg/environments"
	"github.com/porgull/go-search/pkg/search"
)

type AStar struct {
	frontier          map[string]environments.Node
	cost              map[string]int
	costWithHeuristic map[string]int
	steps             int
}

func (a AStar) Run(sctx search.Context, e environments.Environment) (search.Result, error) {
	a.SetStart(e.Start())

	node, err := a.FindGoal(e)
	if err != nil {
		return search.Result{}, err
	}

	return search.Result{
		Node:  node,
		Steps: a.steps,
	}, nil
}

func (a *AStar) SetStart(start environments.Node) {
	a.frontier = make(map[string]environments.Node, 512)
	a.frontier[start.Name()] = start

	a.cost = make(map[string]int, 512)
	a.cost[start.Name()] = 0

	a.costWithHeuristic = make(map[string]int, 512)
	a.costWithHeuristic[start.Name()] = start.Heuristic()

	a.steps = 0
}

func (a *AStar) FindGoal(e environments.Environment) (environments.Node, error) {
	for len(a.frontier) > 0 {
		a.steps++
		currentNode := a.LowestCostWithHeuristic()

		if e.IsGoalNode(currentNode) {
			return currentNode, nil
		}

		a.RemoveNodeFromFrontier(currentNode)

		currentNodeCost := a.cost[currentNode.Name()]
		for _, child := range currentNode.Children() {
			childCost := currentNodeCost + child.Cost()

			previousChildCost, seen := a.cost[child.Name()]

			// if we found a better route to this node OR this
			// node hasn't been seen yet
			if (!seen) || (seen && previousChildCost > childCost) {
				a.cost[child.Name()] = childCost
				a.costWithHeuristic[child.Name()] = childCost + child.Heuristic()
				a.AddNodeToFrontier(child)
			}
		}
	}
	return nil, fmt.Errorf("frontier is empty; searched entire space, but could not find goal state")
}

func (a *AStar) RemoveNodeFromFrontier(node environments.Node) {
	delete(a.frontier, node.Name())
}

func (a *AStar) AddNodeToFrontier(node environments.Node) {
	a.frontier[node.Name()] = node
}

func (a *AStar) LowestCostWithHeuristic() environments.Node {
	var lowestNode environments.Node = nil
	var lowestTotalScore int

	for _, node := range a.frontier {
		if (lowestTotalScore > a.costWithHeuristic[node.Name()]) || (lowestNode == nil) {
			lowestNode = node
			lowestTotalScore = a.costWithHeuristic[node.Name()]
		}
	}

	return lowestNode
}
