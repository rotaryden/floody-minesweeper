package main

// denotes the given cel is a hole
const ThisIsHoleMarker = 1000

type CellState int

// CellState now is open/closed but chosen not to be boolean to be potentially extendable
const (
	CellStateClosed = iota
	CellStateOpen
)


// Cell sructure, intended to be minimal - HoleNumbers is enough to hold any free/hole configurations
type Cell struct {
	State CellState
	// HolesNumber == 0 - free cell
	// > 0 < HoleMarker - Hole-adjacent cell
	// == HoleMarker = is a hole
	HolesNumber int
}

