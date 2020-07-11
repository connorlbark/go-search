package search

import (
	"fmt"
	"strings"

	"github.com/porgull/go-search/pkg/environments"
)

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
