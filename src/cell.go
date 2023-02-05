package main

// denotes the given cel is a mine
const ThisIsMineMarker = 1000

type CellState int

// CellState now is open/closed but chosen not to be boolean to be potentially extendable
const (
	CellStateClosed = iota
	CellStateOpen
)


// Cell sructure, intended to be minimal - MineNumbers is enough to hold any free/mine configurations
type Cell struct {
	State CellState
	// MinesNumber == 0 - free cell
	// > 0 < MineMarker - Mine-adjacent cell
	// == MineMarker = is a mine
	MinesNumber int
}

