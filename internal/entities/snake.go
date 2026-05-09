package entities

import (
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/engine"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

type SnakeNode struct {
	engine.Cell
}

func (sn *SnakeNode) Bounds() engine.Cell {
	return sn.Cell
}

func newSnakeNode(x, y int) *SnakeNode {
	return &SnakeNode{
		Cell: engine.NewCell(x, y, terminal.Green),
	}
}

type SnakeDirection int

const (
	Up SnakeDirection = 1 << iota
	Down
	Left
	Right
)

type Snake struct {
	Nodes           []*SnakeNode // this should be private
	direction       SnakeDirection
	queuedDirection SnakeDirection
}

func NewSnake(length int, x, y int) *Snake {
	if length == 0 {
		length = 1
	}
	nodes := make([]*SnakeNode, 0, length)
	for i := range length {
		nodes = append(nodes, newSnakeNode(x+i, y))
	}

	return &Snake{
		Nodes:           nodes,
		direction:       Left,
		queuedDirection: Left,
	}
}

func (s *Snake) SetDirection(direction SnakeDirection) {
	switch direction {
	case Up:
		if s.direction != Down {
			s.queuedDirection = direction
		}
	case Down:
		if s.direction != Up {
			s.queuedDirection = direction
		}
	case Left:
		if s.direction != Right {
			s.queuedDirection = direction
		}
	case Right:
		if s.direction != Left {
			s.queuedDirection = direction
		}
	}
}

func (s *Snake) Move() {
	s.direction = s.queuedDirection
	switch s.direction {
	case Up:
		s.moveNodes(0, -1)
	case Down:
		s.moveNodes(0, 1)
	case Left:
		s.moveNodes(-1, 0)
	case Right:
		s.moveNodes(1, 0)
	}
}

func (s *Snake) moveNodes(xdelta, ydelta int) {
	var (
		prevX    int
		prevY    int
		currentX int
		currentY int
	)

	for i, node := range s.Nodes {
		cell := &node.Cell
		currentX, currentY = cell.X, cell.Y
		if i == 0 {
			prevX, prevY = cell.X, cell.Y
			cell.X += xdelta
			cell.Y += ydelta
		} else {
			cell.X, cell.Y = prevX, prevY
			prevX, prevY = currentX, currentY
		}
	}
}

func (s *Snake) Grow() {
	lastNode := s.Nodes[len(s.Nodes)-1]
	newNode := newSnakeNode(lastNode.X, lastNode.Y)
	s.Nodes = append(s.Nodes, newNode)
}

func (s *Snake) Head() *SnakeNode {
	return s.Nodes[0]
}
