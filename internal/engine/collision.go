package engine

import (
	"errors"
	"math/rand"
)

func (a *Area) getBusyPointMap() map[int]map[int]struct{} {
	busy := make(map[int]map[int]struct{}, len(a.Bounders))
	for _, b := range a.Bounders {
		cell := b.Bounds()
		_, ok := busy[cell.X]
		if !ok {
			busy[cell.X] = make(map[int]struct{})
		}
		busy[cell.X][cell.Y] = struct{}{}
	}
	return busy
}

// CalculateFreePoints returns a slice of unoccupied [Point] objects in current area.
// Very inefficient implementation.
func (a *Area) CalculateFreePoints() []Point {
	busy := a.getBusyPointMap()
	freePoints := make([]Point, 0, a.Cols*a.Rows)
	for x := 0; x < a.Cols; x++ {
		for y := 0; y < a.Rows; y++ {
			xRow, ok := busy[x]
			if !ok {
				freePoints = append(freePoints, Point{x, y})
				continue
			}
			_, ok = xRow[y]
			if !ok {
				freePoints = append(freePoints, Point{x, y})
				continue
			}
		}
	}
	return freePoints
}

func (a *Area) Collides(i, j Bounder) bool {
	cellI, cellJ := i.Bounds(), j.Bounds()
	return cellI.X == cellJ.X && cellI.Y == cellJ.Y
}

func (a *Area) GetRandomFreePoint() (Point, error) {
	freePoints := a.CalculateFreePoints()
	if len(freePoints) == 0 {
		return Point{}, errors.New("no free point available")
	}
	return freePoints[rand.Intn(len(freePoints))], nil
}
