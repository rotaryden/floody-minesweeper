package main

import (
	"fmt"
	"os"
)

func main(){
	for {
		settings := NewGameUI()
		f, err := NewField(settings.Height, settings.Width, settings.HolesNumber)
		if err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}
		for point := GameTurnUI(f); f.OpenCell(point) == GameStateInProgress; point = GameTurnUI(f) {
			
		}
	}
}