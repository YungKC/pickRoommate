package main

import (
	"container/heap"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

/*
n person
2 person dorm rooms
match room-mates to satisfy most requests

p1, p2, ..., pn
slot1, slot2, ..., slotn
slot1 and slot2 -> room1

*/

type roommates [2]int
type selectorFunc func(slotAssignment []int) (int, []roommates)

var pref = map[roommates]int{
	roommates{2, 3}: 1,
	roommates{5, 7}: -1,
	roommates{0, 1}: -1,
	roommates{0, 2}: -1,
	roommates{0, 3}: -1,
	roommates{0, 4}: -1,
	roommates{0, 5}: -1,
	roommates{0, 6}: -1,
	roommates{0, 7}: -1,
	roommates{0, 8}: 1,
	roommates{0, 9}: 1,
}

var logger = log.New(ioutil.Discard, "Log: ", log.Ltime|log.Lshortfile)

func initPref(count int) {
	var weight int
	for i := 0; i < count; i++ {
		for j := i + 1; j < count; j++ {
			lastDigit := (i + j) % 10
			if lastDigit == 8 {
				weight = -100
			} else if lastDigit == 9 {
				weight = -50
			} else if lastDigit == 4 {
				weight = 1000
			} else {
				weight = 0
			}
			if weight != 0 {
				pref[roommates{i, j}] = weight
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().Unix() * 123)
	numPeople := 100
	initPref(numPeople)
	iterations := 1

	lowestCost := 1000000
	var finalAssignment []roommates
	var finalChoice []int

	for i := 0; i < iterations; i++ {
		choice := getRandomChoice(numPeople)
		logger.Println("choice:", choice)

		selector := iterateRoommateChoicesSelector
		//selector := simpleSelector

		cost, assignment := selector(choice)
		if cost < lowestCost {
			lowestCost = cost
			finalChoice = choice
			finalAssignment = assignment
		}
		if i%100 == 0 {
			logger.Println("cost:  ", i, cost, lowestCost)
		}
	}
	fmt.Println("Final Choice: ", lowestCost, "\n", finalChoice, "\n", finalAssignment)

}

// input assignment
func getCostAndAssignment(slotAssignment []int) (int, []roommates) {
	roomAssignment := slotToRoom(slotAssignment)
	cost := getRoommateCost(roomAssignment)
	return cost, roomAssignment
}

func simpleSelector(slotAssignment []int) (int, []roommates) {
	return getCostAndAssignment(slotAssignment)
}

func geneticSelector(numPerson int) (int, []roommates) {
	numVariations := 10000
	choices := make([][]int, numVariations)
	for i := 0; i < numVariations; i++ {
		choices[i] = getRandomChoice(numPerson)
	}
	bestScore, numTopSolutions, selectedChoices := evaluateGeneration(numVariations, choices)

	for i := 1; i < 1000; i++ {
		fmt.Println("Generation ", i)
		choices = getNextGeneration(numPerson, numTopSolutions, numVariations, selectedChoices)
		fmt.Println(choices[0])
		bestScore, numTopSolutions, selectedChoices = evaluateGeneration(numVariations, choices)
	}
	return bestScore, slotToRoom(selectedChoices[numTopSolutions-1])
}

func getNextGeneration(numPerson, numTopSolutions int, numVariations int, selectedChoices [][]int) [][]int {
	// keep the last best 10
	result := make([][]int, numVariations)
	for i := 0; i < 10; i++ {
		result[i] = selectedChoices[numTopSolutions-i-1]
	}
	currentCount := 10
	for i := 0; i < numTopSolutions; i++ {
		for j := 0; j < numTopSolutions && currentCount < numVariations; j++ {
			if i == j {
				continue
			} else {
				cutLocation := rand.Intn(numPerson)
				if cutLocation < 1 {
					cutLocation = 1
				}
				//				fmt.Println("Joining: ", currentCount, cutLocation, i, j)
				result[currentCount] = append(selectedChoices[i][:cutLocation-1], selectedChoices[j][cutLocation:]...)
				currentCount++
			}
		}
	}
	for i := currentCount; i < numVariations; i++ {
		cutLocation := rand.Intn(numPerson)
		result[i] = append(selectedChoices[rand.Intn(numTopSolutions)][:cutLocation-1], selectedChoices[rand.Intn(numTopSolutions)][cutLocation:]...)
	}
	return result
}

func evaluateGeneration(numVariations int, choices [][]int) (bestScore int, numReturned int, selectedChoices [][]int) {
	numTopSolutions := 100
	pq := make(CappedPriorityQueue, 0, numTopSolutions+1)
	heap.Init(&pq)
	var cost int
	var solution []roommates
	for i := 0; i < numVariations; i++ {
		choice := choices[i]
		cost, solution = getCostAndAssignment(choice)
		item := &Item{
			value:    choice,
			priority: cost,
		}
		heap.Push(&pq, item)
		// keep the top 100 solutions
		if len(pq) > numTopSolutions {
			heap.Pop(&pq)
		}
	}

	topChoices := make([][]int, numTopSolutions)
	index := 0
	for {
		result := heap.Pop(&pq)
		if result == nil {
			break
		} else {
			data := result.(*Item)
			cost = data.priority
			topChoices[index] = data.value.([]int)
			//			fmt.Println(cost)
		}
		index++
	}
	solution = slotToRoom(topChoices[numTopSolutions-1])
	fmt.Println("genetic: ", cost, solution)
	return cost, numTopSolutions, topChoices
}

func iterateRoommateChoicesSelector(slotAssignment []int) (int, []roommates) {
	bestCost, bestRoommateAssignment := getCostAndAssignment(slotAssignment)
	workingRoommateAssignment := make([]roommates, len(bestRoommateAssignment))

	copy(workingRoommateAssignment, bestRoommateAssignment)
	logger.Println("starting iteration: ", workingRoommateAssignment)
	for i := 1; i < len(slotAssignment); i++ {
		room1 := i / 2
		order1 := i % 2
		person := workingRoommateAssignment[room1][order1]
		switchedRoom := -1
		switchedOrder := -1
		for j := i + 1; j < len(slotAssignment); j++ {
			room2 := j / 2
			if room1 == room2 {
				continue
			}
			order2 := j % 2
			candidate := workingRoommateAssignment[room2][order2]
			workingRoommateAssignment[room1][order1] = workingRoommateAssignment[room2][order2]
			workingRoommateAssignment[room2][order2] = person
			curCost := getRoommateCost(workingRoommateAssignment)
			if curCost <= bestCost {
				switchedRoom = room2
				switchedOrder = order2
				bestCost = curCost
			}
			workingRoommateAssignment[room2][order2] = candidate
		}
		if switchedRoom >= 0 {
			tmp := workingRoommateAssignment[switchedRoom][switchedOrder]
			workingRoommateAssignment[switchedRoom][switchedOrder] = person
			workingRoommateAssignment[room1][order1] = tmp
			bestRoommateAssignment[switchedRoom][switchedOrder] = person
			bestRoommateAssignment[room1][order1] = tmp
		} else {
			workingRoommateAssignment[room1][order1] = person
		}
	}

	return bestCost, bestRoommateAssignment
}

func getRoommateCost(roomAssignment []roommates) int {
	result := 0
	for _, match := range roomAssignment {
		min := match[0]
		max := match[1]
		if match[1] < min {
			min = match[1]
			max = match[0]
		}
		searchMatch := roommates{min, max}
		result += pref[searchMatch]
	}
	return result
}

func slotToRoom(slotAssignment []int) []roommates {
	length := len(slotAssignment)
	result := make([]roommates, (length+1)/2)
	slots := make([]int, length)
	for i := 0; i < length; i++ {
		slots[i] = i
	}
	for i, v := range slotAssignment {
		curSlot := slots[v]
		result[curSlot/2][(curSlot+1)%2] = i
		if v == 0 {
			slots = slots[1:]
		} else if v == len(slots) {
			slots = slots[0 : len(slots)-1]
		} else {
			slots = append(slots[0:v], slots[v+1:]...)
		}

	}
	return result
}

func getRandomChoice(count int) []int {
	choice := make([]int, count)
	for i := 0; i < count; i++ {
		tmp := rand.Int() % (count - i)
		choice[i] = tmp
	}
	return choice
}
