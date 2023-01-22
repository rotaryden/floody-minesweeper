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
		// mutex is needed for potential concurrent usage
		lock: sync.Mutex{}, 
	}
}

