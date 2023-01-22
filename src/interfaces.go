package main

type GameState int

const (
	GameStateInProgress GameState = iota
	GameStateLoose
	GameStateWin
)

// common part of interfaces (interface mixin)
type IMeasurable interface {
	// abstract width
	GetWidth() int
	// abstract height
	GetHeight() int
}

// IFloodableField interface is an abstraction of the game field, owns states of cells
// This abstraction decouples flood-fill algorithm and allows it to run on spaces with various logic over cells
// (logic in the Proxx case:
//	  we need to flood-fill all free and hole-adjacent cells tjose all are adjucent to the current free cells region being opened
type IFloodableField interface {
	IMeasurable
	// Main flood-fill predicates
	// test if the cell can be filled with a given state, 
	// "has the old color" in terms of flood-fill, 
	// - is free or hole-adjacent (so  by the fill())
	IsFillable(x, y int, isFirstCell bool) bool
	// knows how to fill the cell with a new state
	Fill(x, y int)
}

// IField interface intended for decoupling game field consumers, such as UI
type IField interface {
	IMeasurable
	// return a pointer to Cell object by coordinates
	GetCell(int, int) *Cell
	// returns Game State
	GetState() GameState
}


