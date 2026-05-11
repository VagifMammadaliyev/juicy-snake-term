package logic

import (
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

func NewLogic() *Logic {
	area := engine.NewArea(31, 31)

	logic := &Logic{
		area:    area,
		bricks:  make([]*entities.Brick, 0, area.Cols*2+area.Rows*2-4),
		players: make(map[string]Player),
		foods:   make([]*entities.Food, 0, 10),
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

func (l *Logic) addFood() {
	if len(l.foods) < 1 {
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
	for pid, player := range l.players {
		for _, anotherPlayer := range l.players {
			for i, node := range anotherPlayer.Snake.Nodes {
				if i == 0 && player == anotherPlayer {
					continue // skip head of the same snake
				}

				if l.area.Collides(player.Snake.Head(), node) {
					delete(l.players, pid)
				}

			}
		}
	}

	// check brick collision with any snake
	for _, b := range l.bricks {
		for pid, player := range l.players {
			snake := player.Snake
			if l.area.Collides(snake.Head(), b) {
				delete(l.players, pid)
			}
		}
	}

	// check food collision with any snake
	eatenFoods := make(map[int]struct{}, len(l.players)) // at most "player" amount of foods can be theoritically eaten during one tick
	for i, food := range l.foods {
		for _, player := range l.players {
			snake := player.Snake

			if l.area.Collides(snake.Head(), food) {
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
		s := player.Snake

		for _, n := range s.Nodes {
			l.area.Bounders = append(l.area.Bounders, n)
		}
	}

	l.addFood()
	for _, food := range l.foods {
		l.area.Bounders = append(l.area.Bounders, food)
	}
}

func (l *Logic) GetStateForPlayer(id string) (*engine.EncodedArea, error) {
	l.updateLock.Lock()
	defer l.updateLock.Unlock()

	player, ok := l.players[id]
	if !ok {
		return nil, fmt.Errorf("player not found: %s", id)

	}

	snakeHead := player.Snake.Head()

	encodedArea := l.area.ToEncodedArea()
	encodedArea.PlayerSnakeHead = snakeHead.Point
	encodedArea.RemoveInvisibleCells(snakeHead.Point, 21)

	return encodedArea, nil
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
