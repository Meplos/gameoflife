package gol

import (
	"log"
	"time"

	"github.com/Meplos/GameOfLife/board"
	"github.com/google/uuid"
)

const FPS = 30

type Game struct {
	ID uuid.UUID
	B  *board.Board
}

func NewGame(w, h uint) *Game {
	return &Game{
		ID: uuid.New(),
		B:  board.NewBoard(w, h),
	}
}

func (g *Game) Init() {
	g.B.Randomize(0.15)
	g.B.IsPaused = true
}

func (g *Game) Start(pipe chan board.BoardState) {
	g.B.Play()
	g.Run(pipe)

}

func (g *Game) Run(pipe chan board.BoardState) {
	ticker := time.NewTicker(time.Second / FPS)
	for t := range ticker.C {
		if g.B.IsPaused {
			return
		}
		log.Printf("Tick %v, Alive: %v Running:%v\n", t, g.B.AliveCount(), !g.B.IsPaused)
		g.B.Next()
		pipe <- g.B.ToBoardState()
	}
}

func (g *Game) Restart() {
	g.Init()
}

func (g *Game) Stop() {
	g.B.Pause()
}
