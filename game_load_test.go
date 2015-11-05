package main

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
	game, err := Load(gameData)
	if err != nil {
		t.Errorf("Error loading game: %s", err)
	}
	if gameData != game.String() {
		t.Errorf("expected: %s, actual: %s", gameData, game.String())
	}
}

func TestContainsDuplicates(t *testing.T) {
	var hasDupsTests = []struct {
		description string
		input       []int
		hasDups     bool
		dupVal      int
	}{
		{"No vals", []int{}, false, -1},
		{"No dups", []int{1, 2, 3}, false, -1},
		{"No dups", []int{1, 2, 3, 1}, true, 1},
	}
	for _, tc := range hasDupsTests {
		hasDups, dupVal := containsDuplicates(tc.input)
		if hasDups != tc.hasDups {
			t.Errorf("%s: hasDuplicates expected: %v, actual: %v", tc.description, tc.dupVal, hasDups)
		}
		if hasDups && dupVal != tc.dupVal {
			t.Errorf("%s: duplicateValue expected: %v, actual: %v", tc.description, tc.dupVal, hasDups)
		}
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
		{"non int or _ value", "1 _ a _ 6 _ 8 _ 9\n", "Invalid number 'a' for item row:1, column:3"},
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
`, "Row 2 contains duplicate value 8"},
		{"duplicates in column", `1 _ 3 _ _ 6 _ 8 _
_ 5 _ _ 8 _ 1 2 _
7 _ 9 1 _ 3 _ 5 6
_ 3 _ _ 6 7 _ 9 _
6 _ 7 8 _ _ _ 3 _
8 _ 1 _ 3 _ 5 _ 7
_ 4 _ _ 7 8 _ 1 _
6 _ 8 _ _ 2 _ 4 _
_ 1 2 _ 4 5 _ 7 8
`, "Column 1 contains duplicate value 6"},
		{"duplicates in section", `1 _ 3 _ _ 6 _ 8 _
_ 5 _ _ 8 _ 1 2 _
7 _ 9 1 _ 3 _ 5 6
_ 3 3 _ 6 7 _ 9 _
5 _ 7 8 _ _ _ 3 _
8 _ 1 _ 3 _ 5 _ 7
_ 4 _ _ 7 8 _ 1 _
6 _ 8 _ _ 2 _ 4 _
_ 1 2 _ 4 5 _ 7 8
`, "Row 4 contains duplicate value 3"},
	}

	for _, tc := range invalidInputTests {
		_, err := Load(tc.input)
		if err.Error() != tc.expectedError {
			t.Errorf("%s: Expected: %s, actual: %s", tc.description, tc.expectedError, err.Error())
		}
	}
}
