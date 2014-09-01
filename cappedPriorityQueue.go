package main

// An Item is something we manage in a priority queue.
type Item struct {
	value    interface{} // The value of the item; arbitrary.
	priority int         // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A CappedPriorityQueue implements heap.Interface and holds Items.
type CappedPriorityQueue []*Item

func (pq CappedPriorityQueue) Len() int { return len(pq) }

func (pq CappedPriorityQueue) Less(i, j int) bool {
	//	fmt.Println("Less: ", i, j, pq[i].priority, pq[j].priority)
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq CappedPriorityQueue) Swap(i, j int) {
	//	fmt.Println("Swap: ", i, j)
	if i < 0 || j < 0 {
		return
	}
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *CappedPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *CappedPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	if n == 0 {
		return nil
	}
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
