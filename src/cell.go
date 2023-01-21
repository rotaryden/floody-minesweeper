package main

const ThisIsHoleMarker = 1000

type CellState int

const (
	CellStateClosed = iota
	CellStateOpen
)

type Cell struct {
	State CellState
	// HolesNumber == 0 - free cell
	// > 0 < HoleMarker - Hole-adjacent cell
	// == HoleMarker = is a hole
	HolesNumber int
}

