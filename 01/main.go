package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"

	"github.com/jung-kurt/gofpdf"
)

type imgSize struct{ x, y float64 }

const (
	basicCost = 5
	deluxCost = 6
	imgPath   = "img/"
	dataPath  = "data/"
)

var (
	pdfOptPNG = gofpdf.ImageOptions{ImageType: "png"}
	pltSize   = imgSize{x: 512., y: 512.}
)

func main() {

	// Revenues of week, month and year for each dataset
	// e.g. basicData["week"][0] is the revenue of fisrt week of basic cakes
	basicVal := processFlatSheet(parseData(dataPath+"basic.txt"), basicCost)
	deluxVal := processFlatSheet(parseData(dataPath+"delux.txt"), deluxCost)
	totalVal := processFlatSheet(parseData(dataPath+"total.txt"), 1)

	basicQuant := processFlatSheet(parseData(dataPath+"basic.txt"), 1)
	deluxQuant := processFlatSheet(parseData(dataPath+"delux.txt"), 1)

	// ---------------------------------------------------------
	// Plots generation
	// ---------------------------------------------------------

	// Profit of last year (per month)
	pltTitle := "Last year profit (per month)"
	pltData := CPlotData{
		"Basic": basicVal["month"][0:12],
		"Delux": deluxVal["month"][0:12],
		"Total": totalVal["month"][0:12],
	}
	lastYearProfitPerMonth := NewCPlot(pltTitle, "Time", "Profit", pltSize.x, pltSize.y, pltData)
	lastYearProfitPerMonth.MakePNG(imgPath + "lastYearProfitPerMonth.png")

	// Profit of last 3 months (per week)

	// #1
	month1ProfitPerWeek := NewCPlot("Last month profit (per week)", "Time", "Profit", pltSize.x, pltSize.y, CPlotData{
		"Basic": basicVal["week"][0:4],
		"Delux": deluxVal["week"][0:4],
		"Total": totalVal["week"][0:4],
	})
	month1ProfitPerWeek.MakePNG(imgPath + "month1ProfitPerWeek.png")

	// #2
	month2ProfitPerWeek := NewCPlot("2 months ago profit (per week)", "Time", "Profit", pltSize.x, pltSize.y, CPlotData{
		"Basic": basicVal["week"][4:8],
		"Delux": deluxVal["week"][4:8],
		"Total": totalVal["week"][4:8],
	})
	month2ProfitPerWeek.MakePNG(imgPath + "month2ProfitPerWeek.png")

	// #3
	month3ProfitPerWeek := NewCPlot("3 months ago profit (per week)", "Time", "Profit", pltSize.x, pltSize.y, CPlotData{
		"Basic": basicVal["week"][8:12],
		"Delux": deluxVal["week"][8:12],
		"Total": totalVal["week"][8:12],
	})
	month3ProfitPerWeek.MakePNG(imgPath + "month3ProfitPerWeek.png")

	// Cupcakes sold (per day, last month)
	cupcakesSold := NewCPlot("Cupcakes sold (per day, last month)", "Time", "Quantity", pltSize.x, pltSize.y, CPlotData{
		"Basic": basicQuant["day"][0:30],
		"Delux": deluxQuant["day"][0:30],
	})
	cupcakesSold.MakePNG(imgPath + "cupcakesSold.png")

	// ---------------------------------------------------------
	// PDF generation
	// ---------------------------------------------------------

	// How much money did I make last year?
	mlyB, mlyD, mlyT := basicVal["year"][0], deluxVal["year"][0], totalVal["year"][0]

	// How much money do I make in a typical month?
	mAvgB, mAvgD, mAvgT := dataAvg(basicVal["month"]), dataAvg(deluxVal["month"]), dataAvg(totalVal["month"])

	// ---------------------------------------------------------

	pdfName := "MatildaBakery.pdf"

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetAutoPageBreak(true, 5.)
	tr := pdf.UnicodeTranslatorFromDescriptor("") // Allow "€, ñ, ó, ù" and other non-english chars

	// Title
	pdf.SetFont("Arial", "B", 24)
	pdf.Write(12, "Matilda's bakery")
	pdf.Ln(15)

	// Q1
	pdf.SetFont("Arial", "B", 18)
	pdf.Write(9, "How much money did I make last year?")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 14)
	pdf.Write(7, tr("Basic:\t"+fmt.Sprintf("%.2f", mlyB)+" €"))
	pdf.Ln(-1)
	pdf.Write(7, tr("Deluxe:\t"+fmt.Sprintf("%.2f", mlyD)+" €"))
	pdf.Ln(-1)
	pdf.Write(7, tr("Total:\t"+fmt.Sprintf("%.2f", mlyT)+" €"))
	pdf.Ln(15)

	// Q2
	pdf.SetFont("Arial", "B", 18)
	pdf.Write(9, "How much money do I make in a typical month?")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 14)
	pdf.Write(7, tr("Basic:\t"+fmt.Sprintf("%.2f", mAvgB)+" €"))
	pdf.Ln(-1)
	pdf.Write(7, tr("Deluxe:\t"+fmt.Sprintf("%.2f", mAvgD)+" €"))
	pdf.Ln(-1)
	pdf.Write(7, tr("Total:\t"+fmt.Sprintf("%.2f", mAvgT)+" €"))
	pdf.Ln(15)

	// Charts
	pdf.SetFont("Arial", "B", 18)
	pdf.Write(7, tr("Some cool charts "))
	pdf.Ln(25)

	pw, _ := pdf.GetPageSize()
	ps := pw * 0.65
	pdfInsertImg := func(path string) {
		pdf.ImageOptions(path, 30, -1, ps, ps, true, gofpdf.ImageOptions{ImageType: "png"}, 0, "")
	}

	// Get images from its folder
	images, err := ioutil.ReadDir(imgPath)
	if err != nil {
		log.Fatalf("Could not read %q: %v", imgPath, err)
	}
	for _, img := range images {
		pdfInsertImg(imgPath + img.Name())
	}

	if err := pdf.OutputFileAndClose(pdfName); err != nil {
		log.Fatalf("Could not create PDF %q: %v", pdfName, err)
	}

	// ---------------------------------------------------------
}

