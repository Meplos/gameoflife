package board

import (
	"log"
	"math/rand"
	"sync"
)

const P = 0.15

type Board struct {
	W        uint     `json:"w"`
	H        uint     `json:"h"`
	Cells    [][]bool `json:"state"`
	IsPaused bool     `json:"pause"`
	mu       sync.Mutex
}

func NewBoard(width, height uint) *Board {
	return &Board{
		W:        width,
		H:        height,
		Cells:    Array2D(width, height),
		IsPaused: true,
	}

}

func (b *Board) Randomize(percent float32) {
	if percent > 1 || percent < 0 {
		panic("board randomize error. initial population btw 0~1")
	}

	alive := 0
	for y := 0; y < int(b.H); y++ {
		for x := 0; x < int(b.W); x++ {
			v := rand.Float32() < percent
			b.Cells[y][x] = v
			if v {
				alive++
			}
		}
	}
	log.Printf("[Board.Randomize] p:%v alive:%d", percent, alive)
}
func Array2D(width, height uint) [][]bool {

	array := make([][]bool, height)
	for i := range array {
		array[i] = make([]bool, width)
	}
	return array
}

func (b *Board) Next() {
	if b.IsPaused {
		return
	}
	newState := Array2D(b.W, b.H)
	for y := 0; y < int(b.H); y++ {
		for x := 0; x < int(b.W); x++ {
			newState[y][x] = b.processCell(x, y)
		}
	}
	b.Cells = newState
}

func (b *Board) processCell(x, y int) bool {
	alive := b.Cells[y][x]
	s := b.countNeighbors(x, y)
	if s == 3 {
		return true
	}

	if alive && s == 2 {
		return true
	}

	return false
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

			if b.Cells[y+dy][x+dx] {
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
			if cell == true {
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
	b.Cells = Array2D(b.W, b.H)
	b.Randomize(P)

}
