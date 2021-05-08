package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const ByteLengthPerTemp = 128
const ByteLengthPerCount = 4

func (g *TempSensorConfig) updateCurrentTempValue(temp int, index int) {
	g.Sensors[index].value = temp
}

func (g *TempSensorConfig) updateDataCache(currentTemp int, index int) {
	if currentTemp < g.Sampling.Low {
		currentTemp = g.Sampling.Low - g.Sampling.Step
	} else if currentTemp == g.Sampling.High {
		currentTemp = g.Sampling.High - g.Sampling.Step
	} else if currentTemp > g.Sampling.High {
		currentTemp = g.Sampling.High
	}
	deltaTemp := currentTemp - g.Sampling.Low
	fmt.Printf("currentTemp: %d, g.Sampling.Low: %d, deltaTemp: %d\n", currentTemp, g.Sampling.Low, deltaTemp)
	offset := deltaTemp/g.Sampling.Step + 1 + 2
	g.Sensors[index].Cache[offset]++
	fmt.Printf("g.Sensors[%d].Cache[%d] : %d\n", index, offset, g.Sensors[index].Cache[offset])
}

func (g *TempSensorConfig) writeDataToPersistence() {
	fmt.Println("Flushing data cache into persistence storage...")
	//merge all sensors cache together
	g.RawData = make([]byte, 0)
	for index := 0; index < len(g.Sensors); index++ {
		out := []byte(g.Sensors[index].Name)
		g.RawData = append(g.RawData, out[0:8]...)
		for k := 2; k < 32; k++ {
			g.RawData = append(g.RawData, IntToBytes(g.Sensors[index].Cache[k])...)
		}
		out = make([]byte, 0)
	}
	fmt.Println("writeDataToPersistence()->len(g.RawData): ", len(g.RawData))
	fmt.Println("g.RawData:", g.RawData)
	WriteFile(g.Persistence.File, g.RawData)
}

func (g *TempSensorConfig) ReadDataFromPersistence() (string, error) {
	var persistenceFile string
	fmt.Println("History Data File:", dataFile)
	_, err := os.Open(dataFile)
	if err != nil {
		fmt.Printf("History Data File: %s is not exist\n", dataFile)
		persistenceFile = g.Persistence.File
	} else {
		fmt.Printf("History Data File: %s is exist\n", dataFile)
		persistenceFile = dataFile
	}
	content, err := ReadFile(persistenceFile)
	if err != nil {
		return persistenceFile, err
	}
	g.RawData = append(g.RawData, []byte(content)...)
	//fmt.Println("ReadDataFromPersistence() -> len(g.RawData) = ", len(g.RawData))
	if len(g.RawData) < ByteLengthPerTemp*len(g.Sensors) {
		reserveRawLength := ByteLengthPerTemp*len(g.Sensors) - len(g.RawData)
		reserveRawData := make([]byte, reserveRawLength)
		g.RawData = append(g.RawData, reserveRawData...)
		fmt.Println("len(g.RawData) < ByteLengthPerTemp*len(g.Sensors) -> len(g.RawData) = ", len(g.RawData))
	}
	//fmt.Println("ReadDataFromPersistence() -> ByteLengthPerTemp*len(g.Sensors) = ", ByteLengthPerTemp*len(g.Sensors))
	for index := 0; index < len(g.Sensors); index++ {
		// get each 128 bytes from v.RawData
		lowRange := index * ByteLengthPerTemp
		highRange := (index+1)*ByteLengthPerTemp - 1
		rawBuffer := g.RawData[lowRange:highRange]
		// covert 128 bytes to 32 int and write to cache of each sensor
		for k := 0; k < 32; k++ {
			lowPerRange := k * ByteLengthPerCount
			highPerRange := (k + 1) * ByteLengthPerCount
			g.Sensors[index].Cache[k] = BytesToInt(rawBuffer[lowPerRange:highPerRange])
		}
	}
	return persistenceFile, nil
}

