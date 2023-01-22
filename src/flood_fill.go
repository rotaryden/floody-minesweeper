package main

type FillEvent int

const (
	FillEventFinished = iota
	FillEventNothingToFill
)

// -----------------------------------------------------------------------------------------------------------
// Here is a non-recursive stack-based scan line algorithm
// rewritten in Go from here: https://lodev.org/cgtutor/floodfill.html
// it has been chosen because of speed, program stack safety and relative simplicity
// - x, y int - coordinates of the seed to start from
// - field - abstract cell field
func FloodFill(x, y int, field IFloodableField) FillEvent {
	isFirstCell := true

	if !field.IsFillable(x, y, isFirstCell) {
		return FillEventNothingToFill
	}

	width := field.GetWidth()
	height := field.GetHeight()

	stack := NewStack[Point]()

	// Point is a small structure without side effects and references, store by-value
	stack.Push(Point{x, y})

	for point, ok := stack.Pop(); ok; point, ok = stack.Pop() {
		dx := point.X
		lineY := point.Y

		// move to the very left cell of the current connected region on the row
		for dx >= 0 && field.IsFillable(dx, lineY, isFirstCell) {
			dx--
		}
		dx++

		spanAbove := false
		spanBelow := false

		for dx < width && field.IsFillable(dx, lineY, isFirstCell) {
			field.Fill(dx, lineY)
			isFirstCell = false

			// check adjucent cells from the row above and push to stack
			if !spanAbove && lineY > 0 && field.IsFillable(dx, (lineY-1), isFirstCell) {
				stack.Push(Point{dx, lineY - 1})
				spanAbove = true
			} else if spanAbove && lineY > 0 && !field.IsFillable(dx, (lineY-1), isFirstCell) {
				// reject to span above into the above adjacent cell close to the previous adjacent cell, 
				// where we''ve already spawn
				// but will step into the next above cell on the right, if needed
				spanAbove = false
			}

			// check adjucent cells from the row below and push to stack
			if !spanBelow && lineY < height-1 && field.IsFillable(dx, (lineY+1), isFirstCell) {
				stack.Push(Point{dx, lineY + 1})
				spanBelow = true
			} else if spanBelow && lineY < height-1 && !field.IsFillable(dx, (lineY+1), isFirstCell) {
				spanBelow = false
			}

			// move on the current line
			dx++
		}
	}
	return FillEventFinished
}
