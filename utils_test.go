package main

import (
	"bytes"
	"os"
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

func TestIntToBytes(t *testing.T) {
	var inputValue uint32
	inputValue = 12345678
	bytes := IntToBytes(inputValue)
	t.Logf("input integer: %d\n", inputValue)
	outputValue := BytesToInt(bytes)
	t.Logf("output integer: %d\n", outputValue)
	if outputValue == inputValue {
		t.Log("Convert integer to bytes OK")
	} else {
		t.Error("Convert integer to bytes fail")
	}
}

func TestBytesToInt(t *testing.T) {
	inputString := "BABY"
	value := BytesToInt([]byte(inputString))
	t.Logf("input string: %s\n", inputString)
	t.Logf("integer value: %d\n", value)
	outputString := string(IntToBytes(value))
	t.Logf("output string: %s\n", outputString)
	if outputString == inputString {
		t.Log("Convert bytes to integer OK")
	} else {
		t.Error("Convert bytes to integer fail")
	}
}

func TestWriteFile(t *testing.T) {
	fileName := CheckFile("/tmp/temp/test.file")
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	dummyString := "Hello world!"
	WriteFile(fileName, []byte(dummyString))
	content, err := ReadFile(fileName)
	if err != nil {
		t.Error("Write file fail")
	}
	if content == dummyString {
		t.Logf("File content : \n%s\n", content)
		t.Log("Write file OK")
	}
}

func TestReadFile(t *testing.T) {
	fileName := CheckFile("/tmp/temp/test.file")
	content, _err_ := ReadFile(fileName)
	if _err_ != nil {
		t.Error("Read file fail")
	} else {
		t.Logf("File content : \n%s\n", content)
		t.Log("Read file OK")
	}
}