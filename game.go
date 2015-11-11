package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

const square_size = 9
const section_size = 3

type Game struct {
	square          *Square
	CellsSolved     int
	GuessCount      int
	SolutionChannel chan *Game
	DeadLine        time.Time
}

func newGame() *Game {
	g := Game{}
	g.square = NewSquare(square_size)
	return &g
}

func (g Game) copy() *Game {
	ng := Game{}
	ng.square = g.square.Copy()
	ng.CellsSolved = g.CellsSolved
	ng.GuessCount = g.GuessCount
	ng.SolutionChannel = g.SolutionChannel
	ng.DeadLine = g.DeadLine
	return &ng
}

func Load(lines string) (*Game, error) {
	linesRead := 0
	game := newGame()
	for x, line := range strings.Split(lines, "\n") {

		if x >= square_size {
			break
		}
		if line == "" {
			break
		}

		splitted := strings.Split(line, " ")
		if len(splitted) != square_size {
			return nil, fmt.Errorf("Invalid number of columns for row %d: needs %d, actual %d", x+1,
				square_size, len(splitted))
		}
		for y, val := range splitted {
			if val == "_" {
				game.square.Clear(x, y)
				continue
			}
			num, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("Invalid number '%s' for item row:%d, column:%d", val, x+1, y+1)
			}
			if num < 0 || num > square_size {
				return nil, fmt.Errorf("Invalid value %d for item row:%d, column:%d", num, x+1, y+1)

			}
			game.square.Set(x, y, num)
		}
		linesRead++
	}
	if linesRead != square_size {
		return nil, fmt.Errorf("Not enough rows: needs %d, actual %d", square_size, linesRead)
	}

	err := game.validate()
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (g Game) validate() error {
	// check for duplicate values in each row
	for x := 0; x < g.square.Size; x++ {
		row := g.square.GetRowValues(x)
		hasDups, dupVal := containsDuplicates(row)
		if hasDups {
			return fmt.Errorf("Row %d contains duplicate value %d", x+1, dupVal)
		}
	}

	// check for duplicate values in each column
	for y := 0; y < g.square.Size; y++ {
		column := g.square.GetColumnValues(y)
		hasDups, dupVal := containsDuplicates(column)
		if hasDups {
			return fmt.Errorf("Column %d contains duplicate value %d", y+1, dupVal)
		}
	}

	// check for duplicate values in each section
	// only visit the centers of each section
	for x := 1; x < g.square.Size; x += section_size {
		for y := 1; y < g.square.Size; y += section_size {
			section := g.square.GetSectionValues(x, y)
			hasDups, dupVal := containsDuplicates(section)
			if hasDups {
				return fmt.Errorf("Cell %d-%d is in section with duplicate value %d", x+1, y+1, dupVal)
			}
		}
	}

	return nil
}

func containsDuplicates(array []int) (bool, int) {
	for idx, i := range array {
		if contains(array, i, idx) {
			return true, i
		}
	}
	return false, -1
}

func contains(array []int, value int, selfIdx int) bool {
	for idx, i := range array {
		if idx == selfIdx {
			continue
		}
		if value == i {
			return true
		}
	}
	return false
}

func Solve(g *Game, timeout int, minSolutionCount int) ([]*Game, error) {
	// blocking to prevent go-routines to block on reporting solution
	solutionChannel := make(chan *Game, 1000)
	duration := time.Duration(timeout) * time.Second

	// package complrtion params within game
	g.SolutionChannel = solutionChannel
	g.DeadLine = time.Now().Add(duration)

	// Start solving in background
	// Solutions will be reported back over solutionChannel
	go solve(g)

	return waitforCompletion(solutionChannel, duration, minSolutionCount)

}

func waitforCompletion(solutionChannel chan *Game, duration time.Duration, minSolutionCount int) ([]*Game, error) {
	timer := time.After(duration)

	solutions := make([]*Game, 0, 10)
outerLoop:
	for {
		select {
		case newSolution := <-solutionChannel:
			if !solutionExists(solutions, newSolution) {
				fmt.Fprintf(os.Stderr, "Solution is new:\n%s", newSolution)
				solutions = append(solutions, newSolution)
				if len(solutions) >= minSolutionCount {
					fmt.Fprintf(os.Stdout, "Enough solutions received: %d\n", len(solutions))
					break outerLoop
				}
			} else {
				fmt.Fprintf(os.Stderr, "Solution exists")
			}
		case <-timer:
			fmt.Fprintf(os.Stdout, "Timeout expired after %d secs\n", timeout)
			break outerLoop
		}
	}

	if len(solutions) == 0 {
		return solutions, fmt.Errorf("No solutions found")
	}
	return solutions, nil
}

func solutionExists(solutions []*Game, newSolution *Game) bool {
	for _, s := range solutions {
		if reflect.DeepEqual(s.square, newSolution.square) {
			return true
		}
	}

	return false
}

