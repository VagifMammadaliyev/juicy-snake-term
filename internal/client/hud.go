package client

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/engine"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

type debugInfo struct {
	bytesReceived   int
	decodedAreaCols int
	decodedAreaRows int
	boundersCount   int
	tickDuration    time.Duration
}

type hud struct {
	engine.Point

	playerScore    int
	playerPosition engine.Point

	debug  debugInfo
	errors []error
}

func (h *hud) render(w io.Writer) {
	defer terminal.ResetStyle(w)

	playerInfo := [][2]string{
		{"Score", fmt.Sprintf("%d", h.playerScore)},
		{"Position", fmt.Sprintf("%d, %d", h.playerPosition.X, h.playerPosition.Y)},
	}

	maxErrors := 3
	errorInfo := make([][2]string, 0, maxErrors)
	for i := len(h.errors) - 1; i >= 0; i-- {
		errorInfo = append(errorInfo, [2]string{"Error", h.errors[i].Error()})
		if len(errorInfo) >= maxErrors {
			break
		}
	}

	debugInfo := [][2]string{
		{"Bytes received", fmt.Sprintf("%d", h.debug.bytesReceived)},
		{"Decoded area", fmt.Sprintf("%d x %d", h.debug.decodedAreaCols, h.debug.decodedAreaRows)},
		{"Bounders count", fmt.Sprintf("%d", h.debug.boundersCount)},
		{"Tick duration", h.debug.tickDuration.String()},
	}

	maxLines := max(len(playerInfo), len(errorInfo), len(debugInfo))

	lastCol := h.doRenderHud(w, terminal.White, playerInfo, 0, maxLines)
	lastCol = h.doRenderHud(w, terminal.Cyan, debugInfo, lastCol, maxLines)
	lastCol = h.doRenderHud(w, terminal.Red, errorInfo, lastCol, maxLines)
}

func (h *hud) doRenderHud(
	w io.Writer,
	color terminal.Color,
	data [][2]string,
	startCol int16,
	maxLines int,
) int16 {
	if len(data) == 0 {
		return startCol
	}

	defer terminal.ResetStyle(w)
	terminal.MoveCursor(w, h.Point.Y, h.Point.X)
	color.Set(w, terminal.Foreground)

	maxLineLen := 0
	for _, item := range data {
		key, value := item[0], item[1]
		line := fmt.Sprintf("%s: %s", key, value)
		if len(line) > maxLineLen {
			maxLineLen = len(line)
		}
	}

	topBorder := fmt.Sprintf("┌%s┐", strings.Repeat("─", maxLineLen))
	bottomBorder := fmt.Sprintf("└%s┘", strings.Repeat("─", maxLineLen))
	sideBorder := "│"

	terminal.MoveCursorRight(w, startCol)
	fmt.Fprintf(w, "%s\n\r", topBorder)

	for _, item := range data {
		key, value := item[0], item[1]
		line := fmt.Sprintf("%s: %s", key, value)
		charPadding := strings.Repeat(" ", maxLineLen-len(line))
		terminal.MoveCursorRight(w, startCol)
		fmt.Fprintf(w, "%s%s%s%s\n\r", sideBorder, line, charPadding, sideBorder)
	}

	for i := 0; i < maxLines-len(data); i++ {
		terminal.MoveCursorRight(w, startCol)
		fmt.Fprintf(w, "%s%s%s\n\r", sideBorder, strings.Repeat(" ", maxLineLen), sideBorder)
	}

	terminal.MoveCursorRight(w, startCol)
	fmt.Fprintf(w, "%s\n\r", bottomBorder)

	return int16(maxLineLen+2) + startCol
}
