package main

import "testing"

func TestSolve(t *testing.T) {
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
	game, err := Load(gameData)
	if err != nil {
		t.Errorf("Error loading game: %s", err)
	}
	steps, err := game.Solve(0)
	if err != nil {
		t.Errorf("Error solving game: %s", err)
	} else {
		t.Logf("Completed after %d steps", steps)
	}
}
