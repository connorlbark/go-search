package environments

import (
	"fmt"
)

func init() {
	addEnvironmentType("state", &StateEnvironment{})
}

var _ Environment = &StateEnvironment{}
var _ Node = &StateNode{}

// StateEnvironment loads a static environment
// from a json file
type StateEnvironment struct {
	StartNode       string `json:"start_node"`
	GoalNode        string `json:"goal_node"`
	States          States `json:"states"`
	EnvironmentName string `json:"environment_name"`
}

// Start returns the start node
func (l *StateEnvironment) Start() Node {
	return l.States.loadNode(l.StartNode, 0, nil, l)
}

// IsGoalNode checks if the node has the right name
func (l *StateEnvironment) IsGoalNode(n Node) bool {
	return n.(*StateNode).name == l.GoalNode
}

// Name returns the name of the environment
func (l *StateEnvironment) Name() string {
	return l.EnvironmentName
}

// VisualizeSolution prints out the steps it took
func (l *StateEnvironment) VisualizeSolution(solution Node) {
}

// Validate checks if the environment is valid
// (e.g. all child states resolve to a real state)
func (l *StateEnvironment) Validate() error {
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

// StateNode contains the current node data
// and a pointer to the environment
type StateNode struct {
	env    *StateEnvironment
	name   string
	parent *StateNode
	cost   int
}

// Keys returns the steps it took to get here
func (n *StateNode) Steps() []string {

	names := make([]string, 1, 128)
	names[0] = n.name
	nextParent := n.parent
	for nextParent != nil {
		names = append(names, nextParent.name)
		nextParent = nextParent.parent
	}
	reverse(names)
	return names
}

// Name returns the name of the node
func (n *StateNode) Name() string {
	return n.name
}

// Children returns the children of the
// node
func (n *StateNode) Children() []Node {
	currentState := n.env.States[n.name]
	out := make([]Node, len(currentState.Children))

	i := 0
	for childName, cost := range currentState.Children {
		out[i] = n.env.States.loadNode(childName, cost, n, n.env)
		i++
	}

	return out
}

// Parent returns the parent of the node
func (n *StateNode) Parent() Node {
	if n.parent == nil {
		return nil // this is required in order to allow nil comparisons
	}
	return n.parent
}

// Cost returns the cost of the node
func (n *StateNode) Cost() int {
	return n.cost
}

func (n *StateNode) state() State {
	return n.env.States[n.name]
}

// Heuristic returns the node's heuristic value
func (n *StateNode) Heuristic() int {
	return n.state().Heuristic
}

// IsNode compares equality with another node
func (n *StateNode) IsNode(other Node) bool {
	if other == nil || n == nil {
		return false
	}

	if loaded, ok := other.(*StateNode); ok {
		if loaded == nil {
			return false
		}
		return loaded.name == n.name
	}
	return false
}

// States contains all of the states of the environment
type States map[string]State

func (s States) loadNode(name string, cost int, parent *StateNode, env *StateEnvironment) *StateNode {
	return &StateNode{
		env:    env,
		name:   name,
		parent: parent,
		cost:   cost,
	}
}

// State is the json supplied to encode a state
type State struct {
	Children  map[string]int `json:"children"`
	Heuristic int            `json:"heuristic"`
}
