package main

import (
	"testing"
)

func TestInitCmdLine(t *testing.T) {
	InitCmdLine()
	t.Log("Initial command line OK")
}

func TestParseCmdLine(t *testing.T) {
	h = true
	result := ParseCmdLine()
	if result == -1 {
		t.Log("Parse command line for help usage")
	}
	h = false
	v = true
	result = ParseCmdLine()
	if result == -2 {
		t.Log("Parse command line for help version")
	}
	h = false
	v = false
	V = true
	result = ParseCmdLine()
	if result == -3 {
		t.Log("Parse command line for print default")
	}
}
