package main

import (
	"os"
	"testing"
)

func TestUnitTestWriteCsv(t *testing.T) {
	fileName := "/tmp/temp/test.csv"
	UnitTestWriteCsv(fileName)
	file, err := os.Open(fileName)
	if err != nil {
		t.Error("Unit Test for write CSV file fail!")
	} else {
		defer file.Close()
	}
	t.Log("Unit Test for write CSV file OK")
}

func TestUnitTestReadCsv(t *testing.T) {
	fileName := "/tmp/temp/test.csv"
	UnitTestReadCsv(fileName)
	t.Log("Unit Test for read CSV file OK")
}

func TestTempSensorConfig_HandleCsvFile(t *testing.T) {
	fileName := "/tmp/temp/test.csv"
	tempSensor := NewTempSensor()
	err := tempSensor.ParseJsonFile("/tmp/temp/config.json")
	if err != nil {
		t.Error("Test Configure File: Fail!")
	}
	readCsvFile = true
	result := tempSensor.HandleCsvFile(fileName)
	if result == 1 {
		t.Log("Read CSV file OK")
	} else {
		t.Error("Read CSV file fail")
	}
	readCsvFile = false
	writeCsvFile = true
	result = tempSensor.HandleCsvFile(fileName)
	if result == 2 {
		t.Log("Write CSV file OK")
	} else {
		t.Error("Write CSV file fail")
	}
	readCsvFile = false
	writeCsvFile = false
	result = tempSensor.HandleCsvFile(fileName)
	if result == 0 {
		t.Log("Generate CSV file OK")
	}
}