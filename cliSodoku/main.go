package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/MarcGrol/sodoku/solver"
)

var (
	_Verbose      *bool
	_Timeout      *int
	_MinSolutions *int
)

func main() {
	processArgs()

	gameData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading game-data from stdin: %s", err)
		os.Exit(-2)
	}

	game, err := solver.LoadString(string(gameData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading game: %s", err)
		os.Exit(-2)
	}

	solutions, err := solver.Solve(game, *_Timeout, *_MinSolutions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error solving game: %s", err)
		game.DumpGameState()
		os.Exit(-3)
	}

	for _, solution := range solutions {
		fmt.Fprintf(os.Stdout, "%s\n", solution.Dump())
		fmt.Fprintf(os.Stderr, "steps:%d, guesses:%d\n\n",
			solution.CellsToBeSolved, solution.GuessCount)
	}

	os.Exit(0)
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	fmt.Fprintf(os.Stderr, " %s [flags]\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func processArgs() {
	help := flag.Bool("help", false, "Usage information")
	_Verbose = flag.Bool("verbose", false, "Verbose logging to stderr")
	_Timeout = flag.Int("timeout", 10, "Timeout in secs before give up")
	_MinSolutions = flag.Int("solutions", 1, "Nummber of solutions to wait for before give up")

	flag.Parse()

	if help != nil && *help == true {
		printUsage()
	}

	solver.Verbose = *_Verbose

	fmt.Fprintf(os.Stderr, "Using timeout %d and minSolutions: %d\n", *_Timeout, *_MinSolutions)
}
