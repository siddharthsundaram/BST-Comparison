package main

import "sync"

type ConcurrentBuffer struct {
	first *LLNode
	last *LLNode
	buf_lock sync.Mutex
	produce *sync.Cond
	consume *sync.Cond
	length int
	cap int
	done bool
}

type LLNode struct {
	next *LLNode
	prev *LLNode
	val Work
}

type Work struct {
	idx1 int
	idx2 int
}

func NewConcurrentBuffer(cap int) *ConcurrentBuffer {
	b := &ConcurrentBuffer{length: 0, cap: cap, done: false}
	b.produce = sync.NewCond(&b.buf_lock)
	b.consume = sync.NewCond(&b.buf_lock)

	return b
}

func (buf *ConcurrentBuffer) Length() int {
	return buf.length
}

func (buf *ConcurrentBuffer) Cap() int {
	return buf.cap
}

func (buf *ConcurrentBuffer) enqueue(idx1 int, idx2 int) {
	new_work := Work{idx1: idx1, idx2: idx2}
	new_node := &LLNode{val: new_work}

	buf.buf_lock.Lock()
	for buf.Length() == buf.Cap() {
		buf.produce.Wait()
	}

	// Enter critical section and add work to buffer
	if buf.Length() == 0 {
		buf.first = new_node
	} else {
		buf.last.next = new_node
		new_node.prev = buf.last
	}

	buf.last = new_node
	buf.length += 1

	// Signal waiting goroutine to consume
	buf.consume.Signal()
	buf.buf_lock.Unlock()
}

func (buf *ConcurrentBuffer) dequeue() (Work, bool) {
	buf.buf_lock.Lock()
	defer buf.buf_lock.Unlock()

	for buf.Length() == 0 {
		if buf.done {
			return Work{}, false
		}

		buf.consume.Wait()
	}

	// Enter critical section and take work from buffer
	res := buf.first.val
	buf.first = buf.first.next
	buf.length -= 1
	if buf.length > 0 {
		buf.first.prev = nil
	}

	// Signal waiting main thread to produce
	buf.produce.Signal()

	return res, true
}

func (buf *ConcurrentBuffer) Close() {
	buf.done = true
	buf.consume.Broadcast()
}