package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	h bool
	v bool
	V bool
	t bool
	T bool
	r bool
	w bool
	d bool
	s string
	p string
	j string
	c string
	n string
	e uint
)

const VersionOfThisProgram = "0.0.5"
const AuthorInformation = "Babytech"

func showVersion() {
	fmt.Println("Version:", VersionOfThisProgram)
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `gtemp version: gtemp %s
Author: %s
Usage: gtemp [-hvVtTrwd] [-s signal] [-p prefix] [-j json_file] [-c csv_file] [-n notify_file] [-e size] 
Options:
`, VersionOfThisProgram, AuthorInformation)
	flag.PrintDefaults()
}

func InitCmdLine() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.BoolVar(&V, "V", false, "show version and configure options then exit")
	flag.BoolVar(&t, "t", false, "test JSON configuration and exit")
	flag.BoolVar(&T, "T", false, "test JSON configuration, dump it and exit")
	flag.BoolVar(&r, "r", false, "read CSV file and exit")
	flag.BoolVar(&w, "w", false, "write CSV file and exit")
	flag.BoolVar(&d, "d", false, "write dummy temperature into sensor file period")
	// Note default is: -s string，change to: -s signal here
	flag.StringVar(&s, "s", "", "send `signal` to a master process: stop, quit, reopen, reload")
	flag.StringVar(&p, "p", "/tmp/", "set `prefix` path")
	flag.StringVar(&j, "j", "temp/config.json", "set configuration 'input_file` -> json format")
	flag.StringVar(&c, "c", "temp/data/gTemp.csv", "set configuration 'output_file` -> csv format")
	flag.StringVar(&n, "n", "temp/notify.txt", "set notify file to flush data cache to persistent storage <eeprom>")
	flag.UintVar(&e, "e", 128*16, "set the raw data size 'n bytes' of persistent storage <eeprom>")
	// Override default usage function
	flag.Usage = usage
}

func ParseCmdLine() int {
	flag.Parse()
	if h {
		flag.Usage()
		return -1
	}
	if v {
		showVersion()
		return -2
	}
	if V {
		showVersion()
		flag.PrintDefaults()
		return -3
	}
	return 0
}