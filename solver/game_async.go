package solver

import (
	"reflect"
	"time"
)

func Solve(g *Game, timeout int, minSolutionCount int) ([]*Game, error) {
	// non-blocking channel to prevent go-routines to block each other on reporting solution
	solutionChannel := make(chan *Game, 1000)
	duration := time.Duration(timeout) * time.Second

	// Store completion variables within game
	g.solutionChannel = solutionChannel
	g.deadline = time.Now().Add(duration)

	// Start solving in background
	// Solutions will be reported back over solutionChannel
	go solve(g)

	// Wait for a solution
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
				debug("Solution is new:\n")
				solutions = append(solutions, newSolution)
				if len(solutions) >= minSolutionCount {
					debug("Enough solutions received: %d\n", len(solutions))
					break outerLoop
				}
			} else {
				debug("Solution exists")
			}
		case <-timer:
			debug("Timeout expired after %d secs\n", duration)
			break outerLoop
		}
	}

	if len(solutions) == 0 {
		debug("No solutions found")
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
