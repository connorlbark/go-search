package algorithms

import (
	"fmt"

	"github.com/porgull/go-search/pkg/environments"
	"github.com/porgull/go-search/pkg/search"
)

// RecursiveBestFirstSearch implements the RBFS, a modification
// and/or optimization of A*
type RecursiveBestFirstSearch struct {
	queue *PriorityNodeQueue

	fLimits map[string]int

	iterations int
}

// Run runs A* on the environment and returns the result
func (a RecursiveBestFirstSearch) Run(ctx search.Context, e environments.Environment) (search.Result, error) {
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

// initialize RecursiveBestFirstSearch's fields for this environment
func (a *RecursiveBestFirstSearch) setStart(start environments.Node) {
	a.fLimits = make(map[string]int, 512)
	a.fLimits[start.Name()] = -1

	a.iterations = 0

}

// find and return the goal node
func (a *RecursiveBestFirstSearch) findGoal(e environments.Environment) (environments.Node, error) {
	node, _ := a.recurse(e, e.Start(), -1)
	if node == nil {
		return nil, fmt.Errorf("explored entire search space, but could not find goal node")
	}
	return node, nil
}

func (a *RecursiveBestFirstSearch) recurse(e environments.Environment, node environments.Node, fLimit int) (environments.Node, int) {
	if e.IsGoalNode(node) {
		return node, 0
	}
	a.iterations++

	children := node.Children()

	if len(children) == 0 {
		fmt.Println("no children")
		return nil, -1
	}

	for _, child := range children {
		//_, hasPrevFLimit := a.fLimits[child.Name()]
		currChildF := a.totalCost(child) + child.Heuristic()

		a.fLimits[child.Name()] = rbfsmax(currChildF, a.fLimits[node.Name()])

	}

	for {
		bestF, bestIdx := -1, -1

		for idx, child := range children {
			if bestIdx == -1 {
				bestF = a.fLimits[child.Name()]
				bestIdx = idx
			}
			if rbfslessthan(a.fLimits[child.Name()], bestF) {
				bestF = a.fLimits[child.Name()]
				bestIdx = idx
			}
		}

		if rbfsgreaterthan(bestF, fLimit) {
			return nil, bestF
		}

		altF, setAlt := -1, false
		for idx, child := range children {
			if idx == bestIdx {
				continue
			}
			if !setAlt {
				altF = a.fLimits[child.Name()]
				setAlt = true
			}

			if rbfslessthan(a.fLimits[child.Name()], altF) {
				altF = a.fLimits[child.Name()]
			}
		}

		result, bestF := a.recurse(e, children[bestIdx], rbfsmin(fLimit, altF))
		a.fLimits[children[bestIdx].Name()] = bestF // now that search has been conducted,
		if result != nil {
			return result, 0
		}
	}
}

func (a *RecursiveBestFirstSearch) totalCost(child environments.Node) int {
	cost := child.Cost()
	parent := child.Parent()
	for parent != nil {
		cost += parent.Cost()
		parent = parent.Parent()
	}
	return cost
}

func rbfsmax(a, b int) int {
	// -1 == inf
	if a == -1 {
		return b
	}
	if b == -1 {
		return a
	}
	if a > b {
		return a
	}
	return b
}

func rbfsmin(a, b int) int {
	// -1 == inf
	if a == -1 {
		return b
	}
	if b == -1 {
		return a
	}
	if a > b {
		return b
	}
	return a
}

func rbfslessthan(a, b int) bool {
	if a == -1 && b != -1 {
		return false
	}
	if b == -1 && a != -1 {
		return true
	}
	return a < b
}

func rbfsgreaterthan(a, b int) bool {
	if a == -1 && b != -1 {
		return true
	}
	if b == -1 && a != -1 {
		return false
	}
	return a > b
}
