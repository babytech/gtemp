package main

import "testing"

func TestStartMonitorTask(t *testing.T) {
	dummyTemp = true
	StartMonitorTask()
	t.Log("Start tempSensor monitor task OK")
}
