package board

import (
	"log"
	"math/rand"
	"sync"
)

const P = 0.15

type CellState = bool

const (
	ALIVE CellState = true
	DEAD  CellState = false
)

type Change struct {
	X     int       `json:"x"`
	Y     int       `json:"y"`
	State CellState `json:"state"`
}

type BoardState struct {
	W        uint     `json:"w"`
	H        uint     `json:"h"`
	IsPaused bool     `json:"pause"`
	Changes  []Change `json:"changes"`
	Type     string   `json:"type"`
}

type Board struct {
	W        uint          `json:"w"`
	H        uint          `json:"h"`
	Cells    [][]CellState `json:"state"`
	IsPaused bool          `json:"pause"`
	mu       sync.Mutex
	Changes  []Change `json:"changes"`
}

type InitialBoard struct {
	W        uint          `json:"w"`
	H        uint          `json:"h"`
	Cells    [][]CellState `json:"state"`
	IsPaused bool          `json:"pause"`
	Type     string        `json:"type"`
}

func NewBoard(width, height uint) *Board {
	log.Printf("[Board.NewBoard] H: %v W:%v", height, width)
	return &Board{
		W:        width,
		H:        height,
		Cells:    Array2D[CellState](width, height),
		IsPaused: true,
		Changes:  make([]Change, 0),
	}

}

func (b *Board) Randomize(percent float32) {
	if percent > 1 || percent < 0 {
		panic("board randomize error. initial population btw 0~1")
	}

	alive := 0
	for y := 0; y < int(b.H); y++ {
		for x := 0; x < int(b.W); x++ {
			v := DEAD
			i := rand.Float32()
			log.Printf("[Board.Randomize] p:%v i:%v", percent, i)
			if i < percent {
				v = ALIVE
				alive++
			}
			b.Cells[y][x] = v
		}
	}
	log.Printf("[Board.Randomize] p:%v alive:%d", percent, alive)
}
func Array2D[T any](width, height uint) [][]T {

	array := make([][]T, height)
	for i := range array {
		array[i] = make([]T, width)
	}
	return array
}

func (b *Board) Next() {
	b.Changes = make([]Change, 0)
	if b.IsPaused {
		return
	}
	newState := Array2D[CellState](b.W, b.H)
	for y := 0; y < int(b.H); y++ {
		for x := 0; x < int(b.W); x++ {
			newState[y][x] = b.processCell(x, y)
			if newState[y][x] != b.Cells[y][x] {
				b.Changes = append(b.Changes, Change{
					X:     x,
					Y:     y,
					State: newState[y][x],
				})

			}
		}
	}
	b.Cells = newState
}

func (b *Board) processCell(x, y int) CellState {
	e := b.Cells[y][x]
	s := b.countNeighbors(x, y)
	if s == 3 {
		return ALIVE
	}

	if e == ALIVE && s == 2 {
		return ALIVE
	}

	return DEAD
}

func (b *Board) countNeighbors(x, y int) int {

	var count int
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx := x + dx
			ny := y + dy
			if nx < 0 || nx >= int(b.W) || ny < 0 || ny >= int(b.H) {
				continue
			}

			if b.Cells[y+dy][x+dx] == ALIVE {
				count++
			}
		}
	}
	return count
}

func (b *Board) AliveCount() int {

	n := 0

	for _, row := range b.Cells {
		for _, cell := range row {
			if cell == ALIVE {
				n++
			}
		}
	}
	return n

}

func (b *Board) Pause() {
	b.mu.Lock()
	b.IsPaused = true
	b.mu.Unlock()
}

func (b *Board) Play() {
	b.mu.Lock()
	b.IsPaused = false
	b.mu.Unlock()
}

func (b *Board) Restart() {
	b.Pause()
	b.Cells = Array2D[CellState](b.W, b.H)
	b.Randomize(P)
}

func (b *Board) ToBoardState() BoardState {
	return BoardState{
		W:        b.W,
		H:        b.H,
		Changes:  b.Changes,
		IsPaused: b.IsPaused,
		Type:     "change",
	}
}

func (b *Board) ToInitialState() InitialBoard {
	return InitialBoard{
		W:        b.W,
		H:        b.H,
		Cells:    b.Cells,
		IsPaused: b.IsPaused,
		Type:     "init",
	}
}
