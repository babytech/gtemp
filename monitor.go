package main

import (
	"fmt"
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
	g.RawData = make ([]byte, 0)
	for index := 0; index < len(g.Sensors); index++ {
		out := []byte(g.Sensors[index].Name)
		g.RawData = append(g.RawData, out[0:8]...)
		for k :=2; k < 32; k++ {
			g.RawData = append(g.RawData, IntToBytes(g.Sensors[index].Cache[k])...)
		}
		out = make([]byte, 0)
	}
	fmt.Println("len(g.RawData): ",len(g.RawData))
	fmt.Println("g.RawData:",g.RawData)
	WriteFile(g.Persistence.File, g.RawData)
}

func (g *TempSensorConfig) ReadDataFromPersistence() error {
	content, err := ReadFile(g.Persistence.File)
	if err != nil {
		return err
	}
	g.RawData = append(g.RawData, []byte(content)...)
	for index := 0; index < len(g.Sensors); index++ {
		// get each 128 bytes from v.RawData
		lowRange := index * ByteLengthPerTemp
		highRange := (index+1) * ByteLengthPerTemp - 1
		rawBuffer := g.RawData[lowRange:highRange]
		// covert 128 bytes to 32 int and write to cache of each sensor
		for k :=0; k < 32; k++ {
			lowPerRange := k * ByteLengthPerCount
			highPerRange := (k + 1) * ByteLengthPerCount
			g.Sensors[index].Cache[k] = BytesToInt(rawBuffer[lowPerRange:highPerRange])
		}
	}
	return nil
}

func (g *TempSensorConfig) MonitorBody(index int) {
	//start go routine
	go func() {
		fmt.Printf("Start Read TempSensor: %s\nTempSensor path: %s \n",
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
			curTemp = curTemp/1000
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
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR2)
	for {
		s := <-c
		// Receive signal and do custom job
		fmt.Println("get signal:", s)
		g.writeDataToPersistence()
	}
}

func (g *TempSensorConfig) doCreateDummyTemp(index int) {
	for {
		for item := 0; item < (g.Sampling.High-g.Sampling.Low) * 1000/g.Sensors[index].dummyTempIncrement; item++ {
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

func StartMonitorTask() {
	fmt.Printf("===> gtemp [%s]: Start Monitoring Task for capture temperature...\n", VersionOfThisProgram)
	ConfigFileForTempMonitor := prefix + jsonFile
	fmt.Println("Configuration File :", ConfigFileForTempMonitor)
	fmt.Printf("size of eeprom : 0x%x\n", eepromSize)
	tempSensor := NewTempSensor()
	err := tempSensor.ParseJsonFile(ConfigFileForTempMonitor)
	if err != nil {
		fmt.Printf("Test Configure File: %s Fail!\n", ConfigFileForTempMonitor)
		os.Exit(1)
	} else if testJsonFile {
		fmt.Printf("Test Configure File: %s OK!\n", ConfigFileForTempMonitor)
		os.Exit(0)
	} else if TestJsonFile {
		fmt.Printf("Test Configure File: %s OK!\n", ConfigFileForTempMonitor)
		fmt.Println("Parse Configure File Data Result :")
		fmt.Println(tempSensor.Sampling)
		fmt.Println(tempSensor.Persistence)
		fmt.Println(tempSensor.Csv)
		for index := 0; index < len(tempSensor.Sensors); index++ {
			fmt.Println(tempSensor.Sensors[index].Name)
			fmt.Println(tempSensor.Sensors[index].File)
		}
		os.Exit(0)
	}

	// Initial read data from eeprom file
	err = tempSensor.ReadDataFromPersistence()
	if err != nil {
		fmt.Println("ReadDataFromPersistence: Fail!")
		os.Exit(1)
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
	go tempSensor.SignalListen()

	// Handle generate CSV file
	result := tempSensor.HandleCsvFile(prefix + csvFile)
	if result != 0 {
		os.Exit(0)
	}

	// Before board reset using script to set /tmp/temp/notify.txt as 1,
	// Register this notify to trigger write data to eeprom file
	tempSensor.WatchFile(prefix + notifyFile, tempSensor.writeDataToPersistence)
}