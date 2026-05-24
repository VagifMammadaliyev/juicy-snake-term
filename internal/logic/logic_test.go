package logic

import (
	"bytes"
	"testing"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/entities"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/messages"
)

func TestNewLogic_CapsArea(t *testing.T) {
	l := NewLogic(maxSquareAreaCols+100, maxSquareAreaRows+100)
	if l.area.Cols != maxSquareAreaCols {
		t.Errorf("cols: got %d, want %d", l.area.Cols, maxSquareAreaCols)
	}
	if l.area.Rows != maxSquareAreaRows {
		t.Errorf("rows: got %d, want %d", l.area.Rows, maxSquareAreaRows)
	}
}

func TestNewLogic_AddsBorderBricks(t *testing.T) {
	const cols, rows = int16(20), int16(15)
	l := NewLogic(cols, rows)
	want := int(cols*2 + (rows-2)*2)
	if got := len(l.bricks); got != want {
		t.Errorf("brick count: got %d, want %d", got, want)
	}
}

func TestDeleteFromSlice(t *testing.T) {
	foods := []*entities.Food{
		entities.NewFood(0, 0),
		entities.NewFood(1, 1),
		entities.NewFood(2, 2),
		entities.NewFood(3, 3),
	}
	got := deleteFromSlice(foods, map[int]struct{}{1: {}, 3: {}})

	if len(got) != 2 {
		t.Fatalf("len: got %d, want 2", len(got))
	}
	if got[0].X != 0 || got[1].X != 2 {
		t.Errorf("kept wrong items: x=%d, x=%d", got[0].X, got[1].X)
	}
}

func TestAddPlayer(t *testing.T) {
	l := NewLogic(20, 20)

	id, err := l.AddPlayer(entities.Up)
	if err != nil {
		t.Fatalf("AddPlayer: %v", err)
	}
	if _, ok := l.players[id]; !ok {
		t.Errorf("player %q missing from map", id)
	}
}

func TestSetPlayerDirection_UnknownPlayer(t *testing.T) {
	l := NewLogic(20, 20)
	if err := l.SetPlayerDirection("nope", entities.Up); err == nil {
		t.Error("expected error for unknown player")
	}
}

func TestWriteStateForPlayer_UnknownPlayer(t *testing.T) {
	l := NewLogic(20, 20)
	var buf bytes.Buffer
	if err := l.WriteStateForPlayer("nope", &buf); err == nil {
		t.Error("expected error for unknown player")
	}
}

func TestWriteStateForPlayer_WritesAreaUpdateHeader(t *testing.T) {
	l := NewLogic(20, 20)
	id, err := l.AddPlayer(entities.Left)
	if err != nil {
		t.Fatalf("AddPlayer: %v", err)
	}

	var buf bytes.Buffer
	if err := l.WriteStateForPlayer(id, &buf); err != nil {
		t.Fatalf("WriteStateForPlayer: %v", err)
	}
	if buf.Len() == 0 || messages.MsgType(buf.Bytes()[0]) != messages.MsgTypeAreaUpdate {
		t.Errorf("first byte: got %d, want %d", buf.Bytes()[0], messages.MsgTypeAreaUpdate)
	}
}

func TestUpdateState_SnakeMoves(t *testing.T) {
	l := NewLogic(20, 20)
	id, err := l.addPlayerInCoords(entities.Right, 5, 5)
	if err != nil {
		t.Fatalf("addPlayerInCoords: %v", err)
	}

	l.UpdateState()

	p, ok := l.players[id]
	if !ok {
		t.Fatal("player removed unexpectedly")
	}
	if head := p.Snake.Head(); head.X != 6 || head.Y != 5 {
		t.Errorf("head: got (%d,%d), want (6,5)", head.X, head.Y)
	}
}

func TestUpdateState_BrickCollisionRemovesPlayer(t *testing.T) {
	l := NewLogic(20, 20)

	id, err := l.addPlayerInCoords(entities.Up, 5, 1)
	if err != nil {
		t.Fatalf("addPlayerInCoords: %v", err)
	}

	l.UpdateState()

	if _, ok := l.players[id]; ok {
		t.Error("expected player to be removed after colliding with border brick")
	}
}

func TestUpdateState_SnakeRammingAnothersBodyIsRemoved(t *testing.T) {
	l := NewLogic(20, 20)

	p1, err := l.addPlayerInCoords(entities.Left, 5, 5)
	if err != nil {
		t.Fatalf("addPlayerInCoords p1: %v", err)
	}
	l.players[p1].Snake.Grow()

	p2, err := l.addPlayerInCoords(entities.Right, 4, 5)
	if err != nil {
		t.Fatalf("addPlayerInCoords p2: %v", err)
	}

	l.UpdateState()

	if _, ok := l.players[p2]; ok {
		t.Error("expected p2 to be removed after ramming p1's body")
	}
	if _, ok := l.players[p1]; !ok {
		t.Error("p1 should survive — its head moved away from the collision cell")
	}
}

func TestUpdateState_EatingFoodGrowsSnake(t *testing.T) {
	l := NewLogic(20, 20)
	id, err := l.addPlayerInCoords(entities.Right, 5, 5)
	if err != nil {
		t.Fatalf("addPlayerInCoords: %v", err)
	}
	l.foods = append(l.foods, entities.NewFood(6, 5))

	before := len(l.players[id].Snake.Bounders)
	l.UpdateState()
	after := len(l.players[id].Snake.Bounders)

	if after != before+1 {
		t.Errorf("snake length: got %d, want %d", after, before+1)
	}
}
