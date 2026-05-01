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

type Bounder interface {
	Bounds() Cell
}

type Area struct {
	Point
	Cols     int
	Rows     int
	Bounders []Bounder
}

func NewArea(cols, rows int) *Area {
	return &Area{
		Point:    Point{0, 0},
		Cols:     cols,
		Rows:     rows,
		Bounders: make([]Bounder, 0),
	}
}

type DebugCell struct {
	Cell
}

func NewDebugCell(x, y int) DebugCell {
	return DebugCell{
		Cell: Cell{
			Point: Point{
				X: x,
				Y: y,
			},
			Height:  DefaultHeight,
			Width:   DefaultWidth,
			FgColor: terminal.Black,
			BgColor: terminal.Black,
			Symbol:  '-',
		},
	}
}

func (dc DebugCell) Bounds() Cell {
	return dc.Cell
}
