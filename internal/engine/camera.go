package engine

import (
	"io"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

const (
	cameraOffsetCols = 21
	cameraOffsetRows = 11
)

type CenteredCamera struct {
	OffsetCols int16
	OffsetRows int16
	Pivot      Point
}

func NewCenteredCamera(pivot Point) CenteredCamera {
	return CenteredCamera{
		OffsetCols: cameraOffsetCols,
		OffsetRows: cameraOffsetRows,
		Pivot:      pivot,
	}
}

// setCameraCenter sets the center of the camera to the given pivot point.
// This method allows set center for specific player.
// This method should be called on cloned area, so other players can have different camera centers.
func (a *Area) setCameraCenter(pivot Point) {
	a.cameraCenter = pivot
}

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
		maxX, maxY, minX, minY int16
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
		X: a.Point.X - 1,
		Y: a.Point.Y - 1,
	}

	for _, b := range a.Bounders {
		cell := b.Bounds()
		if cell.X <= maxX && cell.X > minX &&
			cell.Y <= maxY && cell.Y > minY {
			// this cell should be visible, apply the diff and render.
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

// removeInvisibleCells removes cells that should not render for given camera pivot and offsets.
// Recommended to [Area.Clone] the area before calling, so original area state is not lost.
func (a *Area) removeInvisibleCells(pivot Point) {
	visibleCells := make([]Bounder, 0, len(a.Bounders))

	maxX := pivot.X + cameraOffsetCols
	maxY := pivot.Y + cameraOffsetRows
	minX := pivot.X - cameraOffsetCols
	minY := pivot.Y - cameraOffsetRows

	for _, bounders := range a.Bounders {
		c := bounders.Bounds()
		if c.X <= maxX && c.X > minX &&
			c.Y <= maxY && c.Y > minY {
			visibleCells = append(visibleCells, c)
		}
	}

	a.Bounders = visibleCells
	a.Cols = cameraOffsetCols*2 + 1
	a.Rows = cameraOffsetRows*2 + 1
}

// PrepareForCamera prepares area to be encoded for a sprific player,
// given the snake head of that player.
func (a *Area) PrepareForCamera(snakeHead Point) {
	a.removeInvisibleCells(snakeHead)
	a.setCameraCenter(snakeHead)
}
