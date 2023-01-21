package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const letters = "abcdefghijklmnopqrstuvwxyz"
const lettersBase = int('a')

func getPrintableField(f IField) []string {
	h := f.GetHeight()
	w := f.GetWidth()
	result := make([]string, h+1)
	// letter denotion of columns
	result[0] = letters[0:w]
	for y := 0; y < h; y++ {
		var sb strings.Builder
		// print line number
		sb.WriteString(fmt.Sprintf("%2d ", y))
		for x := 0; x < w; x++ {
			pc := f.GetCell(x, y)
			switch {
			case pc.State == CellStateClosed:
				sb.WriteRune('🙫')
			case pc.HolesNumber == ThisIsHoleMarker:
				sb.WriteRune('⦿')
			case pc.HolesNumber == 0:
				sb.WriteRune('⛶')
			default:
				sb.WriteRune(rune(pc.HolesNumber))
			}
		}
		result[y+1] = sb.String()
	}
	return result
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
	fmt.Println("==================================")
	fmt.Println("============= Proxx ==============")
	fmt.Println("==================================")
}

func printFooter(f IField) {
	fmt.Println("----------------------------------")
	fmt.Printf("Status: %s\n", getPrintableGameStatus(f))
	fmt.Println("----------------------------------\n\n")
}

func getValidAnswer[T any](reader *bufio.Reader, question string, predicate func(string) (bool, T)) T {
	for {
		fmt.Println("")
		fmt.Print(question)
		a, _, _ := reader.ReadLine()
		if isValid, res := predicate(string(a)); isValid {
			return res
		}
		fmt.Println("Input error!")
	}
}

func GameTurnUI(f IField) Point {
	printHeader()
	lines := getPrintableField(f)
	for l := range lines {
		fmt.Println(l)
	}
	printFooter(f)
	r := bufio.NewReader(os.Stdin)

	return getValidAnswer(
		r,
		"Enter cell to open in form of dd-l, where d-> 0...9 and l-> 'a'...'z': ",
		func(ans string) (bool, Point) {
			parts := strings.Split(ans, "-")
			p := Point{}
			if len(parts[1]) > 1 {
				return false, p
			}
			// convert character code into an X coordinate 
			p.X = int(byte(parts[1][0]) - byte(lettersBase))
			var e error
			// convert first part - number strin into Y coordinate
			p.Y, e = strconv.Atoi(parts[0])
			return e == nil && p.X < GameSettingsMaxWidth, p
		},
	)

}

func NewGameUI() GameSettings {
	gs := GameSettings{}
	printHeader()
	fmt.Println("*** Create New Game ****")
	r := bufio.NewReader(os.Stdin)

	gs.Height = getValidAnswer(
		r,
		fmt.Sprintf("Game field Height (%d < height < %d): ", GameSettingsMinHeight, GameSettingsMaxHeight),
		func(ans string) (bool, int) {
			res, e := strconv.Atoi(ans)
			return e == nil && res > GameSettingsMinHeight && res < GameSettingsMaxHeight, res
		},
	)

	gs.Width = getValidAnswer(
		r,
		fmt.Sprintf("Game field Width (%d < width < %d): ", GameSettingsMinWidth, GameSettingsMaxWidth),
		func(ans string) (bool, int) {
			res, e := strconv.Atoi(ans)
			return e == nil && res > GameSettingsMinWidth && res < GameSettingsMaxWidth, res
		},
	)

	gs.HolesNumber = getValidAnswer(
		r,
		fmt.Sprintf("Number of holes (0 < width < %d): ", gs.Height * gs.Width),
		func(ans string) (bool, int) {
			res, e := strconv.Atoi(ans)
			return e == nil && res > 0 && res < gs.Height * gs.Width, res
		},
	)

	fmt.Println("\n")
	
	return gs
}
