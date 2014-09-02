package main

import (
	"container/heap"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"runtime"
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

type selectorFunc func(numPerson int, numIterations int) (int, []roommates)

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
	runtime.GOMAXPROCS(4)
	fmt.Println("GOMAXPROCS is ", runtime.GOMAXPROCS(0))
	rand.Seed(time.Now().Unix() * 123)
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
	numPeople := 1000
	initPref(numPeople)

	var selectors = []selectorFunc{simpleSelector, iterateRoommateChoicesSelector, geneticSelector}
	selector := selectors[2]
	cost, assignment := selector(numPeople, 10000)

	fmt.Println("Final Choice: ", cost, "\n", assignment)

}

// input assignment
func getCostAndAssignment(slotAssignment []int) (int, []roommates) {
	roomAssignment := slotToRoom(slotAssignment)
	cost := getRoommateCost(roomAssignment)
	return cost, roomAssignment
}

func simpleSelector(numPerson int, numIterations int) (int, []roommates) {
	lowestCost := 1000000
	var finalAssignment []roommates

	for i := 0; i < numIterations; i++ {
		slotAssignment := getRandomChoice(numPerson)
		cost, assignment := getCostAndAssignment(slotAssignment)
		if cost < lowestCost {
			lowestCost = cost
			finalAssignment = assignment
		}
		if i%100 == 0 {
			logger.Println("cost:  ", i, cost, lowestCost)
		}
	}
	return lowestCost, finalAssignment
}

func geneticSelector(numPerson int, numIterations int) (int, []roommates) {
	numVariations := 100
	choices := make([][]int, numVariations)
	for i := 0; i < numVariations; i++ {
		choices[i] = getRandomChoice(numPerson)
	}
	bestScore, numTopSolutions, selectedChoices := evaluateGeneration(numVariations, choices)

	for i := 1; i < numIterations; i++ {
		//		fmt.Println("Generation ", i)
		choices = getNextGeneration(numPerson, numTopSolutions, numVariations, selectedChoices)
		//		fmt.Println(choices[0])
		bestScore, numTopSolutions, selectedChoices = evaluateGeneration(numVariations, choices)
	}
	bestSolution := slotToRoom(selectedChoices[numTopSolutions-1])
	//	fmt.Println("Genetic: \n", bestScore, bestSolution)
	return bestScore, bestSolution
}

func getNextGeneration(numPerson, numTopSolutions int, numVariations int, selectedChoices [][]int) [][]int {

	//	fmt.Println("GetNextGeneration: ", numTopSolutions, numVariations)
	result := make([][]int, numVariations)
	curIndex := numVariations - 1

	for i := 0; i < numTopSolutions; i++ {
		result[curIndex] = make([]int, numPerson)
		copy(result[curIndex], selectedChoices[numTopSolutions-i-1])
		//		result[curIndex] = selectedChoices[numTopSolutions-i-1]
		// apply a mutation here
		if i > 1 {
			randomBit := rand.Intn(numPerson)
			//			fmt.Print(i, curIndex, randomBit, result[curIndex])
			result[curIndex][randomBit] = rand.Intn(numPerson - randomBit)
			//			fmt.Println(result[curIndex])
		}
		curIndex--
	}
	for i := 0; i < numTopSolutions; i++ {
		for j := 0; j < numTopSolutions; j++ {
			if i == j {
				continue
			} else {
				cutLocation := rand.Intn(numPerson-4) + 1
				//				fmt.Println("Joining: ", cutLocation, i, j)
				//				fmt.Println(selectedChoices[i], selectedChoices[j])
				result[curIndex] = make([]int, numPerson)
				copy(result[curIndex], selectedChoices[i])

				result[curIndex] = append(result[curIndex][:cutLocation], selectedChoices[j][cutLocation:]...)
				// apply point mutation here
				pointMutation := rand.Intn(numPerson)
				result[curIndex][pointMutation] = rand.Intn(numPerson - pointMutation)
				//				fmt.Println(result[curIndex])
				curIndex--
			}
		}
	}
	for i := curIndex; i >= 0; i-- {
		cutLocation := rand.Intn(numPerson-4) + 2
		result[i] = make([]int, numPerson)
		copy(result[i], selectedChoices[rand.Intn(numVariations)])
		result[i] = append(result[i][:cutLocation+1], selectedChoices[rand.Intn(numVariations)][cutLocation+1:]...)
		// apply point mutation here
		pointMutation := rand.Intn(numPerson)
		result[i][pointMutation] = rand.Intn(numPerson - pointMutation)

	}
	//	for i := 0; i < numVariations; i++ {
	//		fmt.Println(i, result[i])
	//	}
	return result
}

func evaluateGeneration(numVariations int, choices [][]int) (int, int, [][]int) {
	numTopSolutions := int(math.Sqrt(float64(numVariations)))
	if numTopSolutions < 10 {
		numTopSolutions = 10
	}
	pq := make(CappedPriorityQueue, 0, numTopSolutions+1)
	heap.Init(&pq)
	lowestCost := 9999999
	var cost int
	//var solution []roommates
	//	fmt.Println("EvaluateGeneration")
	for i := 0; i < numVariations; i++ {
		choice := choices[i]
		cost, _ = getCostAndAssignment(choice)
		//		fmt.Println(i, cost, choice)
		if cost < lowestCost {
			lowestCost = cost
		}
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
			//			fmt.Println(index, cost, topChoices[index])
		}
		index++
	}
	//	solution = slotToRoom(topChoices[numTopSolutions-1])
	//	fmt.Println("genetic: ", cost, lowestCost)
	//fmt.Println("topChoice: ", topChoices[numTopSolutions-1])
	return cost, numTopSolutions, topChoices
}

func iterateRoommateChoicesSelector(numPerson int, numIterations int) (int, []roommates) {
	// numIterations is ignored since we will iterate through implicitly
	slotAssignment := getRandomChoice(numPerson)
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
