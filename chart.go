package main

import (
	"encoding/csv"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
)

func drawLineDaily(w http.ResponseWriter, r *csv.Reader) {
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
			Title: "Temperature daily",
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
	time := make([]string, 0)
	temp := make(map[string]([]opts.LineData))
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
		fmt.Println("row: ",row)
		time = append(time,row[0])
		for i:=1; i<len(row); i++ {
			temp[name[i]] = append(temp[name[i]],opts.LineData{Value: row[i]})
		}
	}

	line.SetXAxis(time).
	SetSeriesOptions(charts.WithLineChartOpts(
		opts.LineChart{
			Smooth: true,
		}),
	)
	for k, v := range temp {
		if k != "time" {
			line.AddSeries(k, v)
			fmt.Println("=============k", k, "v", v)
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
			Title:    "temp sensor statistics",
			Subtitle: "MF14 ISR2103.46",
			Left:     "10%",
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
		fmt.Println("row: ",row)
		if rowLine != 0 {
			items := make([]opts.BarData, 0)
			for i := 1; i <= 30; i++ {
				items = append(items, opts.BarData{Value: row[i]})
			}
			bar.AddSeries(row[0], items)
		} else {
			bar.SetXAxis([]string{"<-40", "-40~-35", "-35~-30", "-30~-25", "-25~-20", "-20~-15", "-15~-10", "-10~-5",
				"-5~0", "0~5", "5~10", "10~15", "15~20", "20~25", "25~30", "30~35", "35~40", "40~45", "45~50",
				"50~55", "55~60", "60~65", "65~70", "70~75", "75~80", "80~85", "85~90","90~95", "95~100", ">100"})
		}
		rowLine++
	}
	_ = bar.Render(w)
}

func drawLine(w http.ResponseWriter, r *csv.Reader) {
	line := charts.NewLine()
	line.SetGlobalOptions(
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
			Title:    "temp sensor statistics",
			Subtitle: "MF14 ISR2103.46",
			Left:     "10%",
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
		fmt.Println("row: ",row)
		if rowLine != 0 {
			items := make([]opts.LineData, 0)
			for i := 1; i <= 30; i++ {
				items = append(items, opts.LineData{Value: row[i]})
			}
			line.AddSeries(row[0], items)
		} else {
			line.SetXAxis([]string{"<-40", "-40~-35", "-35~-30", "-30~-25", "-25~-20", "-20~-15", "-15~-10", "-10~-5",
				"-5~0", "0~5", "5~10", "10~15", "15~20", "20~25", "25~30", "30~35", "35~40", "40~45", "45~50",
				"50~55", "55~60", "60~65", "65~70", "70~75", "75~80", "80~85", "85~90","90~95", "95~100", ">100"})
		}
		rowLine++
	}
	_ = line.Render(w)
}

func statisticsChart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tag := params["tag"]
	session := params["session"]
	fmt.Printf("#####tag:%s, session:%s\n", tag, session)
	fmt.Println("======dir:", "data/"+tag+"/"+session)
	inputFile := fmt.Sprintf("/tmp/temp/data/%v/%v.csv", tag, session)
	fmt.Println("inputFile:", inputFile)

	fs, err := os.Open(inputFile)
	if err != nil {
		http.Error(w, "File not found.", 404)
		return
	}
	defer fs.Close()
	rr := csv.NewReader(fs)
	rr.FieldsPerRecord = -1

	switch m {
	case "bar":
		drawBar(w, rr)
	case "line":
		drawLine(w, rr)
	default:
		drawBar(w, rr)
	}
}

func dailyChart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tag := params["tag"]
	session := params["session"]
	fmt.Printf("#####tag:%s, session:%s\n", tag, session)
	fmt.Println("======dir:", "data/"+tag+"/"+session)
	inputFile := fmt.Sprintf("/tmp/temp/data/%v/%v.csv", tag, session)
	fmt.Println("inputFile:", inputFile)

	fs, err := os.Open(inputFile)
	if err != nil {
		http.Error(w, "File not found.", 404)
		return
	}
	defer fs.Close()
	rr := csv.NewReader(fs)
	rr.FieldsPerRecord = -1

	drawLineDaily(w, rr)
}

func chartHttpServer() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/", statisticsChart)
	router.HandleFunc("/{tag}/{session}/statics", statisticsChart)
	router.HandleFunc("/{tag}/{session}/daily", dailyChart)
	_ = http.ListenAndServe(":4321", router)
}

func ChartMain() {
	go chartHttpServer()
}