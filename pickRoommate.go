package main

import (
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
