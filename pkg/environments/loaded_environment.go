package environments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// MustLoadEnvironmentFrom loads an environment into memory, panicking on an error
func MustLoadEnvironmentFrom(path string) Environment {
	env, err := LoadEnvironmentFrom(path)
	if err != nil {
		panic(err)
	}
	return env
}

// LoadEnvironmentFrom loads an environment into memory
func LoadEnvironmentFrom(path string) (Environment, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	env := &LoadedEnvironment{}
	err = json.Unmarshal(b, env)
	if err != nil {
		return nil, err
	}

	return env, nil
}

var _ Environment = &LoadedEnvironment{}
var _ Node = &LoadedNode{}

type LoadedEnvironment struct {
	StartNode       string `json:"start_node"`
	GoalNode        string `json:"goal_node"`
	States          States `json:"states"`
	EnvironmentName string `json:"environment_name"`
	visualize       bool
}

func (l *LoadedEnvironment) Start() Node {
	return l.States.loadNode(l.StartNode, 0, nil, l)
}

func (l *LoadedEnvironment) IsGoalNode(n Node) bool {
	return n.(*LoadedNode).name == l.GoalNode
}

func (l *LoadedEnvironment) Name() string {
	return l.EnvironmentName
}

func (l *LoadedEnvironment) Visualize(enabled bool) {
	l.visualize = enabled
}

func (l *LoadedEnvironment) Unmarshal(b []byte) error {
	if err := json.Unmarshal(b, l); err != nil {
		return err
	}

	if _, ok := l.States[l.StartNode]; !ok {
		return fmt.Errorf("start state %s missing from states", l.StartNode)
	}

	if _, ok := l.States[l.GoalNode]; !ok {
		return fmt.Errorf("goal state %s missing from states", l.GoalNode)
	}

	// validate that no child are missing
	for parentName, state := range l.States {
		for childName := range state.Children {
			if _, ok := l.States[childName]; !ok {
				return fmt.Errorf("missing child %s in states (from parent %s)", childName, parentName)
			}
		}
	}

	return nil
}

type LoadedNode struct {
	env    *LoadedEnvironment
	name   string
	parent *LoadedNode
	cost   int
}

func (n *LoadedNode) Keys() []string {
	if n.parent == nil {
		return make([]string, 0, 128) // make it 128 by default so it doesn't have to reallocate
	}

	return append(n.parent.Keys(), n.name)
}

func (n *LoadedNode) Name() string {
	return n.name
}

func (n *LoadedNode) Children() []Node {
	currentState := n.env.States[n.name]
	out := make([]Node, len(currentState.Children))

	i := 0
	for childName, cost := range currentState.Children {
		out[i] = n.env.States.loadNode(childName, cost, n, n.env)
		i++
	}

	return out
}

func (n *LoadedNode) Parent() Node {
	return n.parent
}

func (n *LoadedNode) Cost() int {
	return n.cost
}

func (n *LoadedNode) state() LoadedState {
	return n.env.States[n.name]
}

func (n *LoadedNode) Heuristic() int {
	return n.state().Heuristic
}

func (n *LoadedNode) Visualize() {
	fmt.Println("TODO")
}

type States map[string]LoadedState

func (s States) loadNode(name string, cost int, parent *LoadedNode, env *LoadedEnvironment) *LoadedNode {
	return &LoadedNode{
		env:    env,
		name:   name,
		parent: parent,
		cost:   cost,
	}
}

type LoadedState struct {
	Children  map[string]int `json:"children"`
	Heuristic int            `json:"heuristic"`
}
