package environments

// Environment defines a starting node and a goal node
type Environment interface {
	// GoalNodeName returns the name of the goal node
	IsGoalNode(Node) bool

	// Start returns the start node
	Start() Node

	// Name returns the environment name
	Name() string

	// VisualizeSolution prints the solution
	// to the console
	VisualizeSolution(Node)

	// Validate is ran before the environment
	// is used to ensure that there aren't any
	// setup errors with the env
	Validate() error
}

// Node is a single node of the search space
type Node interface {
	// Name returns the unique name of this
	// node
	Name() string

	// Parent returns the node's parent
	Parent() Node

	// Children returns the node's children
	Children() []Node

	// Cost returns the cost of getting to this node
	// from the parent
	Cost() int

	// Heuristic returns the heuristic value to the goal
	// node
	Heuristic() int

	// Keys returns the steps taken to reach this node
	Steps() []string

	// IsNode returns if nodes are equivalent
	IsNode(Node) bool
}
