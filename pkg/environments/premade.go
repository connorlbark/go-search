package environments

import "fmt"

var (
	premadeEnvironments = map[string]Environment{
		"bucharest": nil,
	}
)

func GetEnvironment(name string) (Environment, error) {
	if env, ok := premadeEnvironments[name]; ok {
		return env, nil
	}
	return nil, fmt.Errorf("Cannot find environment with name %s", name)
}

// PremadeEnvironments lists the available pre-made environments
func PremadeEnvironments() []string {
	names := make([]string, len(premadeEnvironments))
	for name := range premadeEnvironments {
		names = append(names, name)
	}
	return names
}
