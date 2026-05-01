package terminal

import (
	"fmt"
	"io"
)

type Color int

const (
	Black Color = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

type GraphicRendition int

const (
	Foreground GraphicRendition = 38
	Background GraphicRendition = 48
)

func (c Color) Set(w io.Writer, gr GraphicRendition) {
	fmt.Fprintf(w, "%c[%d;5;%dm", Esc, gr, c)
}

func ResetStyle(w io.Writer) {
	fmt.Fprintf(w, "%c[0m", Esc)
}
