package environments

import (
	"fmt"
	"strings"
)

func init() {
	addEnvironmentType("grid", &GridEnvironment{})
}

type GridEnvironment struct {
	GridName string   `json:"grid_name"`
	Grid     []string `json:"grid"`

	gridSize Vector2D
	start    *Vector2D
	end      *Vector2D
}

func (g *GridEnvironment) loadNode(pnt Vector2D, parent *GridNode, direction string, env *GridEnvironment) *GridNode {
	return &GridNode{
		point:     pnt,
		env:       env,
		parent:    parent,
		direction: direction,
	}
}

func (g *GridEnvironment) Name() string {
	return g.GridName
}

func (g *GridEnvironment) Start() Node {
	return g.loadNode(*g.start, nil, "start", g)
}

func (g *GridEnvironment) getNeighbors(node *GridNode) []Node {
	pnt := node.point

	possiblePoints := map[string]Vector2D{
		"up":    pnt.Add(Vector2D{y: -1}),
		"down":  pnt.Add(Vector2D{y: 1}),
		"left":  pnt.Add(Vector2D{x: -1}),
		"right": pnt.Add(Vector2D{x: 1}),
	}

	out := make([]Node, 0)
	for name, vec := range possiblePoints {
		pnt := g.getPoint(vec)
		if pnt.Passable() == false {
			continue
		}

		out = append(out, g.loadNode(vec, node, name, g))
	}

	return out
}

func (g *GridEnvironment) passable(pnt Vector2D) bool {
	return g.getPoint(pnt).Passable()
}

func (g *GridEnvironment) getPoint(pnt Vector2D) gridPoint {
	if pnt.WithinBounds(g.gridSize) == false {
		return Unpassable
	}

	return gridPoint(g.Grid[pnt.y][pnt.x])
}

func (g *GridEnvironment) VisualizeSolution(n Node) {
	solutionGrid := g.copyGrid()

	if solutionGridNode, ok := n.(*GridNode); ok {
		parent := solutionGridNode
		for parent != nil {
			pnt := parent.point

			col := solutionGrid[pnt.y]
			colBytes := []rune(col)
			colBytes[pnt.x] = rune(Path)
			col = string(colBytes)

			solutionGrid[pnt.y] = col

			parent = parent.parent
		}
	}

	fmt.Println(strings.Join(solutionGrid, "\n"))

}

func (g *GridEnvironment) copyGrid() []string {
	out := make([]string, len(g.Grid))
	copy(out, g.Grid)
	return out
}

func (g *GridEnvironment) IsGoalNode(n Node) bool {
	if gridNode, ok := n.(*GridNode); ok {
		return gridNode.point.Equals(*g.end)
	}
	return false
}

func (g *GridEnvironment) Validate() error {
	if g.Grid == nil || len(g.Grid) == 0 {
		return fmt.Errorf("must supply grid when using type grid")
	}

	g.gridSize = Vector2D{
		x: len(g.Grid[0]),
		y: len(g.Grid),
	}

	for y, gridRow := range g.Grid {
		if len(gridRow) != g.gridSize.x {
			return fmt.Errorf("expected all rows to have same size (%d), but row %d was of length %d", g.gridSize.x, y, len(gridRow))
		}
		for x, char := range gridRow {
			point := gridPoint(char)
			if point.Valid() == false {
				return fmt.Errorf("point at (%d,%d) had invalid value: %s", x, y, string(char))
			}

			if point == Start {
				if g.start != nil {
					return fmt.Errorf("multiple start points (char:%s)", string(Start))
				}
				g.start = &Vector2D{
					x: x,
					y: y,
				}
			} else if point == End {
				if g.end != nil {
					return fmt.Errorf("multiple end points (char:%s)", string(End))
				}
				g.end = &Vector2D{
					x: x,
					y: y,
				}
			}
		}
	}

	if g.start == nil {
		return fmt.Errorf("could not find start point (char:%s)", string(Start))
	} else if g.end == nil {
		return fmt.Errorf("could not find end point (char:%s)", string(End))
	}

	return nil
}

type GridNode struct {
	point     Vector2D
	env       *GridEnvironment
	parent    *GridNode
	direction string
}

func (g *GridNode) Heuristic() int {
	return g.point.ManhattanDistanceTo(*g.env.end)
}

func (g *GridNode) Children() []Node {
	return g.env.getNeighbors(g)
}

func (g *GridNode) Name() string {
	return fmt.Sprintf("(%d,%d)", g.point.x, g.point.y)
}

func (g *GridNode) Cost() int {
	return g.env.getPoint(g.point).Cost()
}

func (g *GridNode) Steps() []string {
	// get directions

	names := make([]string, 1, 128)
	names[0] = g.direction
	nextParent := g.parent
	for nextParent != nil {
		names = append(names, nextParent.direction)
		nextParent = nextParent.parent
	}
	reverse(names)
	return names
}

func (g *GridNode) IsNode(other Node) bool {
	if otherGridNode, ok := other.(*GridNode); ok {
		if otherGridNode == nil || g == nil {
			return false
		}
		return otherGridNode.point.Equals(g.point)
	}
	return false
}

func (g *GridNode) Parent() Node {
	if g.parent == nil {
		return nil // this is required in order to allow nil comparisons
	}
	return g.parent
}

type Vector2D struct {
	x int
	y int
}

func (v Vector2D) Add(other Vector2D) Vector2D {
	return Vector2D{
		x: v.x + other.x,
		y: v.y + other.y,
	}
}

func (v Vector2D) WithinBounds(bounds Vector2D) bool {
	if v.x < 0 || v.y < 0 {
		return false
	}
	if v.x >= bounds.x || v.y >= bounds.y {
		return false
	}
	return true
}

func (v Vector2D) ManhattanDistanceTo(other Vector2D) int {
	return int(abs(int32(other.x-v.x)) + abs(int32(other.y-v.y)))
}

func (v Vector2D) Equals(other Vector2D) bool {
	return v.x == other.x && v.y == other.y
}

type gridPoint rune

const (
	Start      gridPoint = '*'
	End        gridPoint = '!'
	Unpassable gridPoint = 'x'

	LowCost  gridPoint = '.'
	MidCost  gridPoint = ','
	HighCost gridPoint = '#'

	Path gridPoint = '●'
)

var (
	valid = []gridPoint{
		Start,
		End,
		Unpassable,

		LowCost,
		MidCost,
		HighCost,
	}

	cost = map[gridPoint]int{
		Start:    1,
		End:      1,
		LowCost:  1,
		MidCost:  2,
		HighCost: 3,
	}
)

func (g gridPoint) Passable() bool {
	return g != Unpassable
}

func (g gridPoint) Cost() int {
	return cost[g]
}

func (g gridPoint) Valid() bool {
	return g.in(valid)
}

func (g gridPoint) in(arr []gridPoint) bool {
	for _, pnt := range arr {
		if g == pnt {
			return true
		}
	}
	return false
}

func abs(n int32) int32 {
	y := n >> 31
	return (n ^ y) - y
}
