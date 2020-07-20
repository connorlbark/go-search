package algorithms

import (
	"fmt"
	"strconv"

	"github.com/porgull/go-search/pkg/environments"
	"github.com/porgull/go-search/pkg/search"
)

// IterativeDeepening implements the depth limited
// search algorithm. It can take the `initial_depth`
// custom argument, but by the default the value will be 0
type IterativeDeepening struct {
	queue *PriorityNodeQueue

	iterations int

	initialDepth int
	maxDepth     int
}

// Run runs A* on the environment and returns the result
func (a IterativeDeepening) Run(ctx search.Context, e environments.Environment) (search.Result, error) {
	if err := a.setParams(ctx.CustomSearchParams); err != nil {
		return search.Result{}, err
	}

	return a.getResult(e)
}

func (a *IterativeDeepening) setParams(params search.CustomSearchParams) error {
	initialDepthStr, ok := params["initial_depth"]
	if !ok {
		initialDepthStr = "0"
	}

	parsedInitialDepth, err := strconv.ParseInt(initialDepthStr, 10, 32)
	if err != nil {
		return fmt.Errorf("Could not parse 'initial_depth' as integer: %w", err)
	}

	a.initialDepth = int(parsedInitialDepth)

	maxDepthStr, ok := params["max_depth"]
	if !ok {
		maxDepthStr = "-1"
	}

	parsedMaxDepth, err := strconv.ParseInt(maxDepthStr, 10, 32)
	if err != nil {
		return fmt.Errorf("Could not parse 'max_depth' as integer: %w", err)
	}

	a.maxDepth = int(parsedMaxDepth)

	return nil
}

// find and return the goal node
func (a *IterativeDeepening) getResult(e environments.Environment) (search.Result, error) {
	currDepth := a.initialDepth
	for {
		if a.maxDepth != -1 && currDepth >= a.maxDepth {
			return search.Result{}, fmt.Errorf("reached max depth before finding goal node")
		}

		depthLimited := DepthLimited{}
		depthLimitedCtx := search.Context{
			CustomSearchParams: search.CustomSearchParams{
				"depth_limit": strconv.Itoa(currDepth),
			},
		}
		result, err := depthLimited.Run(depthLimitedCtx, e)
		if err == nil {
			result.Iterations = a.iterations + result.Iterations
			return result, nil
		}

		a.iterations = a.iterations + result.Iterations

		currDepth++
	}
}
