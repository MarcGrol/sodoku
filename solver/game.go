package solver

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	SQUARE_SIZE  = 9
	SECTION_SIZE = 3
)

var (
	mutex sync.Mutex
)

type Game struct {
	CellsToBeSolved int
	GuessCount      int
	Steps           []Step
	square          *Square
	solutionChannel chan *Game
	deadline        time.Time
	timeoutSecs     int
}

type Step struct {
	X       int
	Y       int
	Z       Value
	Initial bool
	IsGuess bool
}

func newGame() *Game {
	g := Game{}
	g.square = NewSquare(SQUARE_SIZE)
	return &g
}

func (g Game) copy() *Game {
	ng := Game{}
	ng.square = g.square.Copy()
	ng.CellsToBeSolved = g.CellsToBeSolved
	ng.GuessCount = g.GuessCount
	ng.solutionChannel = g.solutionChannel
	ng.deadline = g.deadline
	toGo := ng.deadline.Second() - time.Now().Second()
	ng.timeoutSecs = toGo
	for _, s := range g.Steps {
		ng.Steps = append(ng.Steps, Step{X: s.X, Y: s.Y, Z: s.Z})
	}
	return &ng
}

func (g *Game) set(x int, y int, z Value, initial bool, isGuess bool) {
	g.square.Set(x, y, z)
	g.Steps = append(g.Steps, Step{X: x, Y: y, Z: z, Initial: initial, IsGuess: isGuess})
	if isGuess == true {
		g.GuessCount++
	}
}

func (g *Game) countEmptyValues() int {
	count := 0
	g.square.Iterate(func(x int, y int, z Value) error {
		if !g.square.Has(x, y) {
			count++
		}
		return nil
	})
	return count
}

func (g Game) Dump() string {
	return g.square.String()
}

func (g Game) String() string {

	return g.square.String()
}

func (g *Game) DumpGameState() {
	mutex.Lock()
	defer mutex.Unlock()

	for x := 0; x < g.square.Size; x++ {
		if (x % 3) == 0 {
			fmt.Fprintf(os.Stderr, "___________________________________________________________________________________________________________________\n")
		}
		for y := 0; y < g.square.Size; y++ {
			if (y % 3) == 0 {
				fmt.Fprintf(os.Stderr, "| ")
			}
			if g.square.Has(x, y) {
				fmt.Fprintf(os.Stderr, "%-12d", g.square.Get(x, y))
			} else {
				mergedCandidates := findCellCandidates(g, x, y)
				alternatives := fmt.Sprintf("%v", mergedCandidates)
				fmt.Fprintf(os.Stderr, "%-12s", alternatives)

			}
		}
		fmt.Fprintf(os.Stderr, "|\n")
	}
	fmt.Fprintf(os.Stderr, "___________________________________________________________________________________________________________________\n")
	fmt.Fprintf(os.Stderr, "\n\n")
}
