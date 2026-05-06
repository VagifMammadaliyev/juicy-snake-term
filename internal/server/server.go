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

type Player struct {
	Snake *entities.Snake
	Addr  *net.UDPAddr
}

type GameServer struct {
	conn *net.UDPConn

	area     *engine.Area
	Controls chan terminal.Control

	updateLock            sync.Mutex
	encodedAreaBufferPool sync.Pool

	players map[string]*Player
	bricks  []*entities.Brick
	foods   []*entities.Food
}

func NewGameServer(conn *net.UDPConn) *GameServer {
	area := engine.NewArea(1111, 101)
	bricks := make([]*entities.Brick, 0, area.Cols*2+area.Rows*2-4)
	players := make(map[string]*Player)
	controls := make(chan terminal.Control, 2)
	foods := make([]*entities.Food, 0, 500)

	game := &GameServer{
		conn:     conn,
		area:     area,
		bricks:   bricks,
		Controls: controls,
		players:  players,
		encodedAreaBufferPool: sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
		foods: foods,
	}

	game.addBricks()
	return game
}

func (g *GameServer) addBricks() {
	for x := range g.area.Cols {
		g.bricks = append(g.bricks, entities.NewBrick(x, 0))
		g.bricks = append(g.bricks, entities.NewBrick(x, g.area.Rows-1))
	}

	for y := range g.area.Rows - 2 {
		g.bricks = append(g.bricks, entities.NewBrick(0, y+1))
		g.bricks = append(g.bricks, entities.NewBrick(g.area.Cols-1, y+1))
	}
}

func (g *GameServer) addFood() {
	if len(g.foods) < 500 {
		randomPoint, err := g.area.GetRandomFreePoint()
		if err != nil {
			// no area for add food...
			// in practice game should end very soon, as the whole space is occupied
			// by snakes and there is a big probability of collision
			return
		}
		g.foods = append(g.foods, entities.NewFood(randomPoint.X, randomPoint.Y))
	}
}

func (g *GameServer) deleteFoods(foodIndexToRemove map[int]bool) []*entities.Food {
	j := 0
	for i, val := range g.foods {
		if !foodIndexToRemove[i] {
			g.foods[j] = val
			j++
		}
	}
	return g.foods[:j]
}

func (g *GameServer) update() {
	g.updateLock.Lock()
	defer g.updateLock.Unlock()

	for _, player := range g.players {
		player.Snake.Move()
	}

	// check collision with other snakes and itself
	for addr, player := range g.players {
		for anotherAddr, anotherPlayer := range g.players {
			for i, node := range anotherPlayer.Snake.Nodes {
				if i == 0 && addr == anotherAddr {
					continue // skip head of the same snake
				}

				if g.area.Collides(player.Snake.Nodes[0], node) {
					delete(g.players, addr)
				}
			}
		}
	}

	// check food collision with any snake
	fmt.Printf("food amount: %d\n", len(g.foods))
	eatenFoods := make(map[int]bool, len(g.players)) // at most "player" amount of foods can be theoritically eaten during one tick
	for i, food := range g.foods {
		for _, player := range g.players {
			snake := player.Snake

			if g.area.Collides(snake.Nodes[0], food) {
				snake.Grow()
				eatenFoods[i] = true
			}
		}
	}

	g.foods = g.deleteFoods(eatenFoods)
	fmt.Printf("food amount after delete: %d\n", len(g.foods))

	// check brick collision with any snake
	for _, b := range g.bricks {
		for key, player := range g.players {
			snake := player.Snake
			if g.area.Collides(snake.Nodes[0], b) {
				delete(g.players, key)
			}
		}
	}

	g.area.Bounders = g.area.Bounders[:0]

	for _, b := range g.bricks {
		g.area.Bounders = append(g.area.Bounders, b)
	}

	for _, player := range g.players {
		s := player.Snake

		for _, n := range s.Nodes {
			g.area.Bounders = append(g.area.Bounders, n)
		}
	}

	g.addFood()
	for _, food := range g.foods {
		g.area.Bounders = append(g.area.Bounders, food)
	}
}

func (g *GameServer) getAreaBuff() (*bytes.Buffer, func()) {
	buff := g.encodedAreaBufferPool.Get().(*bytes.Buffer)
	buff.Reset()
	return buff, func() {
		g.encodedAreaBufferPool.Put(buff)
	}
}

func (g *GameServer) updatePlayers() {
	encAreaBuff, putBack := g.getAreaBuff()
	defer putBack()

	g.updateLock.Lock()

	addrs := make([]*net.UDPAddr, 0, len(g.players))
	snakeHeads := make([]*entities.SnakeNode, 0, len(g.players))
	for _, player := range g.players {
		addrs = append(addrs, player.Addr)
		snakeHeads = append(snakeHeads, player.Snake.Nodes[0])
	}
	encodedArea := g.area.ToEncodedArea()

	g.updateLock.Unlock()

	for i := range len(addrs) {
		addr := addrs[i]
		snakeHead := snakeHeads[i]
		encodedArea.PlayerSnakeHead = snakeHead.Point
		originalCells := encodedArea.RemoveInvisibleCells(snakeHead.Point, 21)

		err := encodedArea.Encode(encAreaBuff)
		if err != nil {
			// shouldn't happen, but if so, just don't update this player at current tick
			log.Printf("can't encode the area: %v", err)
			continue
		}

		g.conn.SetWriteDeadline(time.Now().Add(80 * time.Millisecond))
		_, err = g.conn.WriteToUDP(encAreaBuff.Bytes(), addr)

		if err != nil {
			// it's OK we will try on the next main tick
			log.Printf("can't update the player %s: %v", addr.String(), err)
			continue
		}

		// so we can filter our for next player
		encAreaBuff.Reset()
		encodedArea.Cells = originalCells
	}
}

func (g *GameServer) Run() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	statTicker := time.NewTicker(2 * time.Second)
	defer statTicker.Stop()

	for {
		select {
		case <-ticker.C:
			g.update()
			g.updatePlayers()

		case <-statTicker.C:
			log.Printf("stats: players connected: %d", len(g.players))
		}
	}
}

func (g *GameServer) HandlePlayerConnection(conn *net.UDPConn, data *bytes.Buffer, clientAddr *net.UDPAddr) {
	g.updateLock.Lock()
	defer g.updateLock.Unlock()

	player, ok := g.players[clientAddr.String()]
	if !ok {
		log.Printf("new player connecting: %s", clientAddr.String())
		randomPoint, err := g.area.GetRandomFreePoint()
		if err != nil {
			// can't get free space for new player
			// TODO: should tell the player that they can't connect
			return
		}
		player = &Player{
			Snake: entities.NewSnake(1, randomPoint.X, randomPoint.Y),
			Addr:  clientAddr,
		}
		g.players[clientAddr.String()] = player
	}

	decoder := gob.NewDecoder(data)
	var playerDirection entities.SnakeDirection
	if err := decoder.Decode(&playerDirection); err != nil {
		// shouldn't happen if happens, we just ignore
		return
	}
	log.Printf("setting player %s direction to %d", clientAddr.String(), playerDirection)

	player.Snake.SetDirection(playerDirection)
}
