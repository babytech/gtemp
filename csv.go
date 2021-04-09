package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func readCsv(fileName string) {
	fs, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("can not open the file, err is %+v", err)
	}
	defer fs.Close()
	r := csv.NewReader(fs)
	// For big file, read for each line
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			log.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		fmt.Println(row)
	}
	fmt.Println("----------------------------------")
	fs1, _ := os.Open(fileName)
	r1 := csv.NewReader(fs1)
	content, err := r1.ReadAll()
	if err != nil {
		log.Fatalf("can not readall, err is %+v", err)
	}
	for _, row := range content {
		fmt.Println(row)
	}
}

func writeCsv(file string) {
	fileName := CheckFile(file)
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// Write as UTF-8 BOM
	_, _ = f.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(f)
	data := [][]string{
		{"1", "LiuBei", "23"},
		{"2", "ZhangFei", "23"},
		{"3", "GunYu", "23"},
		{"4", "ZhaoYun", "23"},
		{"5", "HuangZhang", "23"},
		{"6", "MaChao", "23"},
	}
	_ = w.WriteAll(data)
	w.Flush()
}

func (g *TempSensorConfig) genCsv(file string) {
	fileName := CheckFile(file)
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// Write as UTF-8 BOM
	_, _ = f.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(f)
	var data [][]string
	var row []string
	title := []string{"name", "<-40", "-40~-35", "-35~-30", "-30~-25", "-25~-20", "-20~-15", "-15~-10", "-10~-5",
		"-5~0", "0~5", "5~10", "10~15", "15~20", "20~25", "25~30", "30~35", "35~40", "40~45", "45~50",
		"50~55", "55~60", "60~65", "65~70", "70~75", "75~80", "80~90", "90~95", "95~100", ">100"}
	data = append(data, title)
	for index := 0; index < len(g.Sensors); index++ {
		row = append(row, g.Sensors[index].Name)
		for k := 2; k < 32; k++ {
			row = append(row, strconv.Itoa(int(g.Sensors[index].Cache[k])))
		}
		data = append(data, row)
		row = make([]string,0)
	}
	//fmt.Println("[][]:", data)
	_ = w.WriteAll(data)
	w.Flush()
}

func exportCsv(filePath string, data [][]string) {
	fp, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Create File["+filePath+"]Failed,%v", err)
		return
	}
	defer fp.Close()
	// Write as UTF-8 BOM
	_, _ = fp.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(fp)
	_ = w.WriteAll(data)
	w.Flush()
}

func UnitTestWriteCsv(fileName string) {
	columns := [][]string{{"Name", "Phone", "Company", "Job", "Join-time"}, {"1", "2", "Baby,Baby,Baby", "4", "5"}}
	exportCsv(fileName,columns)
}

func UnitTestReadCsv(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		// record has the type []string
		fmt.Println(record)
	}
}

func (g *TempSensorConfig) HandleCsvFile(fileName string) int {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("FileName is empty from cmdline. Error: ", err)
		fmt.Println("FileName is fetching from Json file: ", g.Csv.File)
		fileName = g.Csv.File
	}
	defer file.Close()
	fmt.Println("Handle csv File :", fileName)
	// Handle csv file jobs
	if r {
		fmt.Println("ReadCsv :", fileName)
		readCsv(fileName)
		return 1
	} else if w {
		fmt.Println("WriteCsv :", fileName)
		writeCsv(fileName)
		return 2
	}
	go func() {
		for {
			//time.Sleep(time.Second * 30)
			time.Sleep(time.Second * time.Duration(g.Csv.Interval))
			g.genCsv(fileName)
		}
	}()
	return 0
}