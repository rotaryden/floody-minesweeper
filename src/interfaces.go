package main

type GameState int

const (
	GameStateInProgress GameState = iota
	GameStateLoose
	GameStateWin
)


// IFloodableField interface is an abstraction of the game field, owns states of cells
// This abstraction allows to flood-fill spaces with various logic over cells
// (logic in our case:
//	  we need to flood-fill all clear-free cells + cells with countrs on the closure border of the free region)
// - IsFillable() - test if the cell is fillable (so can be filled with a given state by the fill())
// - Fill()  - knows how to fill the cell with a new state
// - GetWidth() int - abstract width
// - GetHeight() int - abstract height
type IFloodableField interface {
	IsFillable(x, y int) bool
	Fill(x, y int)
	GetWidth() int
	GetHeight() int
}

// IField interface intended for game field consumers, such as UI
// - GetCell() - return a pointer to Cell object by coordinates
type IField interface {
	IFloodableField
	GetCell(int, int) *Cell
	GetState() GameState

}


