package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"testing"
)

func TestCappedPriorityQueue(t *testing.T) {
	pq := make(CappedPriorityQueue, 0, 5)
	fmt.Println("Init...", pq)
	heap.Init(&pq)
	for i := 0; i < 10; i++ {
		randNum := rand.Intn(200) - 100
		item := &Item{
			value:    randNum,
			priority: randNum,
		}
		fmt.Println("Pushing ", item.value)
		heap.Push(&pq, item)
	}
	count := 0
	for {
		result := heap.Pop(&pq)
		if result == nil {
			break
		} else {
			data := result.(*Item)
			fmt.Println("priority: ", data.priority, ", value: ", data.value)
			count++
		}
	}
	if count != 5 {
		t.Error("Expected size ", 5, ", got ", count)
	}
}