func (g *TempSensorConfig) MonitorBody(index int) {
	//start go routine
	go func() {
		fmt.Printf("\nStart Read TempSensor: %s\nTempSensor path: %s \n",
			g.Sensors[index].Name, g.Sensors[index].File)
		for {
			//read sensors value from sysfs
			s, err := ReadFile(g.Sensors[index].File)
			if err != nil {
				fmt.Printf("TempSensor[%d] readFile: %s FAIL!\n", index, g.Sensors[index].File)
			} else {
				fmt.Printf("TempSensor[%d] readFile: %s OK!\n", index, g.Sensors[index].File)
				fmt.Println("FileContent: ", s)
			}
			str := strings.Replace(s, "\n", "", -1)
			curTemp, _ := strconv.Atoi(str)
			curTemp = curTemp / 1000
			fmt.Println("curTemp = ", curTemp)
			g.updateCurrentTempValue(curTemp, index)
			g.updateDataCache(curTemp, index)
			time.Sleep(time.Second * time.Duration(g.Sampling.Interval))
		}
	}()
}

func (g *TempSensorConfig) StartTimer() {
	t := time.NewTimer(time.Second * time.Duration(g.Persistence.Interval))
	go func() {
		for {
			select {
			case <-t.C:
				fmt.Printf("\nTimer (%ds) occur!\n", g.Persistence.Interval)
				g.writeDataToPersistence()
				t.Reset(time.Second * time.Duration(g.Persistence.Interval))
			}
		}
	}()
}

func (g *TempSensorConfig) SignalListen() {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGUSR2)
		for {
			s := <-c
			// Receive signal and do custom job
			fmt.Println("get signal:", s)
			g.writeDataToPersistence()
		}
	}()
}

func (g *TempSensorConfig) doCreateDummyTemp(index int) {
	for {
		for item := 0; item < (g.Sampling.High-g.Sampling.Low)*1000/g.Sensors[index].dummyTempIncrement; item++ {
			WriteFile(g.Sensors[index].File, []byte(strconv.Itoa(g.Sensors[index].dummyTemp)))
			g.Sensors[index].dummyTemp = g.Sensors[index].dummyTemp + g.Sensors[index].dummyTempIncrement
			time.Sleep(time.Second * time.Duration(3))
		}
	}
}

func (g *TempSensorConfig) createDummyTemp() {
	g.Sensors[0].dummyTemp = -50000
	g.Sensors[0].dummyTempIncrement = 1000
	g.Sensors[1].dummyTemp = -5000
	g.Sensors[1].dummyTempIncrement = 2000
	go g.doCreateDummyTemp(0)
	go g.doCreateDummyTemp(1)
	fmt.Println("create dummy temperature for sensors...")
}

