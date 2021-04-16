package main

import (
	"encoding/csv"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/gorilla/mux"
	"io"
	"math/rand"
	"net/http"
	"os"
)

func genLiquidItems(data []float32) []opts.LiquidData {
	items := make([]opts.LiquidData, 0)
	for i := 0; i < len(data); i++ {
		items = append(items, opts.LiquidData{Value: data[i]})
	}
	return items
}

func drawLiquidPin(w http.ResponseWriter, r *csv.Reader) *charts.Liquid {
	liquid := charts.NewLiquid()
	liquid.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "温度传感器温度水位",
		}),
	)
	liquid.AddSeries("liquid", genLiquidItems([]float32{0.3, 0.4, 0.5})).
		SetSeriesOptions(
			charts.WithLiquidChartOpts(opts.LiquidChart{
				IsWaveAnimation: true,
				Shape:           "pin",
			}),
		)
	_ = liquid.Render(w)
	return liquid
}

func drawGauge(w http.ResponseWriter, r *csv.Reader) *charts.Gauge {
	gauge := charts.NewGauge()
	gauge.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "温度传感器温度值"}),
	)
	gauge.AddSeries("ProjectA", []opts.GaugeData{{Name: "#1号温感当前温度", Value: 500}})
	_ = gauge.Render(w)
	return gauge
}

var (
	itemCntPie = 4
	seasons    = []string{"Spring", "Summer", "Autumn ", "Winter"}
)

func generatePieItems() []opts.PieData {
	items := make([]opts.PieData, 0)
	for i := 0; i < itemCntPie; i++ {
		items = append(items, opts.PieData{Name: seasons[i], Value: rand.Intn(100)})
	}
	return items
}

func drawPieRoseRadius(w http.ResponseWriter, r *csv.Reader) *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Rose(Radius)",
		}),
	)
	pie.AddSeries("pie", generatePieItems()).
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
	_ = pie.Render(w)
	return pie
}

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

func drawChart(w http.ResponseWriter, r *http.Request) (*csv.Reader, *os.File) {
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
		return nil, nil
	}
	rr := csv.NewReader(fs)
	rr.FieldsPerRecord = -1
	return rr, fs
}

func statisticsChart(w http.ResponseWriter, r *http.Request) {
	rr, fs := drawChart(w, r)
	switch chartMode {
	case "bar":
		drawBar(w, rr)
	case "line":
		drawLine(w, rr)
	default:
		drawBar(w, rr)
	}
	defer fs.Close()
}

func dailyChart(w http.ResponseWriter, r *http.Request) {
	rr, fs := drawChart(w, r)
	drawLineDaily(w, rr)
	defer fs.Close()
}

func genPages(w http.ResponseWriter, r *csv.Reader) *components.Page {
	page := components.NewPage()
	page.AddCharts(
		drawLiquidPin(w, r),
		drawGauge(w, r),
		drawPieRoseRadius(w, r),
	)
	return page
}

func statusChart(w http.ResponseWriter, r *http.Request) {
	rr, fs := drawChart(w, r)
	page := genPages(w, rr)
	page.SetLayout(components.PageFlexLayout)
	page.Render(io.MultiWriter(fs))
	defer fs.Close()
}

func chartHttpServer() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/", statisticsChart)
	router.HandleFunc("/{tag}/{session}/statics", statisticsChart)
	router.HandleFunc("/{tag}/{session}/daily", dailyChart)
	router.HandleFunc("/{tag}/{session}/status", statusChart)
	_ = http.ListenAndServe(":4321", router)
}

func ChartMain() {
	go chartHttpServer()
}