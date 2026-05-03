package game

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
	g.conn.SetWriteDeadline(time.Now().Add(80 * time.Millisecond))

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
		// by angrility smashing keyboard anyaway
		return err
	}

	return nil
}

func (g *GameClient) listenServer(done chan struct{}) {

	buffer := make([]byte, 1024*41) // 41KB for now, should be enough to handle worst case scenario

	for {
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

		decoder := gob.NewDecoder(bytes.NewReader(buffer[:n]))
		var ea engine.EncodedArea
		if err = decoder.Decode(&ea); err != nil {
			// shouldn't happen, but let's just skip for now
			continue
		}

		area := engine.NewAreaFromEncodedArea(&ea)
		for _, b := range area.Bounders {
			if b == nil {
				fmt.Println("received nil for bounder")
			}
		}
		area.Render(g.screen)
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
