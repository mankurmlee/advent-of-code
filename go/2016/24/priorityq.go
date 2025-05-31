package main

import "container/heap"

type PriorityQueueItem[T any] struct {
	value    T
	priority int
}

type PriorityQueue[T any] []*PriorityQueueItem[T]

func NewPriorityQueue[T any]() PriorityQueue[T] { return PriorityQueue[T]{} }
func (pq PriorityQueue[T]) Len() int            { return len(pq) }
func (pq PriorityQueue[T]) Less(i, j int) bool  { return pq[i].priority < pq[j].priority }
func (pq PriorityQueue[T]) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue[T]) Push(x any)         { *pq = append(*pq, x.(*PriorityQueueItem[T])) }
func (pq *PriorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return item
}

func (pq *PriorityQueue[T]) Enqueue(elem T, priority int) {
	heap.Push(pq, &PriorityQueueItem[T]{elem, priority})
}

func (pq *PriorityQueue[T]) Dequeue() T {
	item := heap.Pop(pq).(*PriorityQueueItem[T])
	return item.value
}
