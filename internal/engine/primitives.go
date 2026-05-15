package engine

import "github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"

type Point struct {
	X int16
	Y int16
}

type Cell struct {
	Point
	Height  int16
	Width   int16
	FgColor terminal.Color
	BgColor terminal.Color
	Symbol  rune
}

func NewCell(x, y int16, color terminal.Color) Cell {
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

// Bounds return the cell itself
// this it to comply to [Bounder] interface
// when constructing [Area] from [EncodedArea]
func (c Cell) Bounds() Cell {
	return c
}

type Bounder interface {
	Bounds() Cell
}

type Area struct {
	Point
	Cols     int16
	Rows     int16
	Bounders []Bounder
}

func NewArea(cols, rows int16) *Area {
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

func NewDebugCell(x, y int16) DebugCell {
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
