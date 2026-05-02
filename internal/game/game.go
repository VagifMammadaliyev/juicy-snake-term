package game

import (
	"bufio"
	"math/rand"
	"time"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/engine"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/entities"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
)

type Game struct {
	screen   *bufio.Writer
	area     *engine.Area
	Controls chan terminal.Control

	snake  *entities.Snake
	bricks []*entities.Brick
	food   *entities.Food
}

func NewGame(buf *bufio.Writer) *Game {
	area := engine.NewArea(40, 40)
	bricks := make([]*entities.Brick, 0, area.Cols*2+area.Rows*2-4)

	game := &Game{
		screen:   buf,
		area:     area,
		snake:    entities.NewSnake(2, area.Cols/2, area.Rows/2),
		bricks:   bricks,
		Controls: make(chan terminal.Control, 2),
		// food:     entities.NewFood(),
	}

	game.addBricks()
	return game
}

func (g *Game) addBricks() {
	for x := range g.area.Cols {
		g.bricks = append(g.bricks, entities.NewBrick(x, 0))
		g.bricks = append(g.bricks, entities.NewBrick(x, g.area.Rows-1))
	}

	for y := range g.area.Rows - 2 {
		g.bricks = append(g.bricks, entities.NewBrick(0, y+1))
		g.bricks = append(g.bricks, entities.NewBrick(g.area.Cols-1, y+1))
	}
}

func (g *Game) addFood() {
	if g.food == nil {
		points := g.area.CalculateFreePoints()
		randomPoint := points[rand.Intn(len(points))]
		g.food = entities.NewFood(randomPoint.X, randomPoint.Y)
	}
}

func (g *Game) Update() bool {
	g.snake.Move()
	if g.food != nil {
		shouldGrow := g.area.Collides(g.snake.Nodes[0], g.food)
		if shouldGrow {
			g.snake.Grow()
			g.food = nil
		}
	}

	for _, b := range g.bricks {
		if g.area.Collides(g.snake.Nodes[0], b) {
			return true
		}
	}

	g.area.Bounders = g.area.Bounders[:0]
	for _, b := range g.bricks {
		g.area.Bounders = append(g.area.Bounders, b)
	}
	for _, n := range g.snake.Nodes {
		g.area.Bounders = append(g.area.Bounders, n)
	}

	g.addFood()
	g.area.Bounders = append(g.area.Bounders, g.food)

	g.area.Render(g.screen)
	g.screen.Flush()
	return false
}

func (g *Game) Run() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			quit := g.Update()
			if quit {
				return
			}

		case control := <-g.Controls:
			if control == terminal.Quit {
				return
			}

			g.snake.SetDirection(entities.SnakeDirection(control))
		}
	}
}
