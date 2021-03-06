package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type ProfileConfig struct {
	Port string
}

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

type FuseConfig struct {
	Path  string
	Mount string
}

type RotorConfig struct {
	Name  string
	Speed string
}

type FanConfig struct {
	Name     string
	Path     string
	Presence string
	Number   string
	Rotors   []RotorConfig
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
	Profile     ProfileConfig
	Sampling    SamplingConfig
	Persistence PersistenceConfig
	Notify      NotifyConfig
	Csv         CsvConfig
	Fuse        FuseConfig
	Fans        []FanConfig
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
