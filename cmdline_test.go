package main

import (
	"testing"
)

func TestInitCmdLine(t *testing.T) {
	InitCmdLine()
	t.Log("Initial command line OK")
}

func TestParseCmdLine(t *testing.T) {
	help = true
	result := ParseCmdLine()
	if result == -1 {
		t.Log("Parse command line for help usage")
	}
	help = false
	version = true
	result = ParseCmdLine()
	if result == -2 {
		t.Log("Parse command line for help version")
	}
	help = false
	version = false
	Version = true
	result = ParseCmdLine()
	if result == -3 {
		t.Log("Parse command line for print default")
	}
}
