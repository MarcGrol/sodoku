package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/MarcGrol/sodoku/solver"
	"github.com/justinas/alice"
)

const TIMEOUT = 2
const MAX_SOLUTIONS = 1
const HARD_EXAMPLE = `9 _ _ _ 2 _ _ _ 5
_ _ _ 9 _ 5 _ _ _
_ _ 7 _ 6 _ 4 _ _
_ 5 _ _ _ _ _ 7 _
8 _ 1 _ 7 _ 6 _ 2
_ 2 _ _ _ _ _ 3 _
_ _ 6 _ 4 _ 1 _ _
_ _ _ 8 _ 3 _ _ _
4 _ _ _ 9 _ _ _ 3
`

var (
	_Timeout      *int
	_MinSolutions *int
	_Verbose      *bool
	_Port         *int
)

type sodokuHandler struct {
}

func (eh *sodokuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		eh.get(w, r)
	case "POST":
		eh.post(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("Only GET and POST are supported"))
	}
}

func (eh *sodokuHandler) get(w http.ResponseWriter, r *http.Request) {
	game, err := solver.LoadString(HARD_EXAMPLE)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
	}
	doSolve(w, r, game)
}

func (eh *sodokuHandler) post(w http.ResponseWriter, r *http.Request) {
	// read incoming steps
	input, err := FromJson(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
	}

	// load
	game, err := solver.LoadSteps(toCoreSteps(input.Steps))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
	}

	// solve and return response
	doSolve(w, r, game)
}

func doSolve(w http.ResponseWriter, r *http.Request, game *solver.Game) {
	coreSolutions, err := solver.Solve(game, TIMEOUT, MAX_SOLUTIONS)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}

	resp := Response{Error: nil}
	for _, coreSolution := range coreSolutions {
		solution := Game{Steps: fromCoreSteps(coreSolution.Steps)}
		resp.Solutions = append(resp.Solutions, solution)
	}
	writeSuccess(w, resp)
}

func toCoreSteps(webSteps []Step) []solver.Step {
	steps := make([]solver.Step, 0, 100)
	for _, step := range webSteps {
		steps = append(steps, solver.Step{X: step.X, Y: step.Y, Z: solver.Value(step.Z), Initial: step.Initial, IsGuess: step.IsGuess})
	}
	return steps
}

func fromCoreSteps(coreSteps []solver.Step) []Step {
	steps := make([]Step, 0, 100)
	for _, step := range coreSteps {
		steps = append(steps, Step{X: step.X, Y: step.Y, Z: int(step.Z), Initial: step.Initial, IsGuess: step.IsGuess})
	}
	return steps
}

func writeError(w http.ResponseWriter, status int, err error) {
	resp := Response{Error: &ErrorDescriptor{Message: err.Error()}}
	// write headers
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	resp.ToJson(w)
}

func writeSuccess(w http.ResponseWriter, resp Response) {
	// write headers
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp.ToJson(w)
}

func main() {
	processArgs()

	solver.Verbose = *_Verbose

	h := &sodokuHandler{}
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
	_Timeout = flag.Int("timeout", 10, "Timeout in secs before give up")
	_MinSolutions = flag.Int("solutions", 1, "Nummber of solutions to wait for before give up")
	_Port = flag.Int("port", 3000, "Port web-server is listening at")

	flag.Parse()

	if help != nil && *help == true {
		printUsage()
	}

	fmt.Fprintf(os.Stderr, "Using timeout %d and minSolutions: %d\n", *_Timeout, *_MinSolutions)
}
