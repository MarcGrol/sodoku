package solver

import "testing"

func TestLoadGameSuccess(t *testing.T) {
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
	if gameData != game.Dump() {
		t.Errorf("expected: %s, actual: %s", gameData, game.Dump())
	}
}

func TestLoadGameStepsOffsetError(t *testing.T) {
	steps := []Step{
		Step{X: -1, Y: 0, Z: 1},
		Step{X: 1, Y: 1, Z: 1},
	}
	_, err := LoadSteps(steps)
	if err.Error() != "Invalid offset: -1-0 for step 0" {
		t.Errorf("expected error: %s, actual: %s", "Invalid offset: -1-0 for step 0", err)
	}
}

func TestLoadGameStepsValueError(t *testing.T) {
	steps := []Step{
		Step{X: 1, Y: 0, Z: 1},
		Step{X: 1, Y: 1, Z: 1},
	}
	_, err := LoadSteps(steps)
	if err.Error() != "Duplicate value 1 for item row:1, column:1 for step 1" {
		t.Errorf("expected error: %s, actual: %s", "Duplicate value 1 for item row:1, column:1 for step 1", err)
	}
}

func TestLoadGameStepsSuccess(t *testing.T) {
	steps := []Step{
		Step{X: 1, Y: 0, Z: 1},
		Step{X: 1, Y: 1, Z: 2},
	}
	_, err := LoadSteps(steps)
	if err != nil {
		t.Errorf("Error loading game: %s", err)
	}
}

func TestLoadGameErrors(t *testing.T) {
	var invalidInputTests = []struct {
		description   string
		input         string
		expectedError string
	}{
		{"not enough rows", "1 _ 3 _ _ 6 _ 8 _\n", "Not enough rows: needs 9, actual 1"},
		{"not enough columns", "1 _ 3 _ _ 6 _ 8\n", "Invalid number of columns for row 1: needs 9, actual 8"},
		{"too many columns", "1 _ 3 _ _ 6 _ 8 _ 9\n", "Invalid number of columns for row 1: needs 9, actual 10"},
		{"non int or _ value", "1 _ a _ 6 _ 8 _ 9\n", "Invalid value 'a' for item row:1, column:3"},
		{"negative value", "1 _ 3 _ -1 _ 8 _ 9\n", "Invalid value -1 for item row:1, column:5"},
		{"too high value", "1 _ 3 _ 4 _ 10 _ 9\n", "Invalid value 10 for item row:1, column:7"},
		{"duplicates in row", `1 _ 3 _ _ 6 _ 8 _
_ 5 _ _ 8 _ 1 8 _
7 _ 9 1 _ 3 _ 5 6
_ 3 _ _ 6 7 _ 9 _
5 _ 7 8 _ _ _ 3 _
8 _ 1 _ 3 _ 5 _ 7
_ 4 _ _ 7 8 _ 1 _
6 _ 8 _ _ 2 _ 4 _
_ 1 2 _ 4 5 _ 7 8
`, "Duplicate value 8 for item row:2, column:8"},
		{"duplicates in column", `1 _ 3 _ _ 6 _ 8 _
_ 5 _ _ 8 _ 1 2 _
7 _ 9 1 _ 3 _ 5 6
_ 3 _ _ 6 7 _ 9 _
6 _ 7 8 _ _ _ 3 _
8 _ 1 _ 3 _ 5 _ 7
_ 4 _ _ 7 8 _ 1 _
6 _ 8 _ _ 2 _ 4 _
_ 1 2 _ 4 5 _ 7 8
`, "Duplicate value 6 for item row:8, column:1"},
		{"duplicates in section", `1 _ 3 _ _ 6 _ 8 _
_ 5 _ _ 8 _ 1 2 _
7 _ 9 1 _ 3 _ 5 6
_ 3 3 _ 6 7 _ 9 _
5 _ 7 8 _ _ _ 3 _
8 _ 1 _ 3 _ 5 _ 7
_ 4 _ _ 7 8 _ 1 _
6 _ 8 _ _ 2 _ 4 _
_ 1 2 _ 4 5 _ 7 8
`, "Duplicate value 3 for item row:4, column:3"},
	}

	for _, tc := range invalidInputTests {
		_, err := LoadString(tc.input)
		if err.Error() != tc.expectedError {
			t.Errorf("%s: Expected: %s, actual: %s", tc.description, tc.expectedError, err.Error())
		}
	}
}
