package entities

import (
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/engine"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

type Food struct {
	engine.Cell
}

func NewFood(x, y int) *Food {
	return &Food{
		Cell: engine.NewCell(x, y, terminal.Magenta),
	}
}

func (f *Food) Bounds() engine.Cell {
	return f.Cell
}
