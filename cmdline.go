package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	help bool
	version bool
	Version bool
	testJsonFile bool
	TestJsonFile bool
	readCsvFile bool
	writeCsvFile bool
	dummyTemp bool
	sendSignal string
	prefix string
	jsonFile string
	csvFile string
	notifyFile string
	chartMode string
	eepromSize uint
)

const VersionOfThisProgram = "0.0.7"
const AuthorInformation = "Babytech"

func showVersion() {
	fmt.Println("Version:", VersionOfThisProgram)
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `gtemp version: gtemp %s
Author: %s
Usage: gtemp [-hvVtTrwd] [-s signal] [-p prefix] [-j json_file] [-c csv_file] [-n notify_file][-m chart_mode <bar/line>] [-e size]
Options:
`, VersionOfThisProgram, AuthorInformation)
	flag.PrintDefaults()
}

func InitCmdLine() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&version, "v", false, "show version and exit")
	flag.BoolVar(&Version, "V", false, "show version and configure options then exit")
	flag.BoolVar(&testJsonFile, "t", false, "test JSON configuration and exit")
	flag.BoolVar(&TestJsonFile, "T", false, "test JSON configuration, dump it and exit")
	flag.BoolVar(&readCsvFile, "r", false, "read CSV file and exit")
	flag.BoolVar(&writeCsvFile, "w", false, "write CSV file and exit")
	flag.BoolVar(&dummyTemp, "d", false, "write dummy temperature into sensor file period")
	// Note default is: -s string，change to: -s signal here
	flag.StringVar(&sendSignal, "s", "", "send `signal` to a master process: stop, quit, reopen, reload")
	flag.StringVar(&prefix, "p", "/tmp/", "set `prefix` path")
	flag.StringVar(&jsonFile, "j", "temp/config.json", "set configuration 'input_file` -> json format")
	flag.StringVar(&csvFile, "c", "temp/data/gTemp.csv", "set configuration 'output_file` -> csv format")
	flag.StringVar(&notifyFile, "n", "temp/notify.txt", "set notify file to flush data cache to persistent storage <eeprom>")
	flag.StringVar(&chartMode, "m", "bar", "set mode for show chart from the output of csv file")
	flag.UintVar(&eepromSize, "e", 128*16, "set the raw data size 'n bytes' of persistent storage <eeprom>")
	// Override default usage function
	flag.Usage = usage
}

func ParseCmdLine() int {
	flag.Parse()
	if help {
		flag.Usage()
		return -1
	}
	if version {
		showVersion()
		return -2
	}
	if Version {
		showVersion()
		flag.PrintDefaults()
		return -3
	}
	return 0
}