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

// Result contains statistics about a single run
// of the search algorithm on an environment
type Result struct {
	Node        environments.Node
	Environment environments.Environment
	Iterations  int

	CustomResultStats map[string]string
}

// Print prints the results to stdout
func (r Result) Print() {
	fmt.Printf("Found node %s in %d iterations.\n", r.Node.Name(), r.Iterations)
	steps := r.Node.Steps()

	fmt.Printf("Steps (%d): %s\n", len(steps), strings.Join(steps, ", "))
	fmt.Println("Total cost of solution:", r.TotalCost())
	r.Environment.VisualizeSolution(r.Node)

	if len(r.CustomResultStats) > 0 {
		fmt.Println("Custom result data for this run:")
		for key, val := range r.CustomResultStats {
			fmt.Printf("%s: %s\n", key, val)
		}
	}
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
