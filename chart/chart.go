package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const versionInformation = "0.1.1"
const authorInformation = "Babytech"
const chartInformation = `
============================================================================
|           __               __     Temperature Monitoring Web Tools       |
|     _____/ /_  ____ ______/ /_               --- written by Golang       |
|    / ___/ __ \/ __ \/ ___/ __/    Author: %s                       |
|   / /__/ / / / /_/ / /  / /_      Version: %s                         |
|   \___/_/ /_/\__,_/_/   \__/      Http & File Server @IP: %s  |
============================================================================
`

const ipPort = "8088"

func externalIP() (net.IP, error) {
	iFaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iFace := range iFaces {
		if iFace.Flags&net.FlagUp == 0 {
			continue
		}
		if iFace.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrS, err := iFace.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrS {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		// not an ipv4 address
		return nil
	}
	return ip
}

func drawLiquidPin(w http.ResponseWriter, r *csv.Reader) []*charts.Liquid {
	var liquids []*charts.Liquid
	var liquidTitle string
	var data []float32
	var ratio float32
	row := make([]string, 0)
	rowData := make([]string, 0)
	item := make([]opts.LiquidData, 0)
	const lowValue = -30
	const highValue = 80
	name, err := r.Read()
	if err != nil && err != io.EOF {
		http.Error(w, "File read failed.", 404)
		fmt.Println("r.Read() Error:", err)
	}
	if err == io.EOF {
		fmt.Println("r.Read() Error:", err)
	}
	for {
		rowData, err = r.Read()
		if err != nil && err != io.EOF {
			http.Error(w, "File read failed.", 404)
			fmt.Println("r.Read() Error:", err)
		}
		if err == io.EOF {
			break
		}
		row = rowData
	}
	//fmt.Println("row: ", row)
	for i := 1; i < len(row); i++ {
		liquid := charts.NewLiquid()
		item = make([]opts.LiquidData, 0)
		value, _ := strconv.Atoi(row[i])
		if value > 0 {
			ratio = float32(value) / highValue
		} else {
			ratio = float32(math.Abs(float64(value) / lowValue))
		}
		data = []float32{ratio, ratio * 0.8, ratio * 1.2}
		for j := 0; j < len(data); j++ {
			item = append(item, opts.LiquidData{Value: data[j]})
		}
		liquidTitle = "温度传感器" + string(name[i]) + "温度水位"
		liquid.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{
				Title: liquidTitle,
			}),
		)
		liquid.AddSeries("liquid", item).
			SetSeriesOptions(
				charts.WithLiquidChartOpts(opts.LiquidChart{
					IsWaveAnimation: true,
					Shape:           "pin",
				}),
			)
		liquids = append(liquids, liquid)
	}
	return liquids
}

func drawEffectScatter(w http.ResponseWriter, r *csv.Reader) *charts.EffectScatter {
	var timeStamp string
	es := charts.NewEffectScatter()
	row := make([]string, 0)
	temp := make([]opts.EffectScatterData, 0)
	name, err := r.Read()
	if err != nil && err != io.EOF {
		http.Error(w, "File read failed.", 404)
		fmt.Println("r.Read() Error:", err)
	}
	if err == io.EOF {
		fmt.Println("r.Read() Error:", err)
	}
	for {
		row, err = r.Read()
		if err != nil && err != io.EOF {
			http.Error(w, "File read failed.", 404)
			fmt.Println("r.Read() Error:", err)
		}
		if err == io.EOF {
			break
		}
		timeStamp = row[0]
		temp = make([]opts.EffectScatterData, 0)
		for i := 1; i < len(row); i++ {
			temp = append(temp, opts.EffectScatterData{Value: row[i]})
		}
	}
	es.SetGlobalOptions(
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "600px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "温度传感器实时温度值" + "    当前时间:" + timeStamp,
			Top:   "top, omitempty",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "温度值",
			SplitLine: &opts.SplitLine{
				Show: true,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "温度传感器",
		}),
	)
	es.SetXAxis(name[1:]).AddSeries("temp", temp,
		charts.WithRippleEffectOpts(opts.RippleEffect{
			Period:    4,
			Scale:     10,
			BrushType: "stroke",
		})).
		SetSeriesOptions(charts.WithLabelOpts(
			opts.Label{
				Show:     true,
				Position: "top",
			}),
		)
	return es
}

