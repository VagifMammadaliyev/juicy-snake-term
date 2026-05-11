package engine

// The code written in this file works by luck honestly, I promise I will go back and check.
// How I made it work: Just changed + to -, - to +, > to >= and so on until visually it seems correct.

import (
	"io"
	"slices"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

type CenteredCamera struct {
	OffsetCols int
	OffsetRows int
	Pivot      Point
}

// TODO: this function has some bugs, very hard to find bugs, find and fix them :(
// but on surface if [CenteredCamera]'s offsets are correctly set it works...
func (a *Area) RenderForCamera(w io.Writer, camera CenteredCamera) {
	terminal.EraseScreen(w)
	offsetCols := camera.OffsetCols
	offsetRows := camera.OffsetRows

	// center of the camera is whatever offset we want, because
	// it should be defined by camera settings
	centerX := camera.OffsetCols
	centerY := camera.OffsetRows

	// if camera requies more space than the area itself
	// then default it to maximum possible for given area.
	// the modulo thing is required so we don't overflow by 1 row/col.
	// nevertheless we want to take into account offsetError in case of
	// overflow prevention, so we render whole space
	offsetColsError := 1 - a.Cols%2
	offsetRowsError := 1 - a.Rows%2

	var (
		maxX, maxY, minX, minY int
	)

	// offset errors should be re-added to either side, not both
	if offsetCols*2+1 > a.Cols {
		offsetCols = centerX - offsetColsError
		maxX = camera.Pivot.X + offsetCols + offsetColsError
	} else {
		maxX = camera.Pivot.X + offsetCols
	}
	minX = camera.Pivot.X - offsetCols

	if offsetRows*2+1 > a.Rows {
		offsetRows = centerY - offsetRowsError
		maxY = camera.Pivot.Y + offsetRows + offsetRowsError
	} else {
		maxY = camera.Pivot.Y + offsetRows
	}
	minY = camera.Pivot.Y - offsetRows

	// difference used to adjust other cells for rendering
	diffX := centerX - camera.Pivot.X
	diffY := centerY - camera.Pivot.Y

	// we need to shift area a little bit so we can draw nice
	// boundary of the camera view
	areaOffsetPoint := Point{
		X: a.Point.X - 1, // why subtract, WHYY??
		Y: a.Point.Y - 1,
	}

	for _, b := range a.Bounders {
		cell := b.Bounds()
		if cell.X <= maxX && cell.X > minX &&
			cell.Y <= maxY && cell.Y > minY {
			// this cell should be visible, apply the diff and render.
			// p.s. for some reason i need to add here, not subtract
			//      although i think i should subtract, but its 2AM so
			//      it works and and that's it
			cell.X += diffX
			cell.Y += diffY
			cell.render(w, areaOffsetPoint)
		}
	}

	// add that camera boundary
	for x := range offsetCols*2 + 1 {
		NewCell(x, 0, terminal.BrightBlack).render(w, a.Point)
		NewCell(x, offsetRows*2, terminal.BrightBlack).render(w, a.Point)
	}

	for y := range offsetRows * 2 {
		NewCell(0, y+1, terminal.BrightBlack).render(w, a.Point)
		NewCell(offsetCols*2, y+1, terminal.BrightBlack).render(w, a.Point)
	}
}

// HACK
func (ea *EncodedArea) RemoveInvisibleCells(pivot Point, offsetX, offsetY int) []Cell {
	visibleCells := make([]Cell, 0, len(ea.Cells))

	maxX := pivot.X + offsetX
	maxY := pivot.Y + offsetY
	minX := pivot.X - offsetX
	minY := pivot.Y - offsetY

	for _, c := range ea.Cells {
		if c.X <= maxX && c.X > minX &&
			c.Y <= maxY && c.Y > minY {
			visibleCells = append(visibleCells, c)
		}
	}

	oldCells := slices.Clone(ea.Cells)
	ea.Cells = visibleCells
	return oldCells
}