func (g *TempSensorConfig) FuseMain(path string) {
	c, err := fuse.Mount(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func(c *fuse.Conn) {
		_ = c.Close()
	}(c)
	srv := fs.New(c, nil)
	fileSys := &FS{
		&Dir{Node: Node{name: "head", inode: NewInode()}, directories: &[]*Dir{}}}
	dirs := make([]*Dir, 0)
	for index := 0; index < len(g.Sensors); index++ {
		dirNode := &Dir{Node: Node{name: g.Sensors[index].Name, inode: NewInode()}, files: &[]*File{
			{Node: Node{name: "temp", inode: NewInode()}, item: g.Sensors[index].File},
		}}
		dirs = append(dirs, dirNode)
	}
	fileSys.root.directories = &dirs
	log.Println("About to serve fs")
	if err := srv.Serve(fileSys); err != nil {
		log.Panicln(err)
	}
}

func (g *TempSensorConfig) StartFuse() error {
	var mountPoint string
	exist, err := PathExists(mountDir)
	if err != nil {
		fmt.Printf("Error: %s Call PathExists fail!\n", mountDir)
		return err
	} else if exist == false {
		fmt.Printf("FUSE mount directory: %s is not exist from cmdline!\n", mountDir)
		fmt.Println("FUSE mount directory is fetching from Json file: ", g.Fuse.Mount)
		mountPoint = g.Fuse.Mount
	} else {
		fmt.Printf("FUSE mount directory: %s is exist from cmdline.\n", mountDir)
		fmt.Println("FUSE mount directory is fetching from cmdline: ", mountDir)
		mountPoint = mountDir
	}
	go func() {
		fmt.Println("FUSE Mount Directory:", mountPoint)
		err := os.MkdirAll(mountPoint, 0711)
		if err != nil {
			log.Println("Error creating directory", mountPoint)
			log.Println(err)
		}
		g.FuseMain(mountPoint)
	}()
	return nil
}

func ParseConfigurationFile(fileName string) *TempSensorConfig {
	tempSensor := NewTempSensor()
	err := tempSensor.ParseJsonFile(fileName)
	if err != nil {
		fmt.Printf("Test Configure File: %s Fail!\n", fileName)
		fmt.Printf("===> gtemp [%s]: exit...\n", VersionInformation)
		return nil
	} else if testJsonFile {
		fmt.Printf("Test Configure File: %s OK!\n", fileName)
	} else if TestJsonFile {
		fmt.Printf("Test Configure JSON File: %s OK!\n", fileName)
		fmt.Println("Parse Configure JSON File Result:")
		fmt.Println(tempSensor.Sampling)
		fmt.Println(tempSensor.Persistence)
		fmt.Println(tempSensor.Csv)
		fmt.Println(tempSensor.Fuse)
		for index := 0; index < len(tempSensor.Sensors); index++ {
			fmt.Println(tempSensor.Sensors[index].Name)
			fmt.Println(tempSensor.Sensors[index].File)
		}
		os.Exit(0)
	}
	return tempSensor
}

func StartMonitorTask() {
	var err error
	var result int
	var configFileForTempMonitor string
	var csvFileForTempMonitor string
	var persistenceFile string
	_, _ = fmt.Fprintf(os.Stderr, welcomeInformation, AuthorInformation, VersionInformation)
	fmt.Printf("===> gtemp [%s]: Start Monitoring Task for capture temperature...\n", VersionInformation)
	fmt.Printf("size of eeprom: 0x%x\n", eepromSize)
	// Parse Configuration File
	configFileForTempMonitor = jsonFile
	fmt.Println("Configuration File:", configFileForTempMonitor)
	tempSensor := ParseConfigurationFile(configFileForTempMonitor)
	if tempSensor == nil {
		fmt.Println("Parse Configuration File Fail!")
		goto failExit
	}
	csvFileForTempMonitor = prefix + csvFile
	fmt.Println("CSV File:", csvFileForTempMonitor)
	// Initial read data from eeprom file
	persistenceFile, err = tempSensor.ReadDataFromPersistence()
	fmt.Println("Persistence File:", persistenceFile)
	if err != nil {
		fmt.Println("ReadDataFromPersistence: Fail!")
		goto failExit
	} else if persistenceFile == dataFile {
		fmt.Println("Generate CSV file: ", csvFileForTempMonitor)
		tempSensor.doGenCsvStatistics(csvFileForTempMonitor)
		fmt.Printf("Upload CSV file to Server -> %v:%d\n", tempSensor.Csv.Ip, &tempSensor.Csv.Port)
		tempSensor.UploadCsvFile(csvFileForTempMonitor)
		goto successExit
	}
	// Startup FUSE function
	err = tempSensor.StartFuse()
	if err != nil {
		fmt.Println("StartFuse: Fail!")
		goto failExit
	}
	// Generate dummy temperature and write to per sensor file
	if dummyTemp {
		tempSensor.createDummyTemp()
	}
	// Start Monitoring for each temperature sensor
	for index := 0; index < len(tempSensor.Sensors); index++ {
		tempSensor.MonitorBody(index)
		fmt.Printf("Ready to call MonitorBody(%d)...\n", index)
	}
	// Create timer to do flush data cache
	tempSensor.StartTimer()
	// Handle signal for write data to eeprom file
	tempSensor.SignalListen()
	// Handle generate CSV file
	result = tempSensor.HandleCsvFile(csvFileForTempMonitor)
	if result != 0 {
		goto failExit
	}
	// Before board reset using script to set /tmp/temp/notify.txt as 1,
	// Register this notify to trigger write data to eeprom file
	tempSensor.WatchFile(notifyFile, tempSensor.writeDataToPersistence)
	return
failExit:
	fmt.Printf("===> gtemp [%s]: NOK exit...\n", VersionInformation)
	os.Exit(1)
successExit:
	fmt.Printf("===> gtemp [%s]: OK exit...\n", VersionInformation)
	os.Exit(0)
}
