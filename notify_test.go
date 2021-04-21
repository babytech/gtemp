package main

import (
	"testing"
	"time"
)

func watchCallback() {
	rawData := make([]byte, 0)
	rawData = append(rawData, []byte("Hello world")...)

	WriteFile("/tmp/temp/persistent", rawData)
}

func TestTempSensorConfig_WatchFile(t *testing.T) {
	tempSensor := NewTempSensor()
	tempSensor.WatchFile("/tmp/temp/notify", watchCallback)
	WriteFile("/tmp/temp/notify", []byte("1"))
	timer := time.NewTimer(time.Second * time.Duration(1))
	go func() {
		for {
			select {
			case <-timer.C:
				content, _err_ := ReadFile("/tmp/temp/persistent")
				if _err_ != nil {
					t.Error("Read file fail")
				} else {
					t.Logf("File content : \n%s\n", content)
					t.Log("Read file OK")
				}
				timer.Reset(time.Second * time.Duration(1))
			}
		}
	}()
}
