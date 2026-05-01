package terminal

import (
	"io"
)

type Control int

const (
	Up Control = 1 << iota
	Down
	Left
	Right
	Quit
)

const (
	CtrlC = 'C' & 0x1f
)

func ParseControl(in []byte, n int) (Control, bool) {
	switch {
	case n == 1 && (in[0] == 'q' || in[0] == CtrlC):
		return Quit, true
	case n == 3 && in[0] == Esc && in[1] == '[':
		switch in[2] {
		case 'A':
			return Up, true
		case 'B':
			return Down, true
		case 'C':
			return Right, true
		case 'D':
			return Left, true
		}
	}
	return 0, false
}

func ListenControl(r io.Reader, keys chan Control, done chan struct{}) {
	buf := make([]byte, 8) // the most complex keyboard combo will be 8 bytes
	for {
		n, err := r.Read(buf)
		if err != nil || n == 0 {
			return
		}
		control, ok := ParseControl(buf, n)
		if ok {
			select {
			case <-done:
				return
			case keys <- control:
			}
		}
	}
}
