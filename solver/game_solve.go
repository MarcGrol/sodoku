package solver

import (
	"sort"
	"time"
)

func wrappedSolve(data interface{}) {
	solve(data.(*Game))
}

func solve(g *Game) {
	maxSteps := SQUARE_SIZE * SQUARE_SIZE

	debug("%p: Start solving\n", g)
	for i := 0; i < maxSteps; i++ {

		// This is a way to cleanly terminate running goroutines when game completed (due to timer or completion)
		if time.Now().After(g.deadline) {
			debug("%p: Abort because deadline expired\n", g)
			return
		}

		// check if cells can be easily solved
		cellsSolvedInStep := step(g)
		if cellsSolvedInStep < 0 {
			// wrong guess upstream, terminate go-routine
			return

		} else if cellsSolvedInStep == 0 {
			// stuck using deterministic approach: start guessing
			guessAndContinue(g)
			return

		} else {
			// this step resulted in cells-solved

			debug("%p: Solved %d cells this loop\n", g, cellsSolvedInStep)

			if g.countEmptyValues() == 0 {
				// we are done: report result back over solution-channel
				debug("%p: Got solution\n", g)
				g.solutionChannel <- g
				return
			}
		}
	}

	// unsolveable
	warning("%p: Abort at cells to go:%d\n", g, g.countEmptyValues())
}

func step(g *Game) int {
	cellsSolved := 0

	for x := 0; x < g.square.Size; x++ {
		for y := 0; y < g.square.Size; y++ {
			if !g.square.Has(x, y) {
				mergedCandidates := findCellCandidates(g, x, y)
				if len(mergedCandidates) == 0 {
					// we have mad a wrong guess somwhere
					debug("%p: Cell %d-%d has zero candidates due to wrong guess upstream\n", g, x+1, y+1)
					return -1
				} else if len(mergedCandidates) == 1 {
					g.set(x, y, mergedCandidates[0], false, false)
					cellsSolved++
				}
			}
		}
	}

	return cellsSolved
}

func guessAndContinue(g *Game) {
	orderedBestGuesses := findCellsWithLeastCandidates(g)

	if len(orderedBestGuesses) > 0 {
		bestGuess := orderedBestGuesses[0]
		for _, cand := range bestGuess.candidates {
			cpy := g.copy()
			debug("%p -> %p: Got stuck -> Try %d-%d with value %d and continue\n",
				g, cpy, bestGuess.x+1, bestGuess.y+1, cand)
			cpy.set(bestGuess.x, bestGuess.y, cand, false, true)
			solveInBackground(cpy)
		}
	} else {
		// unsolveable
		warning("%p: No best guesses found at:%d\n", g, g.countEmptyValues())
	}
}

func findCellsWithLeastCandidates(g *Game) []cell {
	cells := make([]cell, 0, SQUARE_SIZE*SQUARE_SIZE)
	g.square.Iterate(func(x int, y int, z Value) error {
		if !g.square.Has(x, y) {
			cellCandidates := findCellCandidates(g, x, y)
			if len(cellCandidates) > 1 {
				cells = append(cells, cell{x: x, y: y, candidates: cellCandidates})
			}
		}
		return nil
	})
	// cells with least number of candidate-values on top
	sort.Sort(CellByNumberOfCandidates(cells))

	return cells
}

func findCellCandidates(g *Game, x int, y int) []Value {
	mergedValues := mergeValues(
		g.square.GetRowValues(x),
		g.square.GetColumnValues(y),
		g.square.GetSectionValues(x, y))
	return findMissing(mergedValues)
}

func mergeValues(rowValues ValueSet, columnValues ValueSet, sectionValues ValueSet) ValueSet {
	merged := rowValues.Union(columnValues)
	return merged.Union(sectionValues)
}

func findMissing(existing ValueSet) []Value {
	full := NewValueSet(1, 2, 3, 4, 5, 6, 7, 8, 9)
	return full.Difference(existing).ToSlice()
}

type cell struct {
	x          int
	y          int
	candidates []Value
}
type CellByNumberOfCandidates []cell

func (arr CellByNumberOfCandidates) Len() int      { return len(arr) }
func (arr CellByNumberOfCandidates) Swap(i, j int) { arr[i], arr[j] = arr[j], arr[i] }
func (arr CellByNumberOfCandidates) Less(i, j int) bool {
	return len(arr[i].candidates) < len(arr[j].candidates)
}
