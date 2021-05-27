package main

import (
	"os"
)

func main() {
	InitCmdLine()
	result := ParseCmdLine()
	if result != 0 {
		os.Exit(0)
	}
	StartMonitorTask()
	// loop
	select {}
}
