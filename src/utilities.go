package main

// Simple implementation of the Stack using Single-Linked list and Mutex for concurrency - enough for the task, theoretically more efficient then using Go slices as a stack.
// Modified code from here: https://stackoverflow.com/a/40441569

import (
	"sync"
)

type element[T any] struct {
	data T
	next *element[T]
}

// Stack would be generic on elements data type, to avoid interface{}
type Stack[T any] struct {
	lock *sync.Mutex
	head *element[T]
	Size int
}

// Push() puts data on top of the stack. As stack is unlimited (list-based), we dont need error/ok values here (at least for naive implementation)
func (stk Stack[T]) Push(data T) {
	stk.lock.Lock()

	element := new(element[T])
	element.data = data
	temp := stk.head
	element.next = temp
	stk.head = element
	stk.Size++

	stk.lock.Unlock()
}

// Pop() removes and returnes data from the to pof the stack.
// Will return ok = false in case of error - just like e.g. we do for maps : el, ok := myMap["somekey"]
func (stk Stack[T]) Pop() (T, bool) {
	if stk.head == nil {
		var t T
		return t, false
	}
	stk.lock.Lock()
	r := stk.head.data
	stk.head = stk.head.next
	stk.Size--

	stk.lock.Unlock()

	return r, true
}

func NewStack[T any]() *Stack[T] {
	// this will be heap-allocated as it is passed out of the function
	return &Stack[T]{
		head: nil,
		Size: 0,
		lock: &sync.Mutex{},
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

// IField interface is an abstruction of the game field, owns states of cells
// This abstraction allows to flood-fill game fields with different logic over cells
// (logic in our case: 
//	  we need to flood-fill all clear-free cells + cells with countrs on the closure border of the free region)
// - IsFillable() - test if the cell is fillable (so can be filled with a given state by the fill())
// - Fill()  - knows how to fill the cell with a new state
// - GetWidth() int - abstract width
// - GetHeight() int - abstract height
type IField interface {
	IsFillable(x, y int) bool
	Fill(x, y int)
	GetWidth() int
	GetHeight() int
}

// Here is an algorithm rewritten in Go from here: https://lodev.org/cgtutor/floodfill.html
// - x, y int - coordinates of the seed to start from
// - field - abstract cell field
func FloodFill(x, y int, field IField) FillEvent {
	if !field.IsFillable(x, y) {
		return FillEventNothingToFill
	}

	width := field.GetWidth()
	height := field.GetHeight()

	stack := NewStack[Point]()

	// Point is a small structure, store by-value
	stack.Push(Point{x, y})

	for point, ok := stack.Pop(); ok; point, ok = stack.Pop() {
		x1 := point.X
		for x1 >= 0 && field.IsFillable(x1, y) {
			x1--
		}
		x1++
		spanAbove := false
		spanBelow := false

		for x1 < width && field.IsFillable(x1, y) {
			field.Fill(x1, y)

			if !spanAbove && y > 0 && field.IsFillable(x1, (y-1)) {
				stack.Push(Point{x1, y - 1})
				spanAbove = true
			} else if spanAbove && y > 0 && !field.IsFillable(x1, (y-1)) {
				spanAbove = false
			}

			if !spanBelow && y < height-1 && field.IsFillable(x1, (y+1)) {
				stack.Push(Point{x1, y + 1})
				spanBelow = true
			} else if spanBelow && y < height-1 && !field.IsFillable(x1, (y+1)) {
				spanBelow = false
			}
			x1++
		}
	}
	return FillEventFinished
}
