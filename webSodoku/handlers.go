package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/MarcGrol/sodoku/solver"
)

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

type Response struct {
	Error     *ErrorDescriptor `json:"error"`
	Solutions []Game           `json:"solutions"`
}

type ErrorDescriptor struct {
	Message string `json:"message"`
}

type Game struct {
	Board Board  `json:"board"`
	Steps []Step `json:"steps"`
}

type Board struct {
	Rows [solver.SQUARE_SIZE][solver.SQUARE_SIZE]int `json:"rows"`
}

type Step struct {
	X       int  `json:"x"`
	Y       int  `json:"y"`
	Z       int  `json:"z"`
	Initial bool `json:"initial"`
	IsGuess bool `json:"isGuess"`
}

func FromJson(reader io.Reader) (*Game, error) {
	game := Game{}
	err := json.NewDecoder(reader).Decode(&game)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (resp Response) ToJson(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(resp)
}

type sodokuHandler struct {
	timeout      int
	minSolutions int
}

func (sh *sodokuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		sh.get(w, r)
	case "POST":
		sh.post(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("Only GET and POST are supported"))
	}
}

func (sh sodokuHandler) get(w http.ResponseWriter, r *http.Request) {
	game, err := solver.LoadString(HARD_EXAMPLE)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
	}
	sh.doSolve(w, r, game)
}

func (sh sodokuHandler) post(w http.ResponseWriter, r *http.Request) {
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
	sh.doSolve(w, r, game)
}

func (sh sodokuHandler) doSolve(w http.ResponseWriter, r *http.Request, game *solver.Game) {
	coreSolutions, err := solver.Solve(game, sh.timeout, sh.minSolutions)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}

	resp := Response{Error: nil}
	for _, coreSolution := range coreSolutions {
		board := boardFromSteps(coreSolution.Steps)
		solution := Game{
			Board: board,
			Steps: fromCoreSteps(coreSolution.Steps),
		}
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

func boardFromSteps(coreSteps []solver.Step) Board {
	board := Board{}
	for _, step := range coreSteps {
		board.Rows[step.X][step.Y] = int(step.Z)
	}
	return board
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
