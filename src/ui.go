package main

// Very simple terminal UI for the game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const letters = "abcdefghijklmnopqrstuvwxyz"
const lettersBase = int('a')

func printField(f IField) {
	h := f.GetHeight()
	w := f.GetWidth()

	//offset
	fmt.Printf("  ")
	// letter denotion of columns
	for _, c := range letters[0:w] {
		fmt.Printf("%2c", c)
	}
	fmt.Println()
	for y := 0; y < h; y++ {
		// print line number
		fmt.Printf("%2d ", y)
		for x := 0; x < w; x++ {
			pc := f.GetCell(x, y)
			switch {
			case pc.State == CellStateClosed:
				// add space after placeholder "%s " to fix wide rune printing issue
				fmt.Printf("%c ", '🙫')
			case pc.HolesNumber == ThisIsHoleMarker:
				fmt.Printf("%c ", '⦿')
			case pc.HolesNumber == 0:
				fmt.Printf("%c ", '⛶')
			default:
				fmt.Printf("%2d", pc.HolesNumber)
			}
		}
		fmt.Println()
	}
}

func getPrintableGameStatus(f IField) string {
	s := f.GetState()
	switch s {
	case GameStateInProgress:
		return "In progress"
	case GameStateLoose:
		return "You Lost !"
	case GameStateWin:
		return "You Won !!!"
	}
	return "Some error..."
}

func printHeader() {
	fmt.Println("\n\n\n==================================")
	fmt.Println("============= Proxx ==============")
	fmt.Println("==================================")
}

func printTurnFooter(f IField) {
	fmt.Println("----------------------------------")
	fmt.Printf("Status: %s\n", getPrintableGameStatus(f))
	fmt.Println("----------------------------------")
}

func getValidAnswer[T any](question string, predicate func(string) (bool, T)) T {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("")
		fmt.Printf(">>> %s", question)
		a, _, _ := reader.ReadLine()
		if isValid, res := predicate(string(a)); isValid {
			return res
		}
		fmt.Println("Input error!")
	}
}

func GameTurnUI(f IField) Point {
	printHeader()
	printField(f)
	printTurnFooter(f)

	return getValidAnswer(
		"What cell to open? (input e.g. b2, e6, j18): ",
		func(ans string) (bool, Point) {
			p := Point{}
			if len(ans) < 2 {
				return false, p
			}
			letter := ans[0]
			// convert character code into an X coordinate 
			p.X = int(byte(letter) - byte(lettersBase))
			var e error
			// convert second part - number strin into Y coordinate
			p.Y, e = strconv.Atoi(ans[1:])
			return e == nil && p.X < GameSettingsMaxWidth, p
		},
	)

}

func GameEndedUI(f IField) bool {
	printHeader()
	printField(f)
	printTurnFooter(f)

	return getValidAnswer(
		"Begin New Game? (y/n) : ",
		func(ans string) (bool, bool) {
			return true, strings.ToLower(ans) == "y"
		},
	)

}

func NewGameUI() GameSettings {
	gs := GameSettings{}
	printHeader()
	fmt.Println("*** Create New Game ****")

	gs.Height = getValidAnswer(
		fmt.Sprintf("Game field Height (%d < height < %d): ", GameSettingsMinHeight, GameSettingsMaxHeight),
		func(ans string) (bool, int) {
			res, e := strconv.Atoi(ans)
			return e == nil && res > GameSettingsMinHeight && res < GameSettingsMaxHeight, res
		},
	)

	gs.Width = getValidAnswer(
		fmt.Sprintf("Game field Width (%d < width < %d): ", GameSettingsMinWidth, GameSettingsMaxWidth),
		func(ans string) (bool, int) {
			res, e := strconv.Atoi(ans)
			return e == nil && res > GameSettingsMinWidth && res < GameSettingsMaxWidth, res
		},
	)

	gs.HolesNumber = getValidAnswer(
		fmt.Sprintf("Number of holes (0 < width < %d): ", gs.Height * gs.Width),
		func(ans string) (bool, int) {
			res, e := strconv.Atoi(ans)
			return e == nil && res > 0 && res < gs.Height * gs.Width, res
		},
	)

	fmt.Println("\n")
	
	return gs
}
