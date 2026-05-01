package engine

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
// Very inefficent implementation.
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
