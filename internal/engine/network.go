package engine

import (
	"bytes"
	"encoding/gob"
)

type EncodedArea struct {
	Point
	Cols  int
	Rows  int
	Cells []Cell
}

func (a *Area) ToEncodedArea() *EncodedArea {
	encodedArea := EncodedArea{
		Point: a.Point,
		Cols:  a.Cols,
		Rows:  a.Rows,
	}

	cells := make([]Cell, 0, len(a.Bounders))
	for _, b := range a.Bounders {
		cells = append(cells, b.Bounds())
	}

	encodedArea.Cells = cells
	return &encodedArea
}

func (ea *EncodedArea) Encode() (*bytes.Buffer, error) {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	if err := encoder.Encode(*ea); err != nil {
		return nil, err
	}

	return &buff, nil
}

func NewAreaFromEncodedArea(encodedArea *EncodedArea) *Area {
	bounders := make([]Bounder, 0, len(encodedArea.Cells))
	for _, cell := range encodedArea.Cells {
		bounders = append(bounders, cell)
	}

	return &Area{
		Point: Point{
			X: encodedArea.Point.X,
			Y: encodedArea.Point.Y,
		},
		Cols:     encodedArea.Cols,
		Rows:     encodedArea.Rows,
		Bounders: bounders,
	}
}