func drawFanRotor(w http.ResponseWriter, r *csv.Reader) []*charts.EffectScatter {
	var fans []*charts.EffectScatter
	var fan string
	rotors := make([]opts.EffectScatterData, 0)
	row := make([]string, 0)
	name, err := r.Read()
	if err != nil && err != io.EOF {
		http.Error(w, "File read failed.", 404)
		fmt.Println("r.Read() Error:", err)
	}
	if err == io.EOF {
		fmt.Println("r.Read() Error:", err)
	}
	for {
		row, err = r.Read()
		if err != nil && err != io.EOF {
			http.Error(w, "File read failed.", 404)
			fmt.Println("r.Read() Error:", err)
		}
		if err == io.EOF {
			break
		}
		fan = row[0]
		es := charts.NewEffectScatter()
		rotors = make([]opts.EffectScatterData, 0)
		for i := 1; i < len(row); i++ {
			rotors = append(rotors, opts.EffectScatterData{Value: row[i]})
		}
		es.SetGlobalOptions(
			charts.WithTooltipOpts(opts.Tooltip{Show: true}),
			charts.WithInitializationOpts(opts.Initialization{
				Width:  "1200px",
				Height: "600px",
			}),
			charts.WithTitleOpts(opts.Title{
				Title: "风扇" + fan + "马达实时转速",
			}),
			charts.WithYAxisOpts(opts.YAxis{
				Name: "实时转速",
				SplitLine: &opts.SplitLine{
					Show: true,
				},
			}),
			charts.WithXAxisOpts(opts.XAxis{
				Name: "风扇马达",
			}),
		)
		es.SetXAxis(name[1:]).AddSeries("fan", rotors,
			charts.WithRippleEffectOpts(opts.RippleEffect{
				Period:    4,
				Scale:     10,
				BrushType: "fill",
			})).
			SetSeriesOptions(charts.WithLabelOpts(
				opts.Label{
					Show:     true,
					Position: "top",
				}),
			)
		fans = append(fans, es)
	}
	return fans
}

var (
	tempRange = []string{"<-40℃", "-40~-35℃", "-35~-30℃", "-30~-25℃", "-25~-20℃", "-20~-15℃", "-15~-10℃", "-10~-5℃",
		"-5~0℃", "0~5℃", "5~10℃", "10~15℃", "15~20℃", "20~25℃", "25~30℃", "30~35℃", "35~40℃", "40~45℃", "45~50℃",
		"50~55℃", "55~60℃", "60~65℃", "65~70℃", "70~75℃", "75~80℃", "80~85℃", "85~90℃", "90~95℃", "95~100℃", ">100℃"}
)

func drawPieRoseRadius(w http.ResponseWriter, r *csv.Reader) []*charts.Pie {
	var pies []*charts.Pie
	rowLine := 0
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			http.Error(w, "File read failed.", 404)
			fmt.Println("r.Read() Error:", err)
		}
		if err == io.EOF {
			break
		}
		fmt.Println("row: ", row)
		if rowLine != 0 {
			pie := charts.NewPie()
			pie.SetGlobalOptions(
				charts.WithTitleOpts(opts.Title{
					Title: "温度传感器" + row[0] + "温度区间统计图",
				}),
			)
			items := make([]opts.PieData, 0)
			for i, n := range row {
				if i != 0 {
					v, err := strconv.Atoi(n)
					if err != nil {
						fmt.Println(err)
					}
					if v != 0 {
						items = append(items, opts.PieData{Name: tempRange[i-1], Value: v})
						fmt.Println("name:", tempRange[i], "value:", v)
					}
				}
			}
			pie.AddSeries("pie", items).
				SetSeriesOptions(
					charts.WithLabelOpts(opts.Label{
						Show:      true,
						Formatter: "{b}: {c}",
					}),
					charts.WithPieChartOpts(opts.PieChart{
						Radius:   []string{"30%", "75%"},
						RoseType: "radius",
					}),
				)
			pies = append(pies, pie)
		}
		rowLine++
	}
	return pies
}

