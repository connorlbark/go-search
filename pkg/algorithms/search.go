package algorithms

import (
	"fmt"

	"github.com/porgull/go-search/pkg/environments"
	"github.com/porgull/go-search/pkg/search"
)

// Search defines the interface for a search algorithm
type Algorithm interface {
	// Run searches in the environment for the goal node
	Run(sctx search.Context, e environments.Environment) (search.Result, error)
}

var (
	// Algorithms defines the valid
	// search algorithms
	algorithms = map[string]Algorithm{
		"a*": AStar{},
	}
)

func Algorithms() []string {
	names := make([]string, len(algorithms))

	i := 0
	for name := range algorithms {
		names[i] = name
		i++
	}

	return names
}

func GetAlgorithm(name string) (Algorithm, error) {
	if algorithm, ok := algorithms[name]; ok {
		return algorithm, nil
	}
	return nil, fmt.Errorf("could not find algorithm %s", name)
}
