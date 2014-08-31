package main

import (
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
