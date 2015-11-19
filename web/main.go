package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/MarcGrol/sodoku/core"
	"github.com/justinas/alice"
)

const TIMEOUT = 2
const MAX_SOLUTIONS = 1

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

func (eh *sodokuHandler) post(w http.ResponseWriter, r *http.Request) {
	exercise, err := FromJson(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
	}

	game, err := core.LoadSteps(toCoreSteps(exercise.Steps))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
	}

	doSolve(w, r, game)
}

func (eh *sodokuHandler) get(w http.ResponseWriter, r *http.Request) {
	game, err := core.Load(`9 _ _ _ 2 _ _ _ 5
_ _ _ 9 _ 5 _ _ _
_ _ 7 _ 6 _ 4 _ _
_ 5 _ _ _ _ _ 7 _
8 _ 1 _ 7 _ 6 _ 2
_ 2 _ _ _ _ _ 3 _
_ _ 6 _ 4 _ 1 _ _
_ _ _ 8 _ 3 _ _ _
4 _ _ _ 9 _ _ _ 3
`)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
	}
	doSolve(w, r, game)
}

func doSolve(w http.ResponseWriter, r *http.Request, game *core.Game) {
	coreSolutions, err := core.Solve(game, TIMEOUT, MAX_SOLUTIONS)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}

	resp := Response{Error: nil}
	for _, coreSolution := range coreSolutions {
		solution := Solution{Steps: fromCoreSteps(coreSolution.Steps)}
		resp.Solutions = append(resp.Solutions, solution)
	}
	writeSuccess(w, resp)
}

func toCoreSteps(webSteps []Step) []core.Step {
	steps := make([]core.Step, 0, 100)
	for _, step := range webSteps {
		steps = append(steps, core.Step{X: step.X, Y: step.Y, Z: core.Value(step.Z)})
	}
	return steps
}

func fromCoreSteps(coreSteps []core.Step) []Step {
	steps := make([]Step, 0, 100)
	for _, step := range coreSteps {
		steps = append(steps, Step{X: step.X, Y: step.Y, Z: int(step.Z)})
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
	h := &sodokuHandler{}
	// configure middleware around our example
	chain := alice.New(loggingHandler, countingHandler).Then(h)

	mux := http.NewServeMux()
	mux.Handle("/sodoku", chain)

	fmt.Printf("Start listening at http://localhost:3000/\n")
	http.ListenAndServe(":3000", chain)
}
