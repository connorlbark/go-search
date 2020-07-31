package algorithms

import (
	"container/heap"
	"fmt"
	"strconv"

	"github.com/porgull/go-search/pkg/environments"
	"github.com/porgull/go-search/pkg/search"
)

// SMAStar implements the SMA* search algorithm, which
// is A* with a memory bound specified by the
// 'max_frontier_size' custom param.
type SMAStar struct {
	queue *PriorityNodeQueue

	maxFrontierSize int

	cost              map[string]int
	costWithHeuristic map[string]int
	depth             map[string]int
	successors        map[string][]environments.Node

	iterations int
}

// Run runs A* on the environment and returns the result
func (a SMAStar) Run(ctx search.Context, e environments.Environment) (search.Result, error) {
	if err := a.setParams(ctx); err != nil {
		return search.Result{}, err
	}

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

func (a *SMAStar) setParams(ctx search.Context) error {
	if maxFrontierSizeStr, ok := ctx.CustomSearchParams["max_frontier_size"]; ok {
		maxFrontierSize, err := strconv.ParseInt(maxFrontierSizeStr, 10, 32)
		if err != nil {
			return fmt.Errorf("could not parse 'max_frontier_size' as an integer: %w", err)
		}
		a.maxFrontierSize = int(maxFrontierSize)
	} else {
		fmt.Println("WARNING: Custom param 'max_frontier_size' not supplied. Defaulting to 512...")
		a.maxFrontierSize = 512
	}
	return nil
}

// initialize AStar's fields for this environment
func (a *SMAStar) setStart(start environments.Node) {
	a.cost = make(map[string]int, a.maxFrontierSize)
	a.cost[start.Name()] = 0

	a.costWithHeuristic = make(map[string]int, a.maxFrontierSize)
	a.costWithHeuristic[start.Name()] = start.Heuristic()

	a.depth = make(map[string]int, a.maxFrontierSize)
	a.depth[start.Name()] = 0

	a.successors = make(map[string][]environments.Node, a.maxFrontierSize)

	a.iterations = 0

	a.queue = NewPriorityNodeQueue(start, a.costWithHeuristic, PriorityNodeQueueConfig{
		InitialFrontierAllocation: a.maxFrontierSize,
		NegativeOneIsInfinity:     true,
	})
}

func (a *SMAStar) shallowestNodeWithHighestF() environments.Node {
	worst := a.queue.Frontier[len(a.queue.Frontier)-1]
	worstF := a.costWithHeuristic[worst.Name()]
	worstDepth := a.depth[worst.Name()]
	for i := len(a.queue.Frontier) - 2; i >= 0; i-- {
		curr := a.queue.Frontier[i]
		currF := a.costWithHeuristic[curr.Name()]
		if currF != worstF {
			break
		}
		currDepth := a.depth[curr.Name()]
		if currDepth < worstDepth {
			worst = curr
			worstDepth = currDepth
		}
	}
	return worst
}

func (a *SMAStar) removeNodeFromQueue(node environments.Node) {
	delete(a.cost, node.Name())
	delete(a.depth, node.Name())
	//	delete(a.costWithHeuristic, node.Name())
	delete(a.successors, node.Name())
	heap.Remove(a.queue, a.queue.NodeIndexes[node.Name()])
}

func (a *SMAStar) getNextSuccesor(parent environments.Node) environments.Node {
	for _, node := range a.successors[parent.Name()] {
		if _, seen := a.costWithHeuristic[node.Name()]; !seen {
			return node
		}
	}
	return nil
}

func (a *SMAStar) allChildrenExplored(node environments.Node) bool {
	for _, node := range a.getSuccessors(node) {
		if _, explored := a.costWithHeuristic[node.Name()]; !explored {
			return false
		}
	}
	return true
}

func (a *SMAStar) allSuccessorsInQueue(node environments.Node) bool {
	return a.queue.HasAllInQueue(a.getSuccessors(node))
}

func (a *SMAStar) getSuccessors(node environments.Node) []environments.Node {
	if successors, ok := a.successors[node.Name()]; ok {
		return successors
	}
	children := node.Children()
	a.successors[node.Name()] = children
	return children
}

func (a *SMAStar) getExploredSuccessors(parent environments.Node) []environments.Node {
	successors := a.getSuccessors(parent)
	exploredSuccessors := make([]environments.Node, 0, len(successors))

	for _, node := range successors {
		if _, explored := a.costWithHeuristic[node.Name()]; explored {
			exploredSuccessors = append(exploredSuccessors, node)
		}
	}
	return exploredSuccessors
}

func (a *SMAStar) backup(node environments.Node) {
	_, val := a.bestValOfSuccesors(node)
	if val != a.costWithHeuristic[node.Name()] && node.Parent() != nil {
		a.backup(node.Parent())
	}
}

func (a *SMAStar) bestValOfSuccesors(node environments.Node) (environments.Node, int) {
	var best environments.Node = nil
	bestF := -1

	for _, successor := range a.getExploredSuccessors(node) {
		currF := a.costWithHeuristic[successor.Name()]
		if best == nil || currF < bestF {
			best = successor
			bestF = currF
		}
	}
	return best, bestF
}

func (a *SMAStar) removeFromSuccessors(parent, bad environments.Node) {
	successors := a.getSuccessors(parent)
	for i, node := range successors {
		if node.IsNode(bad) {
			// remove the node at i, because it's the bad node
			successors[i] = successors[len(successors)-1]
			successors[len(successors)-1] = nil
			successors = successors[:len(successors)-1]
			break // break out of the loop
		}
	}
}

// find and return the goal node
func (a *SMAStar) findGoal(e environments.Environment) (environments.Node, error) {
	// if nothing in queue/frontier, then it is impossible
	// to find the goal node
	for a.queue.Len() > 0 {
		a.iterations++
		currentNode := heap.Pop(a.queue).(environments.Node)

		if e.IsGoalNode(currentNode) {
			return currentNode, nil
		}

		successor := a.getNextSuccesor(currentNode)
		if successor != nil {
			a.cost[successor.Name()] = a.cost[currentNode.Name()] + successor.Cost()
			a.costWithHeuristic[successor.Name()] = smamax(a.costWithHeuristic[currentNode.Name()], a.cost[successor.Name()]+successor.Heuristic())
			a.depth[successor.Name()] = a.depth[currentNode.Name()] + 1
		}

		if a.allChildrenExplored(currentNode) {
			a.backup(currentNode)
		}

		if a.allSuccessorsInQueue(currentNode) {
			a.removeNodeFromQueue(currentNode)
		}

		if len(a.queue.Frontier) >= a.maxFrontierSize {
			worstNode := a.shallowestNodeWithHighestF()
			a.removeNodeFromQueue(worstNode)
			worstParent := worstNode.Parent()
			if worstParent != nil {
				a.removeFromSuccessors(worstParent, worstNode)
				if a.queue.HasNodeInQueue(worstParent) == false {
					heap.Push(a.queue, worstParent)
				}
			}
		}

		if successor != nil {
			heap.Push(a.queue, successor)
		}

	}
	return nil, fmt.Errorf("frontier is empty; searched entire space up to a maximum depth of the frontier size (%d), but could not find goal state", a.maxFrontierSize)
}

func smamax(a, b int) int {
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
