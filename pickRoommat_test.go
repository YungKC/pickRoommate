package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"testing"
	"time"
)

var testLogger = log.New(ioutil.Discard, "Log: ", log.Ltime|log.Lshortfile)

func TestRandomChoice(t *testing.T) {
	testData := getRandomChoice(10)
	for i := 0; i < 10; i++ {
		if testData[i] > 9-i {
			t.Error("Expected less than ", 9-i, ", got ", testData[i])
		}
	}
}

func TestSlotToRoom(t *testing.T) {
	count := 100
	slotData := getRandomChoice(count)
	testLogger.Println(slotData)
	roomData := slotToRoom(slotData)
	testLogger.Println(roomData)
	checkData := make([]int, count)
	for _, roommates := range roomData {
		for _, index := range roommates {
			checkData[index]++
		}
	}
	testLogger.Println(checkData)
	for i, valueToCheck := range checkData {
		if valueToCheck != 1 {
			t.Error("Expected 1 at ", i, ", got ", valueToCheck)
		}
	}
}

func TestSimpleSelector(t *testing.T) {
	testSelector(t, simpleSelector)
}

func TestIterativeSelector(t *testing.T) {
	testSelector(t, iterateRoommateChoicesSelector)
}

func testSelector(t *testing.T, selector func(slotAssignment []int) (int, []roommates)) {
	rand.Seed(time.Now().Unix() * 123)
	numPeople := 200
	initPref(numPeople)

	choice := getRandomChoice(numPeople)
	_, roomData := selector(choice)

	checkData := make([]int, numPeople)
	for _, roommates := range roomData {
		for _, index := range roommates {
			checkData[index]++
		}
	}
	testLogger.Println(checkData)
	for i, valueToCheck := range checkData {
		if valueToCheck != 1 {
			t.Error("Expected 1 at ", i, ", got ", valueToCheck)
		}
	}

}

func Benchmark1000RandomSelector(b *testing.B) {
	benchmark(b, 1000, simpleSelector)
}

func BenchmarkDeepDiveSelector(b *testing.B) {
	benchmark(b, 1, iterateRoommateChoicesSelector)
}

func benchmark(b *testing.B, iterations int, selector func(slotAssignment []int) (int, []roommates)) {
	b.StopTimer()

	rand.Seed(time.Now().Unix() * 123)
	numPeople := 1000
	initPref(numPeople)

	lowestCost := 1000000
	var finalAssignment []roommates
	var finalChoice []int

	b.StartTimer()

	for loopCount := 0; loopCount < b.N; loopCount++ {
		for i := 0; i < iterations; i++ {
			choice := getRandomChoice(numPeople)

			cost, assignment := selector(choice)
			if cost < lowestCost {
				lowestCost = cost
				finalChoice = choice
				finalAssignment = assignment
			}
		}
		fmt.Println(lowestCost, finalChoice[0], finalAssignment[0])
	}
}
