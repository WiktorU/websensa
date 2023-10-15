package main

// Queue is a generic FIFO queue implementation.
type Queue[T any] struct {
	elements []T
}

// NewQueue creates a new empty queue.
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{elements: make([]T, 0)}
}

// Enqueue adds an element to the end of the queue.
func (q *Queue[T]) Enqueue(item T) {
	q.elements = append(q.elements, item)
}

// Dequeue removes and returns the element from the front of the queue.
func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.elements) == 0 {
		var defaultVal T

		return defaultVal, false
	}

	item := q.elements[0]
	q.elements = q.elements[1:]

	return item, true
}
