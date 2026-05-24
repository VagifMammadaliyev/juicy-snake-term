package engine

import (
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

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
	Symbol  [2]rune // the length should be equal to [cellWidth]
}

const (
	cellHeight int16 = 1
	cellWidth  int16 = 2
)

func NewCell(x, y int16, color terminal.Color) Cell {
	return Cell{
		Point: Point{
			X: x,
			Y: y,
		},
		Height:  cellHeight,
		Width:   cellWidth,
		FgColor: color,
		BgColor: color,
		Symbol:  [2]rune{' ', ' '},
	}
}

func NewCellWithSymbol(x, y int16, color terminal.Color, symbol [2]rune) Cell {
	cell := NewCell(x, y, color)

	cell.Symbol = symbol
	return cell
}

// Bounds return the cell itself
// this it to comply to [Bounder] interface
// when constructing [Area] from network data.
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

	cameraCenter Point
}

func NewArea(cols, rows int16) *Area {
	estimatedCapWithExtra := int(cols)*2 + int(rows)*2 + int(cols)*int(rows)/10

	return &Area{
		Point:    Point{0, 0},
		Cols:     cols,
		Rows:     rows,
		Bounders: make([]Bounder, 0, estimatedCapWithExtra),
	}
}

func (a *Area) Clone() *Area {
	clone := &Area{
		Point:    a.Point,
		Cols:     a.Cols,
		Rows:     a.Rows,
		Bounders: make([]Bounder, len(a.Bounders)),
	}

	copy(clone.Bounders, a.Bounders)
	return clone
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
			Height:  cellHeight,
			Width:   cellWidth,
			FgColor: terminal.Black,
			BgColor: terminal.Black,
			Symbol:  [2]rune{'-', '-'},
		},
	}
}

func (dc DebugCell) Bounds() Cell {
	return dc.Cell
}
