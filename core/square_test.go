package core

import (
	"reflect"
	"testing"
)

func makeLittleSquare() *Square {
	sq := NewSquare(3)
	sq.Set(0, 0, 1)
	sq.Set(0, 1, 2)
	sq.Set(0, 2, 3)
	sq.Set(1, 0, 4)
	sq.Set(1, 1, 5)
	sq.Set(1, 2, 6)
	sq.Set(2, 0, 7)
	sq.Set(2, 1, 8)
	sq.Set(2, 2, 9)
	return sq
}

func makeBigSquare() *Square {
	sq := NewSquare(9)
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			sq.Set(x, y, Value((x+1)*(y+1)))
		}
	}
	return sq
}

func TestCopy(t *testing.T) {
	sq := makeLittleSquare()
	cp := sq.Copy()
	if !reflect.DeepEqual(sq, cp) {
		t.Errorf("Expected: %v, actual: %v", sq, cp)
	}
}

func TestSetGet(t *testing.T) {
	sq := makeLittleSquare()
	sq.Iterate(func(x, y int, z Value) error {
		sq.Set(x, y, Value(x*y))
		return nil
	})

	sq.Iterate(func(x int, y int, z Value) error {
		v := sq.Get(x, y)
		if v != Value(x*y) {
			t.Errorf("Expected: %v, actual: %v", (x * y), v)
		}
		return nil
	})
}

func TestGetRow(t *testing.T) {
	sq := makeLittleSquare()
	sq.Clear(1, 1)

	var tests = []struct {
		row  int
		vals ValueSet
	}{
		{0, NewValueSet(1, 2, 3)},
		{1, NewValueSet(4, 6)},
		{2, NewValueSet(7, 8, 9)},
	}
	for _, tc := range tests {
		vals := sq.GetRowValues(tc.row)
		if !vals.Equal(tc.vals) {
			t.Errorf("Expected: %v, actual: %v", tc.vals, vals)
		}
	}
}

func TestGetColumn(t *testing.T) {
	sq := makeLittleSquare()
	sq.Clear(2, 2)

	var tests = []struct {
		y    int
		vals ValueSet
	}{
		{0, NewValueSet(1, 4, 7)},
		{1, NewValueSet(2, 5, 8)},
		{2, NewValueSet(3, 6)},
	}
	for _, tc := range tests {
		vals := sq.GetColumnValues(tc.y)
		if !vals.Equal(tc.vals) {
			t.Errorf("%d: Expected: %v, actual: %v", tc.y, tc.vals, vals)
		}
	}
}

func TestGetSection(t *testing.T) {
	sq := makeBigSquare()
	var tests = []struct {
		x, y int
		vals ValueSet
	}{
		{0, 0, NewValueSet(1, 2, 3, 4, 6, 9)},
		{8, 6, NewValueSet(49, 56, 63, 64, 72, 81)},
	}
	for _, tc := range tests {
		vals := sq.GetSectionValues(tc.x, tc.y)
		if !vals.Equal(tc.vals) {
			t.Errorf("%d-%d: Expected: %v, actual: %v", tc.x, tc.y, tc.vals, vals)
		}
	}
}

func TestIsAllowed(t *testing.T) {
	sq := makeBigSquare()

	var tests = []struct {
		x, y      int
		z         Value
		isAllowed bool
	}{
		{1, 1, 1, false},
		{1, 1, 2, false},
		{1, 1, 3, false},
		{1, 1, 4, false},
		{1, 1, 5, true},
		{1, 1, 6, false},
		{1, 1, 7, true},
		{1, 1, 8, false},
		{1, 1, 9, false},
	}
	for _, tc := range tests {
		isAllowed := sq.IsAllowed(tc.x, tc.y, tc.z)
		if isAllowed != tc.isAllowed {
			t.Errorf("%d-%d=%d: Expected: %v, actual: %v", tc.x, tc.y, tc.z,
				tc.isAllowed, isAllowed)
		}
	}

}
