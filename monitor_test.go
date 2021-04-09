package main

import "testing"

func TestStartMonitorTask(t *testing.T) {
	d = true
	StartMonitorTask()
	t.Log("Start tempSensor monitor task OK")
}