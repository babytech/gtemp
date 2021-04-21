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
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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

func drawLineDaily(w http.ResponseWriter, r *csv.Reader, fileName string) {
	titleName := "Temperature " + fileName
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
	time := make([]string, 0)
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
		time = append(time, row[0])
		for i := 1; i < len(row); i++ {
			temp[name[i]] = append(temp[name[i]], opts.LineData{Value: row[i]})
		}
	}
	line.SetXAxis(time).
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
			//fmt.Println("=============k", k, "v", v)
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
		fmt.Println("row: ", row)
		if rowLine != 0 {
			items := make([]opts.BarData, 0)
			for i := 1; i <= 30; i++ {
				items = append(items, opts.BarData{Value: row[i]})
			}
			bar.AddSeries(row[0], items).SetSeriesOptions(
				charts.WithMarkPointNameTypeItemOpts(
					opts.MarkPointNameTypeItem{Name: "Maximum", Type: "max"},
					//opts.MarkPointNameTypeItem{Name: "Average", Type: "average"},
					opts.MarkPointNameTypeItem{Name: "Minimum", Type: "min"},
				),
				charts.WithMarkPointStyleOpts(
					opts.MarkPointStyle{Label: &opts.Label{Show: true}},
				),
			)
		} else {
			bar.SetXAxis([]string{"<-40", "-40~-35", "-35~-30", "-30~-25", "-25~-20", "-20~-15", "-15~-10", "-10~-5",
				"-5~0", "0~5", "5~10", "10~15", "15~20", "20~25", "25~30", "30~35", "35~40", "40~45", "45~50",
				"50~55", "55~60", "60~65", "65~70", "70~75", "75~80", "80~85", "85~90", "90~95", "95~100", ">100"})
		}
		rowLine++
	}
	_ = bar.Render(w)
}

func drawChart(w http.ResponseWriter, r *http.Request) (*csv.Reader, *os.File, string) {
	params := mux.Vars(r)
	tag := params["tag"]
	session := params["session"]
	fmt.Printf("#####tag:%s, session:%s\n", tag, session)
	fmt.Println("======dir:", "data/"+tag+"/"+session)
	inputFile := fmt.Sprintf("./tmp/temp/data/%v/%v.csv", tag, session)
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
	rr, fs, _ := drawChart(w, r)
	if rr == nil || fs == nil {
		fmt.Println("Error: drawChart fail!")
		return
	}
	page := genPages(w, rr)
	page.SetLayout(components.PageFlexLayout)
	err := page.Render(io.MultiWriter(fs))
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

func chartHttpServer() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/", statisticsChart)
	router.HandleFunc("/{tag}/{session}/statics", statisticsChart)
	router.HandleFunc("/{tag}/{session}/daily", dailyChart)
	router.HandleFunc("/{tag}/{session}/status", statusChart)
	router.HandleFunc("/upload", uploadHandler)
	log.Fatal(http.ListenAndServe(":4321", router))
}

func main() {
	chartHttpServer()
}