func parseData(path string) []float64 {
	var output []float64

	// Load file in memory
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %q - %v", path, err)
		return output
	}
	defer file.Close()

	// Read line by line and store nums
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		var value float64
		_, err := fmt.Sscanf(sc.Text(), "%f", &value)
		if err == nil {
			output = append(output, value)
		}
	}

	if err := sc.Err(); err != nil {
		log.Fatalf("Could not scan the file %q: %v", path, err)
		return output
	}

	return output
}

func processFlatSheet(arr []float64, mult float64) map[string][]float64 {
	if mult < 1 {
		mult = 1
	}

	currWeek, currMonth, currYear := 0, 0, 0
	daysWeek, daysMonth, daysYear := 0, 0, 0

	var day []float64
	week, month, year := []float64{0.}, []float64{0.}, []float64{0.}

	for i := len(arr) - 1; i >= 0; i-- {
		currVal := arr[i] * mult

		day = append(day, currVal)

		week[currWeek] += currVal
		daysWeek++
		if daysWeek >= 7 {
			daysWeek = 0
			currWeek++
			week = append(week, 0.)
		}

		month[currMonth] += currVal
		daysMonth++
		if daysMonth >= 30 {
			daysMonth = 0
			currMonth++
			month = append(month, 0.)
		}

		year[currYear] += currVal
		daysYear++
		if daysYear >= 365 {
			daysYear = 0
			currYear++
			year = append(year, 0.)
		}
	}

	return map[string][]float64{
		"day":   day,
		"week":  week,
		"month": month,
		"year":  year,
	}
}

//
// Maths

func dataAvg(data []float64) float64 {
	var sum float64
	for _, val := range data {
		sum += val
	}
	avg := sum / float64(len(data))
	return math.Floor(avg*100) / 100 // truncate to two decimals
}

//
// Plotter helpers
