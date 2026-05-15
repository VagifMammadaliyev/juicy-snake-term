package terminal

import (
	"fmt"
	"io"
)

func MoveCursor(w io.Writer, line, column int16) {
	fmt.Fprintf(w, "%c[%d;%dH", Esc, line, column)
}

func HideCursor(w io.Writer) {
	fmt.Fprintf(w, "%c[?25l", Esc)
}

func RestoreCursor(w io.Writer) {
	fmt.Fprintf(w, "%c[?25h", Esc)
}
