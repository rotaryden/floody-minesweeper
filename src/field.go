package main

import (
	"errors"
	"math"
	"math/rand"
)

type Field struct {
	// as we have fixed wodth and height on the start,
	// cells can be kept in the plain array, index arithmetic asumed: y*width+x,
	// also this form is more efficient shuffling elements etc.
	cells  []Cell
	width  int
	height int
	// additional list of pointers pairs to the all holes
	holesRefs []*Cell
	// simple game state
	openCells int
	State GameState
}

func (f *Field) GetWidth() int {
	return f.width
}

func (f *Field) GetHeight() int {
	return f.height
}

func (f *Field) GetCell(x, y int) *Cell {
	return &f.cells[y*f.height+x]
}

func (f *Field) GetState() GameState {
	return f.State
}

// walkNeighbours() walks over all 8 neighbours and runs worker() for each,
func (f *Field) walkNeighbours(x, y int, worker func(*Cell)) {
	// determine real boundaries of the neighbourhood considering field borders
	xStart := int(math.Max(float64(x-1), 0))
	xEnd := int(math.Min(float64(x+1), float64(f.width-1)))
	yStart := int(math.Max(float64(y-1), 0))
	yEnd := int(math.Min(float64(y+1), float64(f.height-1)))

	for dy := yStart; dy <= int(yEnd); dy++ {
		for dx := xStart; dx <= int(xEnd); dx++ {
			// omit current cell
			if dx != x || dy != y {
				worker(f.GetCell(dx, dy))
			}
		}
	}
}

func (f *Field) IsFillable(x, y int) bool {
	// if cell has been open - it should not take part in free-area roll down,
	// as wel las if it is a Hole - holes are handled in the fill() immediatelly
	pc := f.GetCell(x, y)
	if pc.State == CellStateOpen || pc.HolesNumber == ThisIsHoleMarker {
		return false
	}

	// if this is clear free cell - it is just ok
	if pc.HolesNumber == 0 {
		return true
	}

	// but if it is a hole-adjacent cell with a counter - whe should make sure it is adjacent also to a clear free cell
	// this way we will expand only to the borders (with counters) of the free area, no more
	freeNeighborsNumber := 0
	f.walkNeighbours(x, y, func(c *Cell) {
		if c.HolesNumber == 0 {
			freeNeighborsNumber++
		}
	})

	return freeNeighborsNumber > 0
}

// Fill() will be called only if IsFillable() satisfied - so we don't need to do additional checks
func (f *Field) Fill(x, y int) {
	pc := f.GetCell(x, y)
	pc.State = CellStateOpen
	f.openCells++
	// if all cells are open exept holes, the nwe won 
	// - this will be processed at a higher level in OpenCell()
}

// OpenCell opens a cell on the field and changes game state if it is a Hole.
// If the cell has no adjacent holes (clear-free cell), then all free region is opened
func (f *Field) OpenCell(p Point) GameState {	
	pc := f.GetCell(p.X, p.Y)

	if pc.State == CellStateOpen {
		// cell is already opened
		return 0
	}

	if pc.HolesNumber == ThisIsHoleMarker {
		// now, we have to reveal all holes:
		for _, ph := range f.holesRefs {
			ph.State = CellStateOpen
			f.openCells++
		}

		f.State = GameStateLoose
		// game state is clear - return
		return f.State
	}

	if pc.HolesNumber > 0 {
		// this is a hole-adjacent cell, we should open just it
		pc.State = CellStateOpen
		f.openCells++

	} else /*if f.Holes == 0*/ {
		// f.Holes == 0 - this is a clear-free (non-adjacent) cell, 
		// we should flood-fill adjacent free cells and its border with open action
		FloodFill(p.X, p.Y, f)	
	}

	if f.openCells >= len(f.cells) - len(f.holesRefs) {
		// all have been open except holes - we Won!!!
		f.State = GameStateWin
	}

	return f.State
}

// NewField constructs a new game field
func NewField(height, width, holesNumber int) (*Field, error) {
	pfield := new(Field)
	pfield.width = width
	pfield.height = height

	if holesNumber > height * width {
		return nil, errors.New("holesNumber > height * width")
	}

	pfield.State = GameStateInProgress
	pfield.openCells = 0

	cells := make([]Cell, height*width)

	for i := 0; i < holesNumber; i++ {
		cells[i] = Cell{State: CellStateClosed, HolesNumber: ThisIsHoleMarker}
	}

	for i := holesNumber; i < len(cells); i++ {
		cells[i] = Cell{State: CellStateClosed, HolesNumber: 0}
	}

	// Now, make holes to be normally distributed over the field
	rand.Shuffle(len(cells), func(i, j int) {
		cells[i], cells[j] = cells[j], cells[i]
	})

	pfield.cells = cells
	// Holes contains pointers to actual cells
	pfield.holesRefs = make([]*Cell, holesNumber)

	// Now, find all coordinates for shuffled holes and 
	holeIndex := 0
	for y := 0; y < pfield.height; y++ {
		for x := 0; x < pfield.width; x++ {
			pc := pfield.GetCell(x, y)
			// if this is a hole
			if pc.HolesNumber == ThisIsHoleMarker {
				// And hole's cell to the f.holesRefs array for direct tracking
				pfield.holesRefs[holeIndex] = pc
				holeIndex++

				// increment hole-adjacent cells' counters
				pfield.walkNeighbours(x, y, func(pcc *Cell) {
					// if this is not a hole - increase it hole-adjacent counter
					if pcc.HolesNumber != ThisIsHoleMarker {
						pcc.HolesNumber++
					}
				})
			}
		}
	}

	return pfield, nil
}


