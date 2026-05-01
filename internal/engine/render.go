package engine

import (
	"fmt"
	"io"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

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

func (a *Area) Render(w io.Writer) {
	terminal.EraseScreen(w)
	for _, b := range a.Bounders {
		render(w, b, a.Point)
	}
}

func (a *Area) Reset(w io.Writer, b Bounder) {
	terminal.ResetStyle(w)

	cell := b.Bounds()
	point := startPoint(cell, a.Point)

	for y := 0; y < cell.Height; y++ {
		for x := 0; x < cell.Width; x++ {
			terminal.Erase(w, point.Y+y, point.Y+x)
		}
	}
}

func startPoint(c Cell, offset Point) Point {
	startX, startY := c.X*c.Width+1, c.Y*c.Height+1
	return Point{
		X: startX + offset.X*DefaultWidth,
		Y: startY + offset.Y*DefaultHeight,
	}
}

func render(w io.Writer, b Bounder, offset Point) {
	cell := b.Bounds()
	cell.BgColor.Set(w, terminal.Background)
	cell.FgColor.Set(w, terminal.Foreground)
	defer terminal.ResetStyle(w)

	point := startPoint(cell, offset)

	for y := 0; y < cell.Height; y++ {
		for x := 0; x < cell.Width; x++ {
			terminal.MoveCursor(w, point.Y+y, point.X+x)
			fmt.Fprintf(w, "%c", cell.Symbol)
		}
	}
}
