package entities

import (
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/engine"
)

type Brick struct {
	engine.Cell
}

func NewBrick(x, y int16) *Brick {
	return &Brick{
		Cell: engine.NewCell(x, y, BrickColor),
	}
}

func (b *Brick) Bounds() engine.Cell {
	return b.Cell
}
