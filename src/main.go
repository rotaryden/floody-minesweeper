package main

import (
	"fmt"
	"os"
)

func main(){
	fmt.Printf("PID: %d\n\n", os.Getpid())
	for {
		ui := NewUI()
		f, err := NewField(ui.GetSettings())
		ui.SetField(f)
		if err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}
		// make first turn to display field and ask for the first cell to open, then open this cell, 
		// then repeat this process, every time checking the game state after CellOpen()
		for point := ui.GameTurn(); f.OpenCell(point) == GameStateInProgress; point = ui.GameTurn() {}

		// Run GameEndedUI to report outcome of the game and ask question about continuation
		// If user wants to continue - it would return true
		if !ui.GameEnded(){
			fmt.Println("Bye!")
			break
		}
	}
}