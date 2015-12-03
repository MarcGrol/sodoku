package solver

import (
	"fmt"
	"strconv"
	"strings"
)

func LoadSteps(steps []Step) (*Game, error) {
	game := newGame()
	for idx, step := range steps {
		if !game.square.Exists(step.X, step.Y) {
			return nil, fmt.Errorf("Invalid offset: %d-%d for step %d", step.X, step.Y, idx)
		}
		if !game.square.IsAllowed(step.X, step.Y, Value(step.Z)) {
			return nil, fmt.Errorf("Duplicate value %d for item row:%d, column:%d for step %d",
				step.Z, step.X, step.Y, idx)
		}
		game.set(step.X, step.Y, Value(step.Z), true, false)
	}
	game.CellsToBeSolved = game.countEmptyValues()
	return game, nil
}

func LoadString(lines string) (*Game, error) {
	linesRead := 0
	game := newGame()
	for x, line := range strings.Split(lines, "\n") {

		if x >= SQUARE_SIZE {
			break
		}
		if line == "" {
			break
		}

		splitted := strings.Split(line, " ")
		if len(splitted) != SQUARE_SIZE {
			return nil, fmt.Errorf("Invalid number of columns for row %d: needs %d, actual %d", x+1,
				SQUARE_SIZE, len(splitted))
		}
		for y, val := range splitted {
			if val == "_" {
				game.square.Clear(x, y)
				continue
			}
			num, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("Invalid value '%s' for item row:%d, column:%d", val, x+1, y+1)
			}
			if num < 0 || num > SQUARE_SIZE {
				return nil, fmt.Errorf("Invalid value %d for item row:%d, column:%d", num, x+1, y+1)

			}
			if !game.square.IsAllowed(x, y, Value(num)) {
				return nil, fmt.Errorf("Duplicate value %d for item row:%d, column:%d", num, x+1, y+1)
			}
			game.set(x, y, Value(num), true, false)
		}
		linesRead++
	}
	if linesRead != SQUARE_SIZE {
		return nil, fmt.Errorf("Not enough rows: needs %d, actual %d", SQUARE_SIZE, linesRead)
	}

	game.CellsToBeSolved = game.countEmptyValues()
	return game, nil
}
