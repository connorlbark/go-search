package search

import (
	"fmt"
	"strings"

	"github.com/porgull/go-search/pkg/environments"
)

// Context passes in flags to the search algorithm
type Context struct {
	Visualize bool
}

// Result contains statistics about a single run
// of the search algorithm on an environment
type Result struct {
	Node  environments.Node
	Steps int
}

// Print prints the results to stdout
func (r Result) Print() {
	fmt.Printf("Found node %s in %d steps.\n", r.Node.Name(), r.Steps)

	fmt.Printf("Steps: %s\n", strings.Join(r.Node.Keys(), ", "))
}
