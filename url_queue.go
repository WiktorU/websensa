package main

import (
	"net/url"
	"sync"
)

// URLQueue is a string Queue with locks and checking mechanism.
type URLQueue struct {
	mu       sync.Mutex
	queue    *Queue[url.URL]
	usedURLs map[url.URL]bool
}

// NewURLQueue creates a new empty queue.
func NewURLQueue() *URLQueue {
	return &URLQueue{queue: NewQueue[url.URL](), usedURLs: make(map[url.URL]bool)}
}

// Enqueue adds an element to the end of the queue if it is not already used.
func (q *URLQueue) Enqueue(url url.URL) {
	defer q.mu.Unlock()

	q.mu.Lock()
	if !q.usedURLs[url] {
		q.usedURLs[url] = true
		q.queue.Enqueue(url)
	}
}

// Dequeue removes and returns the element from the front of the queue.
func (q *URLQueue) Dequeue() (url.URL, bool) {
	defer q.mu.Unlock()

	q.mu.Lock()
	return q.queue.Dequeue()
}
