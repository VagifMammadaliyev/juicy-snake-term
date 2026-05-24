package engine

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/messages"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

// Encode encodes the [Area] and writes it to buff.
func (a *Area) Encode(buff *bytes.Buffer) error {
	b := make([]byte, 0, 5+len(a.Bounders)*5)
	b = append(b, byte(messages.MsgTypeAreaUpdate))
	b = binary.BigEndian.AppendUint16(b, uint16(a.cameraCenter.X))
	b = binary.BigEndian.AppendUint16(b, uint16(a.cameraCenter.Y))

	for _, bounder := range a.Bounders {
		c := bounder.Bounds()
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

// NewAreaFromBytes creates a new [Area] from byte slice. It also returns center [Point]
// which was set initially in [Area.Encode] method.
func NewAreaFromBytes(data []byte) (*Area, Point, error) {
	n := len(data)
	if n < 4 { // 2 bytes + 2 bytes for PlayerSnakeHead (x,y)
		return nil, Point{}, fmt.Errorf("data too short to create Area: %d bytes", n)
	}

	centerX := int16(binary.BigEndian.Uint16(data[0:2]))
	centerY := int16(binary.BigEndian.Uint16(data[2:4]))
	centerPoint := Point{X: centerX, Y: centerY}

	cells := make([]Cell, 0, (n-4)/5) // each cell is 5 bytes (BgColor, X, Y)
	for i := 4; i < n; i += 5 {
		if i+2 >= n {
			return nil, centerPoint, fmt.Errorf("incomplete cell data at index %d", i)
		}
		bgColor := terminal.Color(int16(data[i]))
		x := int16(binary.BigEndian.Uint16(data[i+1 : i+3]))
		y := int16(binary.BigEndian.Uint16(data[i+3 : i+5]))

		cells = append(cells, NewCell(x, y, bgColor))
	}

	bounders := make([]Bounder, 0, len(cells))
	for _, cell := range cells {
		bounders = append(bounders, cell)
	}

	return &Area{
		Point:    Point{0, 0},
		Cols:     cameraOffsetCols*2 + 1,
		Rows:     cameraOffsetRows*2 + 1,
		Bounders: bounders,
	}, centerPoint, nil
}
