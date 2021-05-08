package main

import "testing"

func TestTempSensorConfig_ParseJsonFile(t *testing.T) {
	tempSensor := NewTempSensor()
	fileName := CheckFile("./config.json")
	err := tempSensor.ParseJsonFile(fileName)
	if err != nil {
		t.Error("Parse JSON file fail!")
	} else {
		t.Log("Parse JSON file OK")
	}
}
