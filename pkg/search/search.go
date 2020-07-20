package search

import (
	"fmt"
	"strings"

	"github.com/porgull/go-search/pkg/environments"
)

// Context passes search args to search algorithms
type Context struct {
	CustomSearchParams CustomSearchParams
}

// CustomSearchParams contains any custom values
// a search algorithm would need, e.g. max depth
// for depth limited search
type CustomSearchParams map[string]string

// ParseCustomSearchParams parses custom search params from
// a user inputted string
func ParseCustomSearchParams(val string) (CustomSearchParams, error) {
	out := make(CustomSearchParams)
	splitKeyVals := strings.Split(val, " ")
	for _, keyval := range splitKeyVals {
		splitKeyVal := strings.SplitN(keyval, "=", 1)
		if len(splitKeyVal) != 2 {
			return CustomSearchParams{}, fmt.Errorf("Custom search parameters have incorrect format: key/val '%s' should have an equals, but does not", keyval)
		}
		key, val := splitKeyVal[0], splitKeyVal[1]
		out[key] = val
	}
	return out, nil
}

// Result contains statistics about a single run
// of the search algorithm on an environment
type Result struct {
	Node        environments.Node
	Environment environments.Environment
	Iterations  int
}

// Print prints the results to stdout
func (r Result) Print() {
	fmt.Printf("Found node %s in %d interations.\n", r.Node.Name(), r.Iterations)
	steps := r.Node.Steps()

	fmt.Printf("Steps (%d): %s\n", len(steps), strings.Join(steps, ", "))
	fmt.Println("Total cost of solution:", r.TotalCost())
	r.Environment.VisualizeSolution(r.Node)
}

// TotalCost returns the total cost of the
// steps taken
func (r Result) TotalCost() int {
	total := 0
	parent := r.Node
	for parent != nil {
		total += parent.Cost()
		parent = parent.Parent()
	}
	return total
}
