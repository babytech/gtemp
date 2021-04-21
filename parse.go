package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type SamplingConfig struct {
	Low      int
	High     int
	Step     int
	Interval int
}

type PersistenceConfig struct {
	Interval int
	File     string
}

type NotifyConfig struct {
	File string
}

type CsvConfig struct {
	Interval int
	File     string
	Ip       string
	Port     string
}

type SensorConfig struct {
	Name               string
	File               string
	value              int
	dummyTemp          int
	dummyTempIncrement int
	Cache              [32]uint32
}

type TempSensorConfig struct {
	RawData     []byte
	Sampling    SamplingConfig
	Persistence PersistenceConfig
	Notify      NotifyConfig
	Csv         CsvConfig
	Sensors     []SensorConfig
}

func NewTempSensor() *TempSensorConfig {
	return &TempSensorConfig{
		RawData: make([]byte, 0),
	}
}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) load(filename string, g interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.New("read file fail")
	}
	err = json.Unmarshal(data, g)
	if err != nil {
		return errors.New("unmarshal fail")
	}
	return nil
}

func (g *TempSensorConfig) ParseJsonFile(filename string) error {
	JsonParse := NewJsonStruct()
	err := JsonParse.load(filename, g)
	if err != nil {
		return err
	}
	return nil
}
