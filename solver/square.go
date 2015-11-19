package solver

import (
	"bytes"
	"errors"
	"fmt"
	"math"
)

const (
	MISSING_VALUE  = 0
	MISSING_SYMBOL = "_"
)

type Square struct {
	data        [][]Value
	Size        int
	SectionSize int
}

func NewSquare(size int) *Square {
	sectionSize := math.Sqrt(float64(size))
	sq := Square{Size: size, SectionSize: int(sectionSize)}
	sq.data = make([][]Value, size)
	for x := 0; x < size; x++ {
		sq.data[x] = make([]Value, size)
	}
	return &sq
}

func (sq Square) Copy() *Square {
	nsq := NewSquare(sq.Size)
	for x := 0; x < sq.Size; x++ {
		for y := 0; y < sq.Size; y++ {
			nsq.Set(x, y, sq.Get(x, y))
		}
	}

	return nsq
}

type IterFunc func(x int, y int, z Value) error

func (sq *Square) Iterate(callback IterFunc) error {
	for x := 0; x < sq.Size; x++ {
		for y := 0; y < sq.Size; y++ {
			err := callback(x, y, sq.data[x][y])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (sq Square) Exists(x int, y int) bool {
	return x >= 0 && x < sq.Size && y >= 0 && y < sq.Size
}

func (sq *Square) Set(x int, y int, value Value) {
	sq.data[x][y] = value
}

func (sq *Square) Clear(x int, y int) {
	sq.data[x][y] = MISSING_VALUE
}

func (sq Square) Has(x int, y int) bool {
	return sq.data[x][y] != MISSING_VALUE
}

func (sq Square) Get(x int, y int) Value {
	return sq.data[x][y]
}

func (sq Square) IsAllowed(x int, y int, z Value) bool {
	rowValues := sq.GetRowValues(x)
	columnValues := sq.GetColumnValues(y)
	sectionValues := sq.GetSectionValues(x, y)

	return !rowValues.Contains(z) && !columnValues.Contains(z) && !sectionValues.Contains(z)
}

func (sq Square) GetRowValues(x int) ValueSet {
	values := NewValueSet()
	for y := 0; y < sq.Size; y++ {
		if sq.Has(x, y) {
			values.Add(sq.data[x][y])
		}
	}
	return values
}

func (sq Square) GetColumnValues(y int) ValueSet {
	values := NewValueSet()
	for x := 0; x < sq.Size; x++ {
		if sq.Has(x, y) {
			values.Add(sq.data[x][y])
		}
	}
	return values
}

type point struct {
	x, y int
}

func (sq Square) GetSectionValues(x int, y int) ValueSet {
	centre, _ := getSectionCentre(x, y)
	return sq.getNeighbourValues(centre)
}

func getSectionCentre(x int, y int) (*point, error) {
	var sections = []struct {
		minInclusive point
		maxInclusive point
		centre       point
	}{
		{point{0, 0}, point{2, 2}, point{1, 1}},
		{point{0, 3}, point{2, 5}, point{1, 4}},
		{point{0, 6}, point{2, 8}, point{1, 7}},

		{point{3, 0}, point{5, 2}, point{4, 1}},
		{point{3, 3}, point{5, 5}, point{4, 4}},
		{point{3, 6}, point{5, 8}, point{4, 7}},

		{point{6, 0}, point{8, 2}, point{7, 1}},
		{point{6, 3}, point{8, 5}, point{7, 4}},
		{point{6, 6}, point{8, 8}, point{7, 7}},
	}

	for _, section := range sections {
		if x >= section.minInclusive.x && x <= section.maxInclusive.x &&
			y >= section.minInclusive.y && y <= section.maxInclusive.y {
			return &section.centre, nil
		}
	}
	return nil, errors.New("No centre found")
}

func (sq Square) getNeighbourValues(sectionCentre *point) ValueSet {
	values := NewValueSet()
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if sq.Exists(sectionCentre.x+dx, sectionCentre.y+dy) &&
				sq.Has(sectionCentre.x+dx, sectionCentre.y+dy) {
				values.Add(sq.Get(sectionCentre.x+dx, sectionCentre.y+dy))
			}
		}
	}
	return values
}

func (sq Square) String() string {
	var buffer bytes.Buffer
	for x := 0; x < sq.Size; x++ {
		for y := 0; y < sq.Size; y++ {
			delimeter := ""
			if y < sq.Size-1 {
				delimeter = " "
			}
			value := fmt.Sprintf(MISSING_SYMBOL)
			if sq.Has(x, y) {
				value = fmt.Sprintf("%d", sq.data[x][y])
			}

			buffer.WriteString(fmt.Sprintf("%s%s", value, delimeter))

		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}
