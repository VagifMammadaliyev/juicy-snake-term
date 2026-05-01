package engine

import "github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"

type Point struct {
	X int
	Y int
}

type Cell struct {
	Point
	Height  int
	Width   int
	FgColor terminal.Color
	BgColor terminal.Color
	Symbol  rune
}

func NewCell(x, y int, color terminal.Color) Cell {
	return Cell{
		Point: Point{
			X: x,
			Y: y,
		},
		Height:  DefaultHeight,
		Width:   DefaultWidth,
		FgColor: color,
		BgColor: color,
		Symbol:  ' ',
	}
}
