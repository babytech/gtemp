package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	help         bool
	version      bool
	Version      bool
	testJsonFile bool
	TestJsonFile bool
	readCsvFile  bool
	writeCsvFile bool
	dummyTemp    bool
	sendSignal   string
	mountDir     string
	prefix       string
	jsonFile     string
	csvFile      string
	dataFile     string
	notifyFile   string
	eepromSize   uint
)

const VersionInformation = "0.1.2"
const AuthorInformation = "Babytech"
const welcomeInformation = `
         __
   _____/  |_  ____   _____ ______      Temperature Monitoring
  / ___\   __\/ __ \ /     \\____ \      --- written by Golang
 / /_/  >  | \  ___/|  Y Y  \  |_> >
 \___  /|__|  \___  >__|_|  /   __/     Author:  %s
/_____/           \/      \/|__|        Version: %s

`

func showVersion() {
	fmt.Println("gtemp version: ", VersionInformation)
}

func usage() {
	var programName = filepath.Base(os.Args[0])
	_, _ = fmt.Fprintf(os.Stderr, `%s version: %s
Author: %s
Usage: gtemp [-hvVtTrwf] [-s signal] [-j json_file] [-m mount_point] [-p prefix] [-c csv_file] [-d data_file] [-n notify_file] [-e size]
Options:
`, programName, VersionInformation, AuthorInformation)
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
	flag.BoolVar(&dummyTemp, "f", false, "write dummy temperature into sensor file period")
	// Note default is: -s string，change to: -s signal here
	flag.StringVar(&sendSignal, "s", "", "send `signal` to a master process: stop, quit, reopen, reload")
	flag.StringVar(&jsonFile, "j", "config.json", "set configuration 'input_file' ->[.json])")
	flag.StringVar(&mountDir, "m", "/tmp/temp/fuse", "set mount point path for FUSE")
	flag.StringVar(&prefix, "p", "./temp/data/", "set `prefix` path")
	flag.StringVar(&csvFile, "c", "MF14/temp.csv", "set statistics <product>/<csv_file> as 'input_file' ->[.csv]")
	flag.StringVar(&dataFile, "d", "./temp/data.bin", "set data 'input_file' ->[.bin] to generate statistics of history temperature")
	flag.StringVar(&notifyFile, "n", "./temp/notify.txt", "set notify 'input_file' ->[.txt] to flush data cache to persistent storage <eeprom>")
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