func drawFanSpeed(w http.ResponseWriter, r *csv.Reader) []*charts.Gauge {
	var gauges []*charts.Gauge
	var gauge *charts.Gauge
	rotorName := make([]string, 0)
	rowLine := 0
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			http.Error(w, "File read failed.", 404)
			fmt.Println("r.Read() Error:", err)
		}
		if err == io.EOF {
			break
		}
		fmt.Println("row: ", row)
		if rowLine == 0 {
			rotorName = append(rotorName, row[1:]...)
		} else {
			for i, n := range row {
				gauge = charts.NewGauge()
				gauge.SetGlobalOptions(
					charts.WithTitleOpts(opts.Title{
						Title: "风扇#" + row[0],
					}),
				)
				items := make([]opts.GaugeData, 0)
				if i != 0 {
					v, err := strconv.Atoi(n)
					if err != nil {
						fmt.Println(err)
					}
					if v != 0 {
						items = append(items, opts.GaugeData{Name: rotorName[i-1] + "马达转速", Value: v})
						fmt.Println("name:", rotorName[i-1], "value:", v)
						gauge.AddSeries("rotor", items)
						gauges = append(gauges, gauge)
					}
				}
			}
		}
		rowLine++
	}
	return gauges
}

func drawLineDaily(w http.ResponseWriter, r *csv.Reader, fileName string) {
	timeStamp := strings.TrimPrefix(fileName, "daily-")
	titleName := "日期:" + timeStamp + "    温度值统计"
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true,
			Right: "10%",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  types.ThemeWesteros,
			Width:  "2400px",
			Height: "600px",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithTitleOpts(opts.Title{
			Title: titleName,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Degree Celsius",
			SplitLine: &opts.SplitLine{
				Show: false,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Time Slot",
		}),
	)
	date := make([]string, 0)
	temp := make(map[string][]opts.LineData)
	name, err := r.Read()
	if err != nil && err != io.EOF {
		http.Error(w, "File read failed.", 404)
		fmt.Println("r.Read() Error:", err)
	}
	if err == io.EOF {
		fmt.Println("r.Read() Error:", err)
	}
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			http.Error(w, "File read failed.", 404)
			fmt.Println("r.Read() Error:", err)
		}
		if err == io.EOF {
			break
		}
		fmt.Println("row: ", row)
		date = append(date, row[0])
		for i := 1; i < len(row); i++ {
			temp[name[i]] = append(temp[name[i]], opts.LineData{Value: row[i]})
		}
	}
	line.SetXAxis(date).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				Smooth: true,
			}),
		)
	for k, v := range temp {
		if k != "time" {
			line.AddSeries(k, v).
				SetSeriesOptions(
					charts.WithLabelOpts(
						opts.Label{Show: true},
					),
				)
		}
	}
	_ = line.Render(w)
}

func drawBar(w http.ResponseWriter, r *csv.Reader) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true,
			Right: "10%",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  types.ThemeWesteros,
			Width:  "1600px",
			Height: "500px",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "温度传感器温度区间发生次数统计图",
			Left:  "10%",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "发生次数",
			SplitLine: &opts.SplitLine{
				Show: false,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "温度区间",
		}))
	rowLine := 0
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			http.Error(w, "File read failed.", 404)
			fmt.Println("r.Read() Error:", err)
		}
		if err == io.EOF {
			break
		}
		fmt.Println("row: ", row)
		if rowLine != 0 {
			items := make([]opts.BarData, 0)
			for i := 1; i <= 30; i++ {
				items = append(items, opts.BarData{Value: row[i]})
			}
			bar.AddSeries(row[0], items).SetSeriesOptions(
				charts.WithMarkPointNameTypeItemOpts(
					opts.MarkPointNameTypeItem{Name: "Maximum", Type: "max"},
				),
				charts.WithMarkPointStyleOpts(
					opts.MarkPointStyle{Label: &opts.Label{Show: true}},
				),
			)
		} else {
			bar.SetXAxis(tempRange)
		}
		rowLine++
	}
	_ = bar.Render(w)
}

func drawChart(w http.ResponseWriter, r *http.Request) (*csv.Reader, *os.File, string) {
	params := mux.Vars(r)
	tag := params["tag"]
	session := params["session"]
	var inputFile string
	fmt.Printf("#####tag:%s, session:%s\n", tag, session)
	fmt.Println("======dir:", "data/"+tag+"/"+session)
	if session == "" {
		dir := "./temp/data/" + tag + "/daily/"
		fileName := FileModTime(dir)
		if fileName == nil {
			http.Error(w, "File not found.", 404)
			return nil, nil, ""
		}
		inputFile = fmt.Sprintf("./temp/data/%v/daily/%v", tag, fileName[0])
	} else if strings.Contains(session, "daily") {
		inputFile = fmt.Sprintf("./temp/data/%v/daily/%v.csv", tag, session)
	} else {
		inputFile = fmt.Sprintf("./temp/data/%v/%v.csv", tag, session)
	}
	fmt.Println("inputFile:", inputFile)
	fs, err := os.Open(inputFile)
	if err != nil {
		http.Error(w, "File not found.", 404)
		return nil, nil, ""
	}
	rr := csv.NewReader(fs)
	rr.FieldsPerRecord = -1
	return rr, fs, session
}

func statisticsChart(w http.ResponseWriter, r *http.Request) {
	rr, fs, _ := drawChart(w, r)
	if rr == nil || fs == nil {
		fmt.Println("Error: drawChart fail!")
		return
	}
	drawBar(w, rr)
	defer fs.Close()
}

func dailyChart(w http.ResponseWriter, r *http.Request) {
	rr, fs, fileName := drawChart(w, r)
	if rr == nil || fs == nil {
		fmt.Println("Error: drawChart fail!")
		return
	}
	drawLineDaily(w, rr, fileName)
	defer fs.Close()
}

func statusChart(w http.ResponseWriter, r *http.Request) {
	rr, fs, _ := drawChart(w, r)
	if rr == nil || fs == nil {
		fmt.Println("Error: drawChart fail!")
		return
	}
	page := components.NewPage()
	for _, n := range drawPieRoseRadius(w, rr) {
		page.AddCharts(n)
	}
	page.SetLayout(components.PageFlexLayout)
	err := page.Render(io.MultiWriter(w))
	if err != nil {
		return
	}
	defer fs.Close()
}

func fanSpeedChart(w http.ResponseWriter, r *http.Request) {
	rr, fs, _ := drawChart(w, r)
	if rr == nil || fs == nil {
		fmt.Println("Error: drawChart fail!")
		return
	}
	page := components.NewPage()
	for _, n := range drawFanSpeed(w, rr) {
		page.AddCharts(n)
	}
	page.SetLayout(components.PageFlexLayout)
	url := "http://" + getServerIP() + ":" + ipPort + "/"
	fmt.Println("URL: ", url)
	page.AssetsHost = url
	err := page.Render(io.MultiWriter(w))
	if err != nil {
		return
	}
	defer fs.Close()
}

func fanRotorChart(w http.ResponseWriter, r *http.Request) {
	rr, fs, _ := drawChart(w, r)
	if rr == nil || fs == nil {
		fmt.Println("Error: drawChart fail!")
		return
	}
	page := components.NewPage()
	for _, n := range drawFanRotor(w, rr) {
		page.AddCharts(n)
	}
	page.SetLayout(components.PageFlexLayout)
	url := "http://" + getServerIP() + ":" + ipPort + "/"
	fmt.Println("URL: ", url)
	page.AssetsHost = url
	err := page.Render(io.MultiWriter(w))
	if err != nil {
		return
	}
	defer fs.Close()
}

func watermarkChart(w http.ResponseWriter, r *http.Request) {
	rr, fs, _ := drawChart(w, r)
	if rr == nil || fs == nil {
		fmt.Println("Error: drawChart fail!")
		return
	}
	page := components.NewPage()
	for _, n := range drawLiquidPin(w, rr) {
		page.AddCharts(n)
	}
	page.SetLayout(components.PageFlexLayout)
	err := page.Render(io.MultiWriter(w))
	if err != nil {
		return
	}
	defer fs.Close()
}

func currentChart(w http.ResponseWriter, r *http.Request) {
	rr, fs, _ := drawChart(w, r)
	if rr == nil || fs == nil {
		fmt.Println("Error: drawChart fail!")
		return
	}
	page := components.NewPage()
	page.AddCharts(drawEffectScatter(w, rr))
	page.SetLayout(components.PageFlexLayout)
	err := page.Render(io.MultiWriter(w))
	if err != nil {
		return
	}
	defer fs.Close()
}

func CreateFile(fileName string) *os.File {
	fileName = "./" + fileName
	dir := filepath.Dir(fileName)
	if _, err := os.Stat(dir); err == nil {
		//fmt.Println("\nDirectory path exists:", dir)
	} else {
		fmt.Println("\nDirectory path not exists: ", dir)
		err := os.MkdirAll(dir, 0711)
		if err != nil {
			log.Println("Error creating directory")
			log.Println(err)
			return nil
		}
	}
	dst, _ := os.Create(fileName)
	return dst
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for {
		part, err := reader.NextPart()
		if part == nil {
			return
		}
		if err == io.EOF {
			break
		}
		if part.FileName() == "" {
			// this is FormData
			data, _ := ioutil.ReadAll(part)
			fmt.Printf("FormData=[%s]\n", string(data))
		} else {
			// This is FileData
			dst := CreateFile(part.FileName())
			_, err__ := io.Copy(dst, part)
			if err__ != nil {
				return
			}
			dst.Close()
		}
	}
}

func FileModTime(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		_, err2 := fmt.Fprintln(os.Stderr, err)
		if err2 != nil {
			return nil
		}
		os.Exit(1)
	}
	var modTime time.Time
	var name []string
	for _, fi := range files {
		if fi.Mode().IsRegular() {
			if !fi.ModTime().Before(modTime) {
				if fi.ModTime().After(modTime) {
					modTime = fi.ModTime()
					name = name[:0]
				}
				name = append(name, fi.Name())
			}
		}
	}
	fmt.Println("===> dir: ", dir)
	fmt.Println("===> names: ", name)
	if len(name) > 0 {
		fmt.Println(modTime, name)
	}
	return name
}

func chartFileServer() {
	p, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	http.Handle("/", http.FileServer(http.Dir(p)))
	err := http.ListenAndServe(":"+ipPort, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func getServerIP() string {
	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}
	return ip.String()
}

func chartHttpServer() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/", statisticsChart)
	router.HandleFunc("/{tag}/{session}/statistics", statisticsChart)
	router.HandleFunc("/{tag}/{session}/status", statusChart)
	router.HandleFunc("/{tag}/{session}/daily", dailyChart)
	router.HandleFunc("/{tag}/{session}/speed", fanSpeedChart)
	router.HandleFunc("/{tag}/{session}/rotor", fanRotorChart)
	router.HandleFunc("/{tag}/temperature/current", currentChart)
	router.HandleFunc("/{tag}/temperature/watermark", watermarkChart)
	router.HandleFunc("/upload", uploadHandler)
	log.Fatal(http.ListenAndServe(":4321", router))
}

func main() {
	_, _ = fmt.Fprintf(os.Stderr, chartInformation, authorInformation, versionInformation, getServerIP())
	go chartHttpServer()
	go chartFileServer()
	// loop
	select {}
}
