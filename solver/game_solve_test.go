package solver

import "testing"

func TestSolveSingleSolution(t *testing.T) {
	gameData := `1 _ 3 _ _ 6 _ 8 _
_ 5 _ _ 8 _ 1 2 _
7 _ 9 1 _ 3 _ 5 6
_ 3 _ _ 6 7 _ 9 _
5 _ 7 8 _ _ _ 3 _
8 _ 1 _ 3 _ 5 _ 7
_ 4 _ _ 7 8 _ 1 _
6 _ 8 _ _ 2 _ 4 _
_ 1 2 _ 4 5 _ 7 8
`
	game, err := LoadString(gameData)
	if err != nil {
		t.Errorf("Error loading game: %s", err)
	}
	solutions, err := Solve(game, 10, 1)
	if err != nil {
		t.Errorf("Error solving game: %s", err)
	} else {
		if len(solutions) != 1 {
			t.Errorf("expected: %d, actual: %d", 1, len(solutions))
		}
	}
}

func TestSolveMultipleSolutions(t *testing.T) {
	gameData := `1 _ 3 _ _ 6 _ 8 _
_ 5 _ _ 8 _ 1 2 _
7 _ 9 1 _ 3 _ 5 6
_ 3 _ _ 6 7 _ 9 _
5 _ 7 8 _ _ _ 3 _
8 _ 1 _ 3 _ 5 _ 7
_ 4 _ _ 7 8 _ 1 _
6 _ 8 _ _ 2 _ 4 _
_ 1 2 _ 4 5 _ 7 8
`
	game, err := LoadString(gameData)
	if err != nil {
		t.Errorf("Error loading game: %s", err)
	}
	solutions, err := Solve(game, 10, 4)
	if err != nil {
		t.Errorf("Error solving game: %s", err)
	} else {
		if len(solutions) != 4 {
			t.Errorf("expected: %d, actual: %d", 4, len(solutions))
		}
	}
}

func TestSolveHardSolution(t *testing.T) {
	gameData := `9 _ _ _ 2 _ _ _ 5
_ _ _ 9 _ 5 _ _ _
_ _ 7 _ 6 _ 4 _ _
_ 5 _ _ _ _ _ 7 _
8 _ 1 _ 7 _ 6 _ 2
_ 2 _ _ _ _ _ 3 _
_ _ 6 _ 4 _ 1 _ _
_ _ _ 8 _ 3 _ _ _
4 _ _ _ 9 _ _ _ 3
`
	game, err := LoadString(gameData)
	if err != nil {
		t.Errorf("Error loading game: %s", err)
	}
	solutions, err := Solve(game, 3, 1)
	if err != nil {
		t.Errorf("Error solving game: %s", err)
	} else {
		if len(solutions) != 1 {
			t.Errorf("expected: %d, actual: %d", 1, len(solutions))
		}
	}
}

func TestSolveHardestSolution(t *testing.T) {
	gameData := `8 _ _ _ _ _ _ _ _
_ _ 3 6 _ _ _ _ _
_ 7 _ _ 9 _ 2 _ _
_ 5 _ _ _ 7 _ _ _
_ _ _ _ 4 5 7 _ _
_ _ _ 1 _ _ _ 3 _
_ _ 1 _ _ _ _ 6 8
_ _ 8 5 _ _ _ 1 _
_ 9 _ _ _ _ 4 _ _
`
	game, err := LoadString(gameData)
	if err != nil {
		t.Errorf("Error loading game: %s", err)
	}
	solutions, err := Solve(game, 1, 1)
	if err != nil {
		t.Errorf("Error solving game: %s", err)
	} else {
		if len(solutions) != 1 {
			t.Errorf("expected: %d, actual: %d", 1, len(solutions))
		}
	}
}

func BenchmarkHardestSolution(b *testing.B) {
	gameData := `8 _ _ _ _ _ _ _ _
_ _ 3 6 _ _ _ _ _
_ 7 _ _ 9 _ 2 _ _
_ 5 _ _ _ 7 _ _ _
_ _ _ _ 4 5 7 _ _
_ _ _ 1 _ _ _ 3 _
_ _ 1 _ _ _ _ 6 8
_ _ 8 5 _ _ _ 1 _
_ 9 _ _ _ _ 4 _ _
`

	for n := 0; n < b.N; n++ {
		game, _ := LoadString(gameData)
		Solve(game, 3, 1)
	}
}
