package main

import (
	"fmt"
	"os"
)

func main(){
	// for debugging
	fmt.Printf("PID: %d\n\n", os.Getpid())
		ui := NewUI()
	
	for {
		ui.Restart()
		// create game field, it is decoupled from UI via IField interface
		f, err := NewField(ui.GetSettings())
		ui.SetField(f)
		if err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}
		// make first turn to display field and ask for the first cell to open, then open this cell, 
		// then repeat this process, every time checking the game state after CellOpen()
		for point, gameStopped := ui.GameTurn(); 
			! gameStopped && f.OpenCell(point) == GameStateInProgress; 
			point, gameStopped = ui.GameTurn() {}

		// Run GameEndedUI to report outcome of the game and ask question about continuation
		// If user wants to continue - it would return true
		if !ui.GameEnded(){
			fmt.Println("Bye!")
			break
		}
	}
}