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

type UI struct {
	gameSettings *GameSettings
	f            IField
	isUnicode    bool
}

func (ui *UI) GetSettings() *GameSettings {
	return ui.gameSettings
}

func (ui *UI) SetField(f IField) {
	ui.f = f
}

func (ui *UI) printField() {
	h := ui.f.GetHeight()
	w := ui.f.GetWidth()

	closedPic := "🙫"
	holePic := "⦿"
	freePic := "⛶"
	picFmt := "%s "
	if !ui.isUnicode {
		closedPic = "##"
		holePic = "**"
		freePic = "[]"
		picFmt = "%s"
	}
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
			pc := ui.f.GetCell(x, y)
			switch {
			case pc.State == CellStateClosed:
				// add space after placeholder "%s " to fix wide rune printing issue
				fmt.Printf(picFmt, closedPic)
			case pc.HolesNumber == ThisIsHoleMarker:
				fmt.Printf(picFmt, holePic)
			case pc.HolesNumber == 0:
				fmt.Printf(picFmt, freePic)
			default:
				fmt.Printf("%d ", pc.HolesNumber)
			}
		}
		fmt.Println()
	}
	fmt.Println("\nLegend:")
	fmt.Printf("%s - closed\n", closedPic)
	fmt.Printf("%s - free\n", freePic)
	fmt.Println("1 - hole-adjacent counter")
	fmt.Printf("%s - hole\n\n", holePic)
}

func (ui *UI) getPrintableGameStatus() string {
	s := ui.f.GetState()
	switch s {
	case GameStateInProgress:
		return "In progress"
	case GameStateLoose:
		return "You Lost !"
	case GameStateWin:
		return "!!! You WON !!!"
	}
	return "Some error..."
}

func (ui *UI) printHeader() {
	fmt.Println("\n\n\n==================================")
	fmt.Println("============= Proxx ==============")
	fmt.Println("==================================")
}

func (ui *UI) printTurnFooter() {
	fmt.Println("----------------------------------")
	fmt.Printf("Status: %s\n", ui.getPrintableGameStatus())
	fmt.Println("----------------------------------")
}

func (ui *UI) GameTurn() Point {
	ui.printHeader()
	ui.printField()
	ui.printTurnFooter()

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

func (ui *UI) GameEnded() bool {
	ui.printHeader()
	ui.printField()
	ui.printTurnFooter()

	return getValidAnswer(
		"Begin New Game? (y/n) : ",
		func(ans string) (bool, bool) {
			return true, strings.ToLower(ans) == "y"
		},
	)

}

func NewUI() *UI {
	ui := new(UI)
	gs := &GameSettings{}
	ui.printHeader()
	fmt.Println("*** Create New Game ****")

	ui.isUnicode = getValidAnswer(
		"Play Unicode version (best on Ubuntu terminal)? (y/n) : ",
		func(ans string) (bool, bool) {
			return true, strings.ToLower(ans) == "y"
		},
	)

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
		fmt.Sprintf("Number of holes (0 < width < %d): ", gs.Height*gs.Width),
		func(ans string) (bool, int) {
			res, e := strconv.Atoi(ans)
			return e == nil && res > 0 && res < gs.Height*gs.Width, res
		},
	)

	fmt.Println("\n")

	ui.gameSettings = gs

	return ui
}
