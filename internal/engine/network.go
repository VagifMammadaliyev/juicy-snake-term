package engine

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

type EncodedArea struct {
	Point
	Cols            int16
	Rows            int16
	Cells           []Cell
	PlayerSnakeHead Point
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

// Encode encodes the [EncodedArea]. It expects [PlayerSnakeHead] to be set.
func (ea *EncodedArea) Encode(buff *bytes.Buffer) error {
	b := make([]byte, 0)
	b = binary.BigEndian.AppendUint16(b, uint16(ea.PlayerSnakeHead.X))
	b = binary.BigEndian.AppendUint16(b, uint16(ea.PlayerSnakeHead.Y))
	for _, c := range ea.Cells {
		b = append(b, byte(c.BgColor))
		b = binary.BigEndian.AppendUint16(b, uint16(c.X))
		b = binary.BigEndian.AppendUint16(b, uint16(c.Y))
	}

	_, err := buff.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func NewEncodedAreaFromBytes(data []byte) (*EncodedArea, error) {
	n := len(data)
	if n < 4 { // 2 bytes + 2 bytes for PlayerSnakeHead (x,y)
		return nil, fmt.Errorf("data too short to create EncodedArea: %d bytes", n)
	}

	snakeHeadX := int16(binary.BigEndian.Uint16(data[0:2]))
	snakeHeadY := int16(binary.BigEndian.Uint16(data[2:4]))
	playerSnakeHead := Point{X: snakeHeadX, Y: snakeHeadY}

	cells := make([]Cell, 0, (n-4)/5) // each cell is 5 bytes (BgColor, X, Y)
	for i := 4; i < n; i += 5 {
		if i+2 >= n {
			return nil, fmt.Errorf("incomplete cell data at index %d", i)
		}
		bgColor := terminal.Color(int16(data[i]))
		x := int16(binary.BigEndian.Uint16(data[i+1 : i+3]))
		y := int16(binary.BigEndian.Uint16(data[i+3 : i+5]))

		cells = append(cells, Cell{
			Point:   Point{X: x, Y: y},
			BgColor: bgColor,
			FgColor: bgColor,
			Symbol:  ' ',
			Height:  DefaultHeight,
			Width:   DefaultWidth,
		})
	}

	return &EncodedArea{
		Point:           Point{0, 0},
		Cols:            DefaultCameraOffsetCols*2 + 1,
		Rows:            DefaultCameraOffsetRows*2 + 1,
		Cells:           cells,
		PlayerSnakeHead: playerSnakeHead,
	}, nil
}

func NewAreaFromEncodedArea(encodedArea *EncodedArea) *Area {
	bounders := make([]Bounder, 0, len(encodedArea.Cells))
	for _, cell := range encodedArea.Cells {
		bounders = append(bounders, cell)
	}

	return &Area{
		Point:    encodedArea.Point,
		Cols:     encodedArea.Cols,
		Rows:     encodedArea.Rows,
		Bounders: bounders,
	}
}
