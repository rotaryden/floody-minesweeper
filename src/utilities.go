package main

import "sync"

// ----------------------------------------------------------------------------------
// Simple implementation of the Stack using Single-Linked list and Mutex for concurrency -
// enough for the sake of the task,
// likely more efficient then using Go slices as a backing structure.
// Modified code from here: https://stackoverflow.com/a/40441569

type elementType[T any] struct {
	data T
	next *elementType[T]
}

// Stack would be generic on elements data type, to avoid interface{}
type Stack[T any] struct {
	lock sync.Mutex
	head *elementType[T]
	Size int
}

// Push() puts data on top of the stack. As stack is unlimited (list-based),
// we dont need error/ok checks (at least for naive implementation)
func (stk *Stack[T]) Push(data T) {
	stk.lock.Lock()
	defer stk.lock.Unlock()

	el := &elementType[T]{data: data}
	temp := stk.head
	el.next = temp
	stk.head = el
	stk.Size++

}

// Pop() removes and returnes data from the top of the stack.
// Will return ok = false in case of emptiness -
// same approach we do e.g. for maps: element, ok := myMap["somekey"]
func (stk *Stack[T]) Pop() (T, bool) {
	if stk.head == nil {
		var t T
		return t, false
	}
	stk.lock.Lock()
	defer stk.lock.Unlock()

	r := stk.head.data
	stk.head = stk.head.next
	stk.Size--

	return r, true
}

func NewStack[T any]() *Stack[T] {
	// &Stack{} will be heap-allocated as the reference is passed out of the function
	return &Stack[T]{
		head: nil,
		Size: 0,
		lock: sync.Mutex{},
	}
}

// Point abstraction
type Point struct {
	X int
	Y int
}

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
	if !field.IsFillable(x, y) {
		return FillEventNothingToFill
	}

	width := field.GetWidth()
	height := field.GetHeight()

	stack := NewStack[Point]()

	// Point is a small structure without side effects, store by-value
	stack.Push(Point{x, y})

	for point, ok := stack.Pop(); ok; point, ok = stack.Pop() {
		dx := point.X
		lineY := point.Y

		for dx >= 0 && field.IsFillable(dx, lineY) {
			dx--
		}
		dx++

		spanAbove := false
		spanBelow := false

		for dx < width && field.IsFillable(dx, lineY) {
			field.Fill(dx, lineY)

			if !spanAbove && lineY > 0 && field.IsFillable(dx, (lineY-1)) {
				stack.Push(Point{dx, lineY - 1})
				spanAbove = true
			} else if spanAbove && lineY > 0 && !field.IsFillable(dx, (lineY-1)) {
				spanAbove = false
			}

			if !spanBelow && lineY < height-1 && field.IsFillable(dx, (lineY+1)) {
				stack.Push(Point{dx, lineY + 1})
				spanBelow = true
			} else if spanBelow && lineY < height-1 && !field.IsFillable(dx, (lineY+1)) {
				spanBelow = false
			}
			dx++
		}
	}
	return FillEventFinished
}
