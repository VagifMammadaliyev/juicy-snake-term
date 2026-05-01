package entities

import (
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/engine"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

type Brick struct {
	engine.Cell
}

func NewBrick(x, y int) *Brick {
	return &Brick{
		Cell: engine.Cell{
			X:       x,
			Y:       y,
			Height:  DefaultHeight,
			Width:   DefaultWidth,
			FgColor: terminal.Red,
			BgColor: terminal.Red,
			Symbol:  ' ',
		},
	}
}

func (b *Brick) Bounds() engine.Cell {
	return b.Cell
}
