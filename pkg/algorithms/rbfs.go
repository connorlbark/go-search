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
	cost    map[string]int

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
	a.iterations = 0
}

// find and return the goal node
func (a *RecursiveBestFirstSearch) findGoal(e environments.Environment) (environments.Node, error) {
	node, _ := a.recurse(e, e.Start(), -1, 0)
	if node == nil {
		return nil, fmt.Errorf("explored entire search space, but could not find goal node")
	}
	return node, nil
}

func (a *RecursiveBestFirstSearch) recurse(e environments.Environment, node environments.Node, fLimit int, totalCost int) (environments.Node, int) {
	if e.IsGoalNode(node) {
		return node, 0
	}
	a.iterations++

	children := node.Children()

	for i, child := range children {
		// TODO: fix
		// remove the parent node, it breaks RBFS b/c it will recurse infinitely.
		// not sure why.
		if node.Parent() != nil && node.Parent().IsNode(child) {
			// replace the parent with the last child
			children[i] = children[len(children)-1]
			children[len(children)-1] = nil
			// lop off the last child
			children = children[:len(children)-1]
		}
	}

	fLimits := make(map[string]int, len(children))

	if len(children) == 0 {
		return nil, -1
	}

	// NOTE: we have to ensure in all of these functions that we don't
	// go back to the parent, otherwise it could keep looping forever. this
	// can be seen in the `maze` environment. maybe there should be some way
	// to ask the environment to not pass the parent as a possible route from
	// the child.
	for _, child := range children {
		// totalCost is a running total of the total cost of reaching the parent node
		pathF := totalCost + child.Cost() + child.Heuristic()

		fLimits[child.Name()] = pathF
	}

	num := 0

	for {
		num++
		bestF, bestIdx := -1, -1

		for idx, child := range children {
			if bestIdx == -1 {
				bestF = fLimits[child.Name()]
				bestIdx = idx
			}
			if rbfslessthan(fLimits[child.Name()], bestF) {
				bestF = fLimits[child.Name()]
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
				altF = fLimits[child.Name()]
				setAlt = true
			}

			if rbfslessthan(fLimits[child.Name()], altF) {
				altF = fLimits[child.Name()]
			}
		}

		result, newBest := a.recurse(e, children[bestIdx], rbfsmin(fLimit, altF), totalCost+children[bestIdx].Cost())
		fLimits[children[bestIdx].Name()] = newBest // now that search has been conducted,
		if result != nil {
			return result, 0
		}
	}
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
