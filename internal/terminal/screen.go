package terminal

import (
	"fmt"
	"io"
)

type ModeAction byte

const (
	Enable  ModeAction = 'h'
	Disable ModeAction = 'l'
)

type BufferMode int

const (
	Alternate BufferMode = 1049
)

func (b BufferMode) Set(w io.Writer, action ModeAction) {
	fmt.Fprintf(w, "%c[?%d%c", Esc, b, action)
}

func EraseScreen(w io.Writer) {
	fmt.Fprintf(w, "%c[2J", Esc)
}

func Erase(w io.Writer, line, column int) {
	MoveCursor(w, line, column)
	fmt.Fprint(w, " ")
}
