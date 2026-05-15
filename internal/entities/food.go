package entities

import (
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/engine"
)

type Food struct {
	engine.Cell
}

func NewFood(x, y int16) *Food {
	return &Food{
		Cell: engine.NewCell(x, y, FoodColor),
	}
}

func (f *Food) Bounds() engine.Cell {
	return f.Cell
}
