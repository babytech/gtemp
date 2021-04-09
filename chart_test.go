package main

import (
	"os/exec"
	"testing"
)

func TestChartMain(t *testing.T) {
	ChartMain()
	cmdStr := `
#!/bin/bash
for var in {1..2}
do
	sleep 1
	curl http://127.0.0.1:4321/MF14/temp
done`
	cmd := exec.Command("bash", "-c", cmdStr+" >> file.log")
	err := cmd.Start()
	if err != nil {
		t.Error("Command start fail")
	}
	err = cmd.Wait()
	if err != nil {
		t.Error("Command wait fail")
	}
	t.Log("curl run OK")
}