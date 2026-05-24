package engine

import (
	"fmt"
	"io"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

func (a *Area) Render(w io.Writer) {
	terminal.EraseScreen(w)
	for _, b := range a.Bounders {
		b.Bounds().render(w, a.Point)
	}

	// freePoints := a.CalculateFreePoints()
	//
	// for _, p := range freePoints {
	// 	render(w, NewDebugCell(p.X, p.Y), a.Point)
	// }
}

func (a *Area) Reset(w io.Writer, b Bounder) {
	terminal.ResetStyle(w)

	cell := b.Bounds()
	point := startPoint(cell, a.Point)

	for y := int16(0); y < cell.Height; y++ {
		for x := int16(0); x < cell.Width; x++ {
			terminal.Erase(w, point.Y+y, point.Y+x)
		}
	}
}

func startPoint(c Cell, offset Point) Point {
	startX, startY := c.X*c.Width+1, c.Y*c.Height+1
	return Point{
		X: startX + offset.X*cellWidth,
		Y: startY + offset.Y*cellHeight,
	}
}

func (c Cell) render(w io.Writer, offset Point) {
	c.BgColor.Set(w, terminal.Background)
	c.FgColor.Set(w, terminal.Foreground)
	defer terminal.ResetStyle(w)

	point := startPoint(c, offset)

	for y := int16(0); y < c.Height; y++ {
		for x := int16(0); x < c.Width; x++ {
			terminal.MoveCursor(w, point.Y+y, point.X+x)
			fmt.Fprintf(w, "%c", c.Symbol)
		}
	}
}
