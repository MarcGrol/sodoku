package solver

import (
	"reflect"
	"time"

	"github.com/MarcGrol/sodoku/workerpool"
)

const (
	WORKER_POOL_SIZE = 100
)

var (
	pool *workerpool.WorkerPool
)

func solveInBackground(g *Game) {
	pool.Execute(wrappedSolve, g, g.timeoutSec)
}

func Solve(g *Game, timeout int, minSolutionCount int) ([]*Game, error) {
	pool = workerpool.NewWorkerPool(WORKER_POOL_SIZE)
	pool.Start()

	// non-blocking channel to prevent go-routines to block each other on reporting solution
	solutionChannel := make(chan *Game, 1000)
	duration := time.Duration(timeout) * time.Second

	// Store completion variables within game
	g.solutionChannel = solutionChannel
	g.deadline = time.Now().Add(duration)
	g.timeoutSec = duration

	// Start solving in background
	// Solutions will be reported back over solutionChannel
	solveInBackground(g)

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

	pool.Stop()

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
