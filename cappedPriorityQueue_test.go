package main

import (
	"container/heap"
	//	"fmt"
	"math/rand"
	"testing"
)

func TestCappedPriorityQueue(t *testing.T) {
	pq := make(CappedPriorityQueue, 0, 5+1)
	heap.Init(&pq)
	for i := 0; i < 100; i++ {
		randNum := rand.Intn(200) - 100
		item := &Item{
			value:    randNum,
			priority: randNum,
		}
		heap.Push(&pq, item)
		if len(pq) > 5 {
			heap.Pop(&pq)
		}
	}
	count := 0
	lastPriority := 99999999
	for {
		result := heap.Pop(&pq)
		if result == nil {
			break
		} else {
			data := result.(*Item)
			//			fmt.Println("priority: ", data.priority, ", value: ", data.value)
			if data.priority > lastPriority {
				t.Error("Priority out of order: ", data.priority, lastPriority)
			} else {
				lastPriority = data.priority
			}
			count++
		}
	}
	if count != 5 {
		t.Error("Expected size ", 5, ", got ", count)
	}
}
