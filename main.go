package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	timeout      *int
	minSolutions *int
)

func main() {
	processArgs()

	gameData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading game-data from stdin: %s", err)
		os.Exit(-2)
	}

	game, err := Load(string(gameData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading game: %s", err)
		os.Exit(-2)
	}

	solutions, err := Solve(game, *timeout, *minSolutions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error solving game: %s", err)
		game.DumpGameState()
		os.Exit(-3)
	}

	for _, solution := range solutions {
		fmt.Fprintf(os.Stderr, "%s\n", solution)
		fmt.Fprintf(os.Stderr, "steps:%d, guesses:%d\n\n",
			solution.CellsSolved, solution.GuessCount)
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
	timeout = flag.Int("timeout", 10, "Timeout in secs before give up")
	minSolutions = flag.Int("solutions", 1, "Nummber of solutions to wait for before give up")

	flag.Parse()

	if help != nil && *help == true {
		printUsage()
	}

	fmt.Fprintf(os.Stderr, "Using timeout %d and minSolutions: %d\n", *timeout, *minSolutions)
}
