package web

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/MarcGrol/ctx"
	"github.com/MarcGrol/logging"
	"github.com/MarcGrol/sodoku/solver"

	"golang.org/x/net/context"
)

const hardExample = `9 _ _ _ 2 _ _ _ 5
_ _ _ 9 _ 5 _ _ _
_ _ 7 _ 6 _ 4 _ _
_ 5 _ _ _ _ _ 7 _
8 _ 1 _ 7 _ 6 _ 2
_ 2 _ _ _ _ _ 3 _
_ _ 6 _ 4 _ 1 _ _
_ _ _ 8 _ 3 _ _ _
4 _ _ _ 9 _ _ _ 3
`

type response struct {
	Error     *errorDescriptor `json:"error"`
	Solutions []game           `json:"solutions"`
}

type errorDescriptor struct {
	Message string `json:"message"`
}

type game struct {
	Board board  `json:"board"`
	Steps []step `json:"steps"`
}

type board struct {
	Rows [solver.SQUARE_SIZE][solver.SQUARE_SIZE]int `json:"rows"`
}

type step struct {
	X       int  `json:"x"`
	Y       int  `json:"y"`
	Z       int  `json:"z"`
	Initial bool `json:"initial"`
	IsGuess bool `json:"isGuess"`
}

func fromJSON(reader io.Reader) (*game, error) {
	g := game{}
	err := json.NewDecoder(reader).Decode(&g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (resp response) toJSON(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(resp)
}

// SodokuHandler exposes sodoku solver as HTTP endpoint
type SodokuHandler struct {
	timeout      int
	minSolutions int
	context      context.Context
}

func (sh *SodokuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// inject context in application
	ctx := ctx.New.CreateContext(r)
	log := logging.New()

	log.Warning(ctx, "Got %s on url %s", r.Method, r.RequestURI)

	switch r.Method {
	case "GET":
		sh.get(w, r)
	case "POST":
		sh.post(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("Only GET and POST are supported"))
	}
}

func (sh SodokuHandler) get(w http.ResponseWriter, r *http.Request) {

	game, err := solver.LoadString(hardExample)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
	}
	sh.doSolve(w, r, game)
}

func (sh SodokuHandler) post(w http.ResponseWriter, r *http.Request) {
	// read incoming steps
	input, err := fromJSON(r.Body)
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

func (sh SodokuHandler) doSolve(w http.ResponseWriter, r *http.Request, gameToSolve *solver.Game) {
	coreSolutions, err := solver.Solve(gameToSolve, sh.timeout, sh.minSolutions)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}

	resp := response{Error: nil}
	for _, coreSolution := range coreSolutions {
		solution := game{
			Board: boardFromSteps(coreSolution.Steps),
			Steps: fromCoreSteps(coreSolution.Steps),
		}
		resp.Solutions = append(resp.Solutions, solution)
	}
	writeSuccess(w, resp)
}

func toCoreSteps(webSteps []step) []solver.Step {
	steps := make([]solver.Step, 0, 100)
	for _, step := range webSteps {
		steps = append(steps, solver.Step{X: step.X, Y: step.Y, Z: solver.Value(step.Z), Initial: step.Initial, IsGuess: step.IsGuess})
	}
	return steps
}

func fromCoreSteps(coreSteps []solver.Step) []step {
	steps := make([]step, 0, 100)
	for _, s := range coreSteps {
		steps = append(steps,
			step{X: s.X, Y: s.Y, Z: int(s.Z), Initial: s.Initial, IsGuess: s.IsGuess},
		)
	}
	return steps
}

func boardFromSteps(coreSteps []solver.Step) board {
	board := board{}
	for _, step := range coreSteps {
		board.Rows[step.X][step.Y] = int(step.Z)
	}
	return board
}

func writeError(w http.ResponseWriter, status int, err error) {
	resp := response{Error: &errorDescriptor{Message: err.Error()}}
	// write headers
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	resp.toJSON(w)
}

func writeSuccess(w http.ResponseWriter, resp response) {
	// write headers
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp.toJSON(w)
}
