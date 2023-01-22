package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Field struct {
	// as we have fixed wodth and height on the start,
	// cells can be kept in the plain array, index arithmetic asumed: y*width+x,
	// also this form is more efficient shuffling elements etc.
	cells  []*Cell
	width  int
	height int
	// additional list of pointers to the all hole cells 
	// - needed for quick reveal of all holes on game over
	holesRefs []*Cell
	// simple game state
	openCells int
	State     GameState
}

func (f *Field) GetWidth() int {
	return f.width
}

func (f *Field) GetHeight() int {
	return f.height
}

// reaches the cell using index arithmetic
func (f *Field) GetCell(x, y int) *Cell {
	if x < 0 || x >= f.width || y < 0 || y >= f.height {
		panic(fmt.Sprintf("x=%d, y=%d out of boundaries!", x, y))
	}
	return f.cells[y * f.width +x ]
}

func (f *Field) GetState() GameState {
	return f.State
}

// walkNeighbours() walks over all 8 neighbours and runs worker() for each,
// unnless worker() returns true before
func (f *Field) walkNeighbours(x, y int, worker func(*Cell) bool) {
	// determine real boundaries of the neighbourhood considering field borders
	xStart := int(math.Max(float64(x-1), 0))
	xEnd := int(math.Min(float64(x+1), float64(f.width-1)))
	yStart := int(math.Max(float64(y-1), 0))
	yEnd := int(math.Min(float64(y+1), float64(f.height-1)))

	for dy := yStart; dy <= int(yEnd); dy++ {
		for dx := xStart; dx <= int(xEnd); dx++ {
			// omit current cell
			if dx != x || dy != y {
				if worker(f.GetCell(dx, dy)) {
					return
				}
			}
		}
	}
}

// This is a central classifying predicate in the flood-fill algorithm
// it is responsible for classify connected free cells area with adjacent border of counter cells (hole-adjacent cells)
func (f *Field) IsFillable(x, y int, isFirstCell bool) bool {
	// if cell has been open - it should not take part in free-area roll down,
	// as wel las if it is a Hole - holes are handled in the fill() immediatelly
	pc := f.GetCell(x, y)
	if pc.State == CellStateOpen || pc.HolesNumber == ThisIsHoleMarker {
		return false
	}

	// if this is THE FIRST clear free cell in the free region being opened - open it.
	//Explanation of the isFirstCell functionality: it should be fillable ONLY if all clear free neighbours are closed, to avoid a bug like this:
	//    a b c d e f g h i j
	//  0 ⛶ 1 🙫 1 ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ 
	//  1 ⛶ 1 1 1 ⛶ ⛶ ⛶ ⛶ 1 1 
	//  2 ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ 1 🙫   <<----- original free region - has been opened e.g. by "e2"
	//  3 ⛶ ⛶ ⛶ ⛶ ⛶ 1 1 1 1 1 
	//  4 ⛶ ⛶ ⛶ ⛶ ⛶ 1 🙫 1 ⛶ ⛶  <<--- this 2nd region has been opened despite it is not connected to the original free region - Bug!
	//  5 ⛶ 1 1 1 ⛶ 1 1 2 1 1 
	//  6 ⛶ 1 🙫 1 ⛶ ⛶ ⛶ 1 🙫 🙫 
	//  7 ⛶ 1 1 1 ⛶ ⛶ ⛶ 1 1 1 
	//  8 ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ 
	//  9 ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ ⛶ 
	if pc.HolesNumber == 0 && isFirstCell {
		return true
	}

	// if it is a hole-adjacent cell with a counter - whe should make sure it is adjacent also to a clear-free (non-adjacent) cell
	// moreover, this cell should be open to mitigate situations like this (h3 is ahjacent to 2 independent free regions):
	// 	   a b c d e f g h i j
	//  0 🙫 2 🙫 1 ⛶ ⛶ ⛶ 1 🙫 🙫
	//  1 🙫 3 1 1 ⛶ ⛶ ⛶ 1 🙫 🙫
	//  2 🙫 1 ⛶ ⛶ ⛶ ⛶ ⛶ 1 🙫 🙫
	//  3 🙫 2 1 1 1 1 1 1 1 1          
	//  4 🙫 🙫 🙫 🙫 🙫 🙫 🙫 1 ⛶ ⛶  <<<----- 2nd disconnected region has opened - Bug!
	//  5 🙫 🙫 🙫 🙫 🙫 2 🙫 3 2 1
	//  6 🙫 🙫 🙫 🙫 🙫 🙫 🙫 🙫 🙫 🙫
	//  7 🙫 🙫 🙫 🙫 🙫 🙫 🙫 🙫 🙫 🙫
	//  8 🙫 🙫 🙫 🙫 🙫 🙫 🙫 🙫 🙫 🙫
	//  9 🙫 🙫 🙫 🙫 🙫 🙫 🙫 🙫 🙫 🙫

	// this way we will expand only to the borders (with counters) of the free area, no more
	hasFreeOpenNeighbour := false
	f.walkNeighbours(x, y, func(c *Cell) bool {
		if c.HolesNumber == 0 && c.State == CellStateOpen {
			hasFreeOpenNeighbour = true
			return true //stop walking through neighbours
		}
		return false
	})

	return hasFreeOpenNeighbour
}

// Fill() will be called only if IsFillable() satisfied - so we don't need to do additional checks
func (f *Field) Fill(x, y int) {
	pc := f.GetCell(x, y)
	pc.State = CellStateOpen
	f.openCells++
	// if all cells are open exept holes, the nwe won
	// - this will be processed at a higher level in OpenCell()
}

func (f *Field) revealHoles() {
	for _, h := range f.holesRefs {
		h.State = CellStateOpen
	}
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
		f.revealHoles()

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

	if f.openCells >= len(f.cells)-len(f.holesRefs) {
		// all have been open except holes - we Won!!!
		f.State = GameStateWin
		// yet still good to show all holes
		f.revealHoles()
	}

	return f.State
}

// NewField constructs a new game field
func NewField(gs *GameSettings) (*Field, error) {
	pfield := new(Field)
	pfield.width = gs.Width
	pfield.height = gs.Height
	holesNumber := gs.HolesNumber

	if holesNumber > gs.Height*gs.Width {
		return nil, errors.New("holesNumber > height * width")
	}

	pfield.State = GameStateInProgress
	pfield.openCells = 0

	cells := make([]*Cell, gs.Height*gs.Width)

	for i := 0; i < holesNumber; i++ {
		cells[i] = &Cell{State: CellStateClosed, HolesNumber: ThisIsHoleMarker}
	}

	for i := holesNumber; i < len(cells); i++ {
		cells[i] = &Cell{State: CellStateClosed, HolesNumber: 0}
	}

	// seed random generator
	rand.Seed(time.Now().UnixNano())
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
				pfield.walkNeighbours(x, y, func(pcc *Cell) bool {
					// if this is not a hole - increase it hole-adjacent counter
					if pcc.HolesNumber != ThisIsHoleMarker {
						pcc.HolesNumber++
					}
					return false
				})
			}
		}
	}

	return pfield, nil
}
