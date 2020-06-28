package environments

import (
	"fmt"

	"github.com/porgull/go-search/pkg/assets"
)

var (
	premadeEnvironments = map[string]func() (Environment, error){
		"bucharest": func() (Environment, error) {
			return LoadEnvironmentFrom(assets.MustOpen("environments/bucharest.json"))
		},
		"corners": func() (Environment, error) {
			return LoadEnvironmentFrom(assets.MustOpen("environments/corners.json"))
		},
		"maze": func() (Environment, error) {
			return LoadEnvironmentFrom(assets.MustOpen("environments/maze.json"))
		},
	}
)

// GetEnvironment gets a premade environment
func GetEnvironment(name string) (Environment, error) {
	if f, ok := premadeEnvironments[name]; ok {
		return f()
	}
	return nil, fmt.Errorf("Cannot find environment with name %s", name)
}

// PremadeEnvironments lists the available pre-made environments
func PremadeEnvironments() []string {
	names := make([]string, len(premadeEnvironments))

	i := 0
	for name := range premadeEnvironments {
		names[i] = name
		i++
	}
	return names
}
