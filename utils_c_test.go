package main

import (
	"bytes"
	"testing"
)

func TestMemCopy(t *testing.T) {
	srcData := []byte("Hello world")
	dstData := make([]byte, len(srcData))
	MemCopy(dstData, srcData)
	t.Logf("dest data: %s\n", string(dstData))
	if bytes.Equal(srcData, dstData) {
		t.Log("Memory copy OK")
	} else {
		t.Error("Memory copy fail")
	}
}

func TestMemMove(t *testing.T) {
	srcData := []byte("Hello World")
	dstData := make([]byte, 16)
	dstData = append(dstData, srcData...)
	t.Logf("[BEFORE] Memory Move --- dest data: %s\n", string(dstData))
	MemMove(dstData[0:], dstData[16:])
	t.Logf("[AFTER] Memory Move --- dest data: %s\n", string(dstData))
	if bytes.Equal(srcData, dstData[0:len(srcData)]) {
		t.Log("Memory copy OK")
	} else {
		t.Error("Memory copy fail")
	}
}
