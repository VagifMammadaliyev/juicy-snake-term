package engine

import (
	"errors"
	"math/rand"
)

func (a *Area) getBusySlice() []bool {
	busy := make([]bool, a.Cols*a.Rows)
	for _, b := range a.Bounders {
		cell := b.Bounds()
		busy[cell.Y*a.Cols+cell.X] = true
	}
	return busy
}

// CalculateFreePoints returns a slice of unoccupied [Point] objects in current area.
// A little bit better implementation that allocating maps.
func (a *Area) CalculateFreePoints() []Point {
	busy := a.getBusySlice()

	freePoints := make([]Point, 0, a.Cols*a.Rows-len(a.Bounders))

	for x := 0; x < a.Cols; x++ {
		for y := 0; y < a.Rows; y++ {
			if !busy[y*a.Cols+x] {
				freePoints = append(freePoints, Point{x, y})
			}
		}
	}
	return freePoints
}

func (a *Area) Collides(i, j Bounder) bool {
	cellI, cellJ := i.Bounds(), j.Bounds()
	return cellI.X == cellJ.X && cellI.Y == cellJ.Y
}

// maxFreePointGuesses is the maximum amount for [GetRandomFreePoint]
// to try to guess. The value should be decided upon empricial evaluation of performance
const maxFreePointGuesses = 5

func (a *Area) GetRandomFreePoint() (Point, error) {
guesses:
	for range maxFreePointGuesses {
		randomX := rand.Intn(a.Cols)
		randomY := rand.Intn(a.Rows)

		// We also can couple bricks knwoledge into area
		// if this will become real performance issue.
		// we can assume that area always has bricks as borders
		// and do the following for random integer generation
		// so we increase the probability of correct point
		//
		// -2 -> because of 2 borders for axis
		// +1 -> because we need to offset, so random point
		//		 is not on the upper or left border
		//
		// randomX := rand.Intn(a.Cols-2) + 1
		// randomY := rand.Intn(a.Rows-2) + 1

		for _, b := range a.Bounders {
			cell := b.Bounds()
			if cell.X == randomX && cell.Y == randomY {
				continue guesses // let's try again
			}
		}
		return Point{randomX, randomY}, nil
	}

	// otherwise let's go for inefficient way
	freePoints := a.CalculateFreePoints()
	if len(freePoints) == 0 {
		return Point{}, errors.New("no free point available")
	}
	return freePoints[rand.Intn(len(freePoints))], nil
}
