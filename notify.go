package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
)

type Callback func()

func (g *TempSensorConfig) WatchFile(fileName string, funcCb Callback) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Watch file is empty from cmdline. Error: ", err)
		fmt.Println("Watch file is fetching from Json file: ", g.Notify.File)
		fileName = g.Notify.File
	}
	defer file.Close()
	watcher := CheckFile(fileName)
	f, err := os.Create(watcher)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("NewWatcher fail!")
		log.Fatal(err)
	}
	defer func(watch *fsnotify.Watcher) {
		_ = watch.Close()
	}(watch)
	err = watch.Add(watcher)
	if err != nil {
		fmt.Println("Watch.Add fail! fileName: ", fileName)
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case ev := <-watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("Create File : ", ev.Name)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("Write File : ", ev.Name)
						funcCb()
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("Remove File : ", ev.Name)
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						log.Println("Rename File : ", ev.Name)
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						log.Println("Modify File : ", ev.Name)
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("Error : ", err)
					return
				}
			}
		}
	}()
}