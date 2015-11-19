package solver

import (
	"bytes"
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

func (sq Square) GetSectionValues(x int, y int) ValueSet {
	cx, cy := sq.getSectionCentre(x, y)
	return sq.getNeighbourValues(cx, cy)
}

func (sq Square) getSectionCentre(x int, y int) (centreX int, centreY int) {
	// TODO would this work for 16x16 or 25x25?
	remainderX := x % sq.SectionSize
	if remainderX == 0 {
		// above centre
		centreX = x + 1
	} else if remainderX == 1 {
		// is centre
		centreX = x
	} else if remainderX == 2 {
		// below centre
		centreX = x - 1
	} else {
		// TODO should not happen
	}

	remainderY := y % sq.SectionSize
	if remainderY == 0 {
		// left of centre
		centreY = y + 1
	} else if remainderY == 1 {
		// is centre
		centreY = y
	} else if remainderY == 2 {
		// right of centre
		centreY = y - 1
	} else {
		// TODO should not happen
	}

	return centreX, centreY
}

func (sq Square) getNeighbourValues(x int, y int) ValueSet {
	values := NewValueSet()
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if sq.Exists(x+dx, y+dy) && sq.Has(x+dx, y+dy) {
				values.Add(sq.Get(x+dx, y+dy))
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
