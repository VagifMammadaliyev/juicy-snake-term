package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/entities"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

type NetworkPlayer struct {
	Player *Player
	Addr   *net.UDPAddr
}

type GameServer struct {
	conn           *net.UDPConn
	networkPlayers map[string]NetworkPlayer
	updateBuffer   bytes.Buffer

	Controls chan terminal.Control
	Logic    *Logic
}

func NewGameServer(conn *net.UDPConn) *GameServer {
	controls := make(chan terminal.Control, 2)

	gameServer := &GameServer{
		conn:           conn,
		networkPlayers: make(map[string]NetworkPlayer),
		updateBuffer:   bytes.Buffer{},

		Controls: controls,
		Logic:    NewLogic(),
	}

	return gameServer
}

func (g *GameServer) updatePlayers() {
	for _, networkPlayer := range g.networkPlayers {
		encodedArea, err := g.Logic.GetUpdateForPLayer(networkPlayer.Player)
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
		_, err = g.conn.WriteToUDP(g.updateBuffer.Bytes(), networkPlayer.Addr)

		if err != nil {
			// it's OK we will try on the next main tick
			log.Printf("can't send player update %s: %v", networkPlayer.Addr.String(), err)
			continue
		}
	}
}

func (g *GameServer) Run() {
	ticker := time.NewTicker(150 * time.Millisecond)
	defer ticker.Stop()
	statTicker := time.NewTicker(2 * time.Second)
	defer statTicker.Stop()

	for {
		select {
		case <-ticker.C:
			g.Logic.update()
			g.updatePlayers()

		case <-statTicker.C:
			log.Printf("stats: players connected: %d", len(g.networkPlayers))
		}
	}
}

func (g *GameServer) HandlePlayerConnection(conn *net.UDPConn, data *bytes.Buffer, clientAddr *net.UDPAddr) {
	networkPlayer, ok := g.networkPlayers[clientAddr.String()]
	if !ok {
		log.Printf("new player connecting: %s", clientAddr.String())
		player, err := g.Logic.AddPlayer()
		if err != nil {
			fmt.Printf("can't add new player: %v\n", err)
			return
		}
		g.networkPlayers[clientAddr.String()] = NetworkPlayer{
			Player: player,
			Addr:   clientAddr,
		}
	}

	decoder := gob.NewDecoder(data)
	var playerDirection entities.SnakeDirection
	if err := decoder.Decode(&playerDirection); err != nil {
		// shouldn't happen if happens, we just ignore
		return
	}

	log.Printf("setting player %s direction to %d", clientAddr.String(), playerDirection)
	g.Logic.SetPlayerDirection(networkPlayer.Player, playerDirection)

}
