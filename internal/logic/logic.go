package logic

import (
	"bytes"
	"fmt"

	"sync"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/engine"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/entities"
	"github.com/google/uuid"
)

type Player struct {
	ID    string
	Snake *entities.Snake
}

type Logic struct {
	area       *engine.Area
	updateLock sync.Mutex

	players map[string]Player
	bricks  []*entities.Brick
	foods   []*entities.Food
}

const (
	maxSquareAreaCols = 8191
	maxSquareAreaRows = 8191
	foodCount         = 70000
)

func NewLogic() *Logic {
	area := engine.NewArea(maxSquareAreaCols, maxSquareAreaRows)

	logic := &Logic{
		area:    area,
		bricks:  make([]*entities.Brick, 0, area.Cols*2+area.Rows*2-4),
		players: make(map[string]Player),
		foods:   make([]*entities.Food, 0, foodCount),
	}

	logic.addBricks()
	return logic
}

func (l *Logic) addBricks() {
	for x := range l.area.Cols {
		l.bricks = append(l.bricks, entities.NewBrick(x, 0))
		l.bricks = append(l.bricks, entities.NewBrick(x, l.area.Rows-1))
	}

	for y := range l.area.Rows - 2 {
		l.bricks = append(l.bricks, entities.NewBrick(0, y+1))
		l.bricks = append(l.bricks, entities.NewBrick(l.area.Cols-1, y+1))
	}
}

// addFood adds food until there is "foodCount" amount of food in the game.
// foodCount value directly correlates with startup time of game server.
func (l *Logic) addFood() {
	for n := len(l.foods); n < foodCount; n++ {
		randomPoint, err := l.area.GetRandomFreePoint()
		if err != nil {
			// no area for add food...
			// in practice game should end very soon, as the whole space is occupied
			// by snakes and there is a big probability of collision
			return
		}
		l.foods = append(l.foods, entities.NewFood(randomPoint.X, randomPoint.Y))
	}
}

func deleteFromSlice[T entities.Food | Player](s []*T, indexesToRemove map[int]struct{}) []*T {
	j := 0
	for i, val := range s {
		if _, ok := indexesToRemove[i]; !ok {
			s[j] = val
			j++
		}
	}
	// no need to clear past j, as we will be appending new values soon
	return s[:j]
}

func (l *Logic) UpdateState() {
	l.updateLock.Lock()
	defer l.updateLock.Unlock()

	for _, player := range l.players {
		player.Snake.Move()
	}

	// check collision with other snakes and itself
players:
	for pid, player := range l.players {
		for aid, anotherPlayer := range l.players {
			if pid == aid {
				if player.Snake.CollidesWithItself() {
					fmt.Printf("player %s collided with itself\n", pid)
					delete(l.players, pid)
					continue players
				}
			} else {
				if player.Snake.Collides(anotherPlayer.Snake) {
					fmt.Printf("player %s collided with player %s\n", pid, aid)
					delete(l.players, pid)
					continue players
				}
			}
		}
	}

	// check brick collision with any snake
	for _, b := range l.bricks {
		for pid, player := range l.players {
			snake := player.Snake
			if engine.Collides(snake.Head(), b) {
				delete(l.players, pid)
			}
		}
	}

	// check food collision with any snake
	eatenFoods := make(map[int]struct{}, len(l.players)) // at most "player" amount of foods can be theoritically eaten during one tick
	for i, food := range l.foods {
		for _, player := range l.players {
			snake := player.Snake

			if engine.Collides(snake.Head(), food) {
				snake.Grow()
				eatenFoods[i] = struct{}{}
			}
		}
	}

	if len(eatenFoods) != 0 {
		l.foods = deleteFromSlice(l.foods, eatenFoods)
	}

	l.area.Bounders = l.area.Bounders[:0]

	for _, b := range l.bricks {
		l.area.Bounders = append(l.area.Bounders, b)
	}

	for _, player := range l.players {
		l.area.Bounders = append(l.area.Bounders, player.Snake.Bounders...)
	}

	// regenerate food after adding players and bricks to area,
	// otherwise it can be generated on top of them
	l.addFood()

	for _, food := range l.foods {
		l.area.Bounders = append(l.area.Bounders, food)
	}
}

func (l *Logic) WriteStateForPlayer(id string, buff *bytes.Buffer) error {
	l.updateLock.Lock()
	defer l.updateLock.Unlock()

	player, ok := l.players[id]
	if !ok {
		return fmt.Errorf("player not found: %s", id)
	}

	snakeHead := player.Snake.Head()

	// the clone method recreates the bounders.
	// in this case case we might actually don't recreate just change
	// the reference to the slice itself, while original slice is kept with [Logic] struct.
	// this can be a good performance gainer, if we happen to have a lot of bounders.
	clonedArea := l.area.Clone()
	clonedArea.RemoveInvisibleCells(
		snakeHead.Point,
		engine.DefaultCameraOffsetCols,
		engine.DefaultCameraOffsetRows,
	)

	err := clonedArea.Encode(buff, snakeHead.Point)
	if err != nil {
		return fmt.Errorf("can't encode area: %w", err)
	}

	return nil
}

func (l *Logic) AddPlayer(direction entities.SnakeDirection) (string, error) {
	l.updateLock.Lock()
	defer l.updateLock.Unlock()

	freePoint, err := l.area.GetRandomFreePoint()
	if err != nil {
		return "", fmt.Errorf("can't add player: %w", err)
	}
	player := Player{
		Snake: entities.NewSnakeWithDirection(1, freePoint.X, freePoint.Y, direction),
	}
	id := uuid.NewString()
	l.players[id] = player
	return id, nil
}

func (l *Logic) SetPlayerDirection(id string, direction entities.SnakeDirection) error {
	l.updateLock.Lock()
	defer l.updateLock.Unlock()

	player, ok := l.players[id]
	if !ok {
		return fmt.Errorf("player not found: %s", id)
	}

	player.Snake.SetDirection(direction)
	return nil
}
