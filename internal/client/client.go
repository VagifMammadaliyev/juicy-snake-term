package client

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/engine"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

const writeDeadline = 80 * time.Millisecond
const maxServerBytes = 1024 // 1KB

type GameClient struct {
	conn     *net.UDPConn
	screen   *bufio.Writer
	Controls chan terminal.Control
}

func NewGameClient(conn *net.UDPConn, buf *bufio.Writer) *GameClient {

	game := &GameClient{
		conn:     conn,
		screen:   buf,
		Controls: make(chan terminal.Control, 2),
	}

	return game
}

func (g *GameClient) sendControls(control terminal.Control) error {
	g.conn.SetWriteDeadline(time.Now().Add(writeDeadline))

	var controlBuffer bytes.Buffer

	encoder := gob.NewEncoder(&controlBuffer)
	if err := encoder.Encode(control); err != nil {
		// shouldn't happen...
		return err
	}

	buffer := make([]byte, controlBuffer.Len())
	controlBuffer.Read(buffer)

	if _, err := g.conn.Write(buffer); err != nil {
		// just ignore, client will try to resend
		// by angrily smashing the keyboard anyway
		return err
	}

	return nil
}

func (g *GameClient) listenServer(done chan struct{}) {

	buffer := make([]byte, maxServerBytes)

	tickDurations := make([]time.Duration, 0, 100)
	for {
		tickStartTime := time.Now()
		n, _, err := g.conn.ReadFromUDP(buffer)

		if err != nil {
			// can't read from server
			// let's just infinitely try to read in a loop,
			// then we will handle timeouts and bla bla...
			log.Printf("didn't receive anything from server: %v\n", err)

			select {
			case <-done:
				return
			default:
				continue
			}
		}
		tickDurations = append(tickDurations, time.Since(tickStartTime))
		if len(tickDurations) > 100 {
			tickDurations = tickDurations[1:]
		}

		ea, err := engine.NewEncodedAreaFromBytes(buffer[:n])
		if err != nil {
			log.Printf("can't decode server update: %v\n", err)
			continue
		}

		area := engine.NewAreaFromEncodedArea(ea)
		for _, b := range area.Bounders {
			if b == nil {
				fmt.Println("received nil for bounder")
			}
		}
		camera := engine.CenteredCamera{
			OffsetCols: engine.DefaultCameraOffsetCols,
			OffsetRows: engine.DefaultCameraOffsetRows,
			Pivot:      ea.PlayerSnakeHead,
		}

		area.RenderForCamera(g.screen, camera)

		terminal.MoveCursor(g.screen, area.Rows+1, 0)
		fmt.Fprintf(g.screen, "Bytes received: %5d\n\r", n)
		fmt.Fprintf(g.screen, "Area size: %2dx%2d\n\r", area.Cols, area.Rows)
		fmt.Fprintf(g.screen, "Bounders: %4d\n\r", len(area.Bounders))
		fmt.Fprintf(g.screen, "Per bounder: %4.2f bytes\n\r", float64(n)/float64(len(area.Bounders)))
		fmt.Fprintf(g.screen, "Player snake head: %2d,%2d\n\r", ea.PlayerSnakeHead.X, ea.PlayerSnakeHead.Y)

		averageTickDuration := time.Duration(0)
		for _, d := range tickDurations {
			averageTickDuration += d
		}
		averageTickDuration /= time.Duration(len(tickDurations))

		fmt.Fprintf(g.screen, "Tick duration: %s\n\r", averageTickDuration)

		g.screen.Flush()
	}
}

func (g *GameClient) Run(done chan struct{}) {
	go func() {
		for control := range g.Controls {
			if control == terminal.Quit {
				close(done)
				g.conn.Close()
				return
			}
			go g.sendControls(control)
		}
	}()

	g.listenServer(done)
}
