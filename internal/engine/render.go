package engine

import (
	"fmt"
	"io"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

type Cell struct {
	X       int
	Y       int
	Height  int
	Width   int
	FgColor terminal.Color
	BgColor terminal.Color
	Symbol  rune
}

type Bounder interface {
	Bounds() Cell
}

func Render(w io.Writer, b Bounder) {
	cell := b.Bounds()
	cell.BgColor.Set(w, terminal.Background)
	cell.FgColor.Set(w, terminal.Foreground)
	defer terminal.ResetStyle(w)

	startX, startY := cell.X*cell.Width+1, cell.Y*cell.Height+1

	for y := 0; y < cell.Height; y++ {
		for x := 0; x < cell.Width; x++ {
			terminal.MoveCursor(w, startY+y, startX+x)
			fmt.Fprintf(w, "%c", cell.Symbol)
		}
	}
}

func Reset(w io.Writer, b Bounder) {
	terminal.ResetStyle(w)

	cell := b.Bounds()
	startX, startY := cell.X*cell.Width+1, cell.Y*cell.Height+1

	for y := 0; y < cell.Height; y++ {
		for x := 0; x < cell.Width; x++ {
			terminal.Erase(w, startY+y, startX+x)
		}
	}
}

func RenderAll(w io.Writer, bounders []Bounder) {
	terminal.EraseScreen(w)
	for _, r := range bounders {
		Render(w, r)
	}
}
