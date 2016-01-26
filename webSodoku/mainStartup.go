// +build !appengine

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/MarcGrol/sodoku/solver"
	"github.com/MarcGrol/sodoku/web"
	"github.com/justinas/alice"
)

var (
	_Timeout      *int
	_MinSolutions *int
	_Verbose      *bool
	_Port         *int
)

func main() {
	processArgs()

	solver.Verbose = *_Verbose

	h := &web.SodokuHandler{
		Timeout:      *_Timeout,
		MinSolutions: *_MinSolutions,
	}
	// configure middleware around our example
	chain := alice.New(loggingHandler, countingHandler).Then(h)

	mux := http.NewServeMux()
	mux.Handle("/sodoku", chain)

	fmt.Printf("Start listening at http://localhost:%d/\n", *_Port)
	http.ListenAndServe(fmt.Sprintf(":%d", *_Port), chain)
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
	_Timeout = flag.Int("timeout", 2, "Timeout in secs before give up")
	_MinSolutions = flag.Int("solutions", 1, "Nummber of solutions to wait for before give up")
	_Port = flag.Int("port", 3000, "Port web-server is listening at")

	flag.Parse()

	if help != nil && *help == true {
		printUsage()
	}

	fmt.Fprintf(os.Stderr, "Using timeout %d and solutions: %d\n", *_Timeout, *_MinSolutions)
}
