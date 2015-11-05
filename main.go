package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var timeout *int

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

	solutions, err := Solve(game, *timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error solving game: %s", err)
		game.DumpGameState()
		os.Exit(-3)
	}
	for _, solution := range solutions {
		fmt.Fprintf(os.Stdout, "%s\n", solution)
		fmt.Fprintf(os.Stdout, "steps:%d, guesses:%d\n\n",
			solution.StepCount, solution.GuessCount)
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
	timeout = flag.Int("timeout", 60, "Timeout in secs before give up")

	flag.Parse()

	if help != nil && *help == true {
		printUsage()
	}
}