func solve(g *Game) {
	maxSteps := square_size * square_size

	for i := 0; i < maxSteps; i++ {

		if time.Now().After(g.DeadLine) {
			fmt.Fprintf(os.Stderr, "%p: Abort because deadline expired\n", g)
			return
		}

		cellsSolvedInStep := g.step()

		fmt.Fprintf(os.Stderr, "%p: Solved %d cells this loop\n", g, cellsSolvedInStep)

		if cellsSolvedInStep < 0 {
			// wrong guess upstream, terminate go-routine
			return
		}

		if cellsSolvedInStep == 0 {
			// stuck using deterministic approach: start guessing
			guessAndContinue(g)
			return
		}

		if g.countEmptyValues() == 0 && g.validate() == nil {
			fmt.Fprintf(os.Stderr, "%p: Got solution\n", g)
			// we are done: report result back over solution-channel
			g.SolutionChannel <- g
			return
		}

		g.CellsSolved += cellsSolvedInStep
	}

	// unsolveable
	fmt.Fprintf(os.Stderr, "%p: Abort after cells solved:%d\n", g, g.CellsSolved)
}

func (g *Game) step() int {
	cellsSolved := 0

	for x := 0; x < g.square.Size; x++ {
		for y := 0; y < g.square.Size; y++ {
			if !g.square.Has(x, y) {
				mergedCandidates := g.findCandidates(x, y)
				if len(mergedCandidates) == 0 {
					// we have mad a wrong guess somwhere
					fmt.Fprintf(os.Stderr, "%p: Cell %d-%d has zero candidates due to wrong guess upstream\n", g, x+1, y+1)
					return -1
				} else if len(mergedCandidates) == 1 {
					g.square.Set(x, y, mergedCandidates[0])
					cellsSolved++
				}
			}
		}
	}

	return cellsSolved
}

func guessAndContinue(g *Game) {
	orderedBestGuesses := g.findCellsWithLeastCandidates()

	if len(orderedBestGuesses) > 0 {
		bestGuess := orderedBestGuesses[0]
		for _, cand := range bestGuess.candidates {
			fmt.Fprintf(os.Stderr, "%p: Try %d-%d to %d and continue\n", g, bestGuess.x+1, bestGuess.y+1, cand)
			cpy := g.copy()
			cpy.square.Set(bestGuess.x, bestGuess.y, cand)
			g.GuessCount++
			go solve(cpy)
		}
	}
}

func (g *Game) findCellsWithLeastCandidates() []cell {
	cells := make([]cell, 0, square_size*square_size)
	g.square.Iterate(func(x int, y int, z int) error {
		if !g.square.Has(x, y) {
			cellCandidates := g.findCandidates(x, y)
			if len(cellCandidates) > 1 {
				cells = append(cells, cell{x: x, y: y, candidates: cellCandidates})
			}
		}
		return nil
	})
	sort.Sort(CellByNumberOfCandidates(cells))

	return cells
}

func (g *Game) findCandidates(x int, y int) []int {
	mergedValues := mergeValues(
		g.square.GetRowValues(x),
		g.square.GetColumnValues(y),
		g.square.GetSectionValues(x, y))
	sort.Ints(mergedValues)
	return findCandidates(mergedValues)
}

func (g *Game) DumpGameState() {

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
				mergedCandidates := g.findCandidates(x, y)
				alternatives := fmt.Sprintf("%v", mergedCandidates)
				fmt.Fprintf(os.Stderr, "%-12s", alternatives)

			}
		}
		fmt.Fprintf(os.Stderr, "|\n")
	}
	fmt.Fprintf(os.Stderr, "___________________________________________________________________________________________________________________\n")
	fmt.Fprintf(os.Stderr, "\n\n")
}

func (g *Game) countEmptyValues() int {
	count := 0
	g.square.Iterate(func(x, y, z int) error {
		if !g.square.Has(x, y) {
			count++
		}
		return nil
	})
	return count
}

func mergeValues(rowValues []int, columnValues []int, sectionValues []int) []int {
	merged := append(rowValues, columnValues...)
	return append(merged, sectionValues...)
}

func findCandidates(existing []int) []int {
	full := makeFull(square_size)
	candidates := minus(full, existing)

	return candidates
}

func makeFull(size int) []int {
	full := make([]int, size)
	for i := 0; i < size; i++ {
		full[i] = (i + 1)
	}
	return full
}

func minus(all []int, other []int) []int {
	stripped := make([]int, 0, len(all))
	for _, a := range all {
		if !contains(other, a, -1) {
			stripped = append(stripped, a)
		}
	}
	return stripped
}

type cell struct {
	x          int
	y          int
	candidates []int
}
type CellByNumberOfCandidates []cell

func (arr CellByNumberOfCandidates) Len() int      { return len(arr) }
func (arr CellByNumberOfCandidates) Swap(i, j int) { arr[i], arr[j] = arr[j], arr[i] }
func (arr CellByNumberOfCandidates) Less(i, j int) bool {
	return len(arr[i].candidates) < len(arr[j].candidates)
}

func (g Game) String() string {
	return g.square.String()
}
