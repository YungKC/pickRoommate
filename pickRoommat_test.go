package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
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

func TestGeneticSelector(t *testing.T) {
	testSelector(t, geneticSelector, 1000)
}

func TestSimpleSelector(t *testing.T) {
	testSelector(t, simpleSelector, 10000)
}

func TestIterativeSelector(t *testing.T) {
	testSelector(t, iterateRoommateChoicesSelector, 1)
}

func testSelector(t *testing.T, selector selectorFunc, numIterations int) {
	numPeople := 200
	initPref(numPeople)

	_, roomData := selector(numPeople, numIterations)

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

func Benchmark100KRandomSelector(b *testing.B) {
	benchmarkSelector(b, 100000, simpleSelector)
}

func BenchmarkDeepDiveSelector(b *testing.B) {
	benchmarkSelector(b, 1, iterateRoommateChoicesSelector)
}

func Benchmark1KGeneticSelector(b *testing.B) {
	benchmarkSelector(b, 1000, geneticSelector)
}

func benchmarkSelector(b *testing.B, iterations int, selector selectorFunc) {
	b.StopTimer()
	numPeople := 100
	initPref(numPeople)
	b.StartTimer()
	var totalCost int
	var assignment []roommates
	for loopCount := 0; loopCount < b.N; loopCount++ {
		cost, result := selector(numPeople, iterations)
		assignment = result
		totalCost = totalCost + cost
	}
	fmt.Println("\n Average Cost: ", totalCost/b.N, "\n", assignment)
}
