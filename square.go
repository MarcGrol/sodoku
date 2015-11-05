package main

import (
	"bytes"
	"fmt"
	"math"
)

const missing_value = 0
const missing_symbol = "_"

type Square struct {
	data        [][]int
	Size        int
	SectionSize int
}

func NewSquare(size int) *Square {
	sectionSize := math.Sqrt(float64(size))
	sq := Square{Size: size, SectionSize: int(sectionSize)}
	sq.data = make([][]int, size)
	for x := 0; x < size; x++ {
		sq.data[x] = make([]int, size)
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

type IterFunc func(x int, y int, z int) error

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

func (sq *Square) Set(x int, y int, value int) {
	sq.data[x][y] = value
}

func (sq *Square) Clear(x int, y int) {
	sq.data[x][y] = missing_value
}

func (sq Square) Has(x int, y int) bool {
	return sq.data[x][y] != missing_value
}

func (sq Square) Get(x int, y int) int {
	return sq.data[x][y]
}

func (sq Square) GetRowValues(x int) []int {
	values := make([]int, 0, sq.Size)
	for y := 0; y < sq.Size; y++ {
		if sq.data[x][y] != missing_value {
			values = append(values, sq.data[x][y])
		}
	}
	return values
}

func (sq Square) GetColumnValues(y int) []int {
	values := make([]int, 0, sq.Size)
	for x := 0; x < sq.Size; x++ {
		if sq.data[x][y] != missing_value {
			values = append(values, sq.data[x][y])
		}
	}
	return values
}

func (sq Square) GetSectionValues(x int, y int) []int {
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

func (sq Square) getNeighbourValues(x int, y int) []int {
	values := make([]int, 0, sq.Size)
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if sq.Exists(x+dx, y+dy) && sq.Has(x+dx, y+dy) {
				values = append(values, sq.Get(x+dx, y+dy))
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
			value := fmt.Sprintf(missing_symbol)
			if sq.data[x][y] != missing_value {
				value = fmt.Sprintf("%d", sq.data[x][y])
			}

			buffer.WriteString(fmt.Sprintf("%s%s", value, delimeter))

		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}
