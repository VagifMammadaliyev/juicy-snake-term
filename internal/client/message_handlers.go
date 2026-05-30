package client

import (
	"fmt"
	"io"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/engine"
)

type msgHandler struct {
	handleFunc func(io.Writer, *hud, []byte) error
	hud        *hud
}

func newMsgHandler(handleFunc func(
	io.Writer,
	*hud,
	[]byte,
) error, hud *hud) *msgHandler {
	return &msgHandler{
		handleFunc: handleFunc,
		hud:        hud,
	}
}

func (m *msgHandler) Handle(w io.Writer, data []byte) error {
	return m.handleFunc(w, m.hud, data)
}

func areaUpdateHandler(w io.Writer, h *hud, buffer []byte) error {
	area, playerSnakeHead, err := engine.NewAreaFromBytes(buffer)
	if err != nil {
		return fmt.Errorf("can't decode server update: %v", err)
	}

	camera := engine.NewCenteredCamera(playerSnakeHead)
	area.RenderForCamera(w, camera)

	di := debugInfo{
		bytesReceived:   len(buffer),
		decodedAreaCols: int(area.Cols),
		decodedAreaRows: int(area.Rows),
		boundersCount:   len(area.Bounders),
	}
	h.debug = di
	h.playerPosition = playerSnakeHead
	h.playerScore = -1
	h.Point = engine.Point{X: 0, Y: area.Rows + 1}

	return nil
}
