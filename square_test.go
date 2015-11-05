package main

import (
	"reflect"
	"testing"
)

func makeSquare() *Square {
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

func TestCopy(t *testing.T) {
	sq := makeSquare()
	cp := sq.Copy()
	if !reflect.DeepEqual(sq, cp) {
		t.Errorf("Expected: %v, actual: %v", sq, cp)
	}
}

func TestGetRow(t *testing.T) {
	sq := makeSquare()
	vals := sq.GetRowValues(1)
	if len(vals) != 3 {
		t.Errorf("Expected: %d, actual: %d", 3, len(vals))
	}
	if vals[0] != 4 {
		t.Errorf("Expected: %d, actual: %d", 4, vals[0])
	}
	if vals[1] != 5 {
		t.Errorf("Expected: %d, actual: %d", 5, vals[1])
	}
	if vals[2] != 6 {
		t.Errorf("Expected: %d, actual: %d", 6, vals[2])
	}
}

func TestGetRowWithMissing(t *testing.T) {
	sq := makeSquare()
	sq.Clear(0, 2)
	vals := sq.GetRowValues(0)
	if len(vals) != 2 {
		t.Errorf("Expected: %d, actual: %d", 2, len(vals))
	}
	if vals[0] != 1 {
		t.Errorf("Expected: %d, actual: %d", 1, vals[0])
	}
	if vals[1] != 2 {
		t.Errorf("Expected: %d, actual: %d", 2, vals[1])
	}
}

func TestGetColumn(t *testing.T) {
	sq := makeSquare()
	vals := sq.GetColumnValues(1)
	if vals[0] != 2 {
		t.Errorf("Expected: %d, actual: %d", 2, vals[0])
	}
	if vals[1] != 5 {
		t.Errorf("Expected: %d, actual: %d", 5, vals[1])
	}
	if vals[2] != 8 {
		t.Errorf("Expected: %d, actual: %d", 8, vals[2])
	}
}

func TestGetColumnWithMissing(t *testing.T) {
	sq := makeSquare()
	sq.Clear(0, 0)
	vals := sq.GetColumnValues(0)
	if len(vals) != 2 {
		t.Errorf("Expected: %d, actual: %d", 2, len(vals))
	}
	if vals[0] != 4 {
		t.Errorf("Expected: %d, actual: %d", 4, vals[0])
	}
	if vals[1] != 7 {
		t.Errorf("Expected: %d, actual: %d", 7, vals[1])
	}
}

func TestGetSection(t *testing.T) {
	sq := NewSquare(9)
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			sq.Set(x, y, x*y)
		}
	}
	{
		expected := []int{1, 2, 2, 4}
		actual := sq.GetSectionValues(0, 0)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected: %v, actual: %v", expected, actual)
		}
	}
	{
		expected := []int{36, 42, 48, 42, 49, 56, 48, 56, 64}
		actual := sq.GetSectionValues(8, 6)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected: %v, actual: %v", expected, actual)
		}
	}
}
