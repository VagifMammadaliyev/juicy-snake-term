package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/engine"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/entities"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

type Logic interface {
	UpdateState()
	GetStateForPlayer(playerID string) (*engine.EncodedArea, error)
	AddPlayer(direction entities.SnakeDirection) (playerID string, err error)
	SetPlayerDirection(playerID string, direction entities.SnakeDirection) error
}

type NetworkPlayer struct {
	PlayerID string
	Addr     *net.UDPAddr
}

type GameServer struct {
	conn           *net.UDPConn
	networkPlayers map[string]NetworkPlayer
	serverLock     sync.Mutex
	updateBuffer   bytes.Buffer

	Controls chan terminal.Control
	Logic    Logic
}

func NewGameServer(conn *net.UDPConn, logic Logic) *GameServer {
	controls := make(chan terminal.Control, 2)

	gameServer := &GameServer{
		conn:           conn,
		networkPlayers: make(map[string]NetworkPlayer),
		updateBuffer:   bytes.Buffer{},

		Controls: controls,
		Logic:    logic,
	}

	return gameServer
}

func (g *GameServer) updatePlayers() {
	g.serverLock.Lock()
	defer g.serverLock.Unlock()

	for _, networkPlayer := range g.networkPlayers {
		encodedArea, err := g.Logic.GetStateForPlayer(networkPlayer.PlayerID)
		if err != nil {
			// either player not found or area can't be encoded
			fmt.Printf("can't get player update: %v. removing...\n", err)
			delete(g.networkPlayers, networkPlayer.Addr.String())
			continue
		}

		g.updateBuffer.Reset()
		err = encodedArea.Encode(&g.updateBuffer)
		if err != nil {
			fmt.Printf("can't encode player update: %v\n", err)
		}

		g.conn.SetWriteDeadline(time.Now().Add(80 * time.Millisecond))
		fmt.Printf("player update size for %s: %d bytes\n", networkPlayer.Addr.String(), g.updateBuffer.Len())
		_, err = g.conn.WriteToUDP(g.updateBuffer.Bytes(), networkPlayer.Addr)

		if err != nil {
			// it's OK we will try on the next main tick
			log.Printf("can't send player update %s: %v", networkPlayer.Addr.String(), err)
			continue
		}
	}
}

const serverTick = 80 * time.Millisecond

func (g *GameServer) Run() {
	ticker := time.NewTicker(serverTick)
	defer ticker.Stop()
	statTicker := time.NewTicker(2 * time.Second)
	defer statTicker.Stop()

	for {
		select {
		case <-ticker.C:
			g.Logic.UpdateState()
			g.updatePlayers()

		case <-statTicker.C:
			log.Printf("stats: players connected: %d", len(g.networkPlayers))
		}
	}
}

func (g *GameServer) HandlePlayerConnection(conn *net.UDPConn, data *bytes.Buffer, clientAddr *net.UDPAddr) {
	g.serverLock.Lock()
	defer g.serverLock.Unlock()

	decoder := gob.NewDecoder(data)
	var playerDirection entities.SnakeDirection
	if err := decoder.Decode(&playerDirection); err != nil {
		// shouldn't happen if happens, we just ignore
		return
	}

	networkPlayer, ok := g.networkPlayers[clientAddr.String()]
	if !ok {
		log.Printf("new player connecting: %s", clientAddr.String())
		playerID, err := g.Logic.AddPlayer(playerDirection)

		if err != nil {
			fmt.Printf("can't add new player: %v\n", err)
			return
		}
		g.networkPlayers[clientAddr.String()] = NetworkPlayer{
			PlayerID: playerID,
			Addr:     clientAddr,
		}
		return
	}

	log.Printf("setting player %s direction to %d", clientAddr.String(), playerDirection)
	g.Logic.SetPlayerDirection(networkPlayer.PlayerID, playerDirection)
}
