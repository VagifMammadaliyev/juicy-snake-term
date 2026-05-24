package client

import (
	"fmt"
	"io"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/engine"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/messages"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

func init() {
	messages.RegisterHandler(messages.MsgTypeAreaUpdate, msgHandler(areaUpdateHandler))
}

type msgHandler func(io.Writer, []byte) error

func (h msgHandler) Handle(w io.Writer, data []byte) error {
	return h(w, data)
}

func areaUpdateHandler(w io.Writer, buffer []byte) error {
	area, playerSnakeHead, err := engine.NewAreaFromBytes(buffer)
	if err != nil {
		return fmt.Errorf("can't decode server update: %v", err)
	}

	for _, b := range area.Bounders {
		if b == nil {
			fmt.Println("received nil for bounder")
		}
	}
	camera := engine.NewCenteredCamera(playerSnakeHead)

	area.RenderForCamera(w, camera)

	terminal.MoveCursor(w, area.Rows+1, 0)
	fmt.Fprintf(w, "Bytes received: %5d\n\r", len(buffer))
	fmt.Fprintf(w, "Area size: %2dx%2d\n\r", area.Cols, area.Rows)
	fmt.Fprintf(w, "Bounders: %4d\n\r", len(area.Bounders))
	fmt.Fprintf(w, "Per bounder: %4.2f bytes\n\r", float64(len(buffer))/float64(len(area.Bounders)))
	fmt.Fprintf(w, "Player snake head: %2d,%2d\n\r", playerSnakeHead.X, playerSnakeHead.Y)

	return nil
}
