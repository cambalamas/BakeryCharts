package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

const logPath = "./logs/"

const (
	basicCost = 5
	deluxCost = 6
)

func main() {
	dPath := func(filename string) string {
		return "data/" + strings.Title(strings.ToLower(filename)) + ".txt"
	}

	// Revenues of week, month and year for each dataset
	// e.g. basicData["week"][0] is the revenue of fisrt week of basic cakes
	basic := processFlatSheet(parseData(dPath("basic")), basicCost)
	delux := processFlatSheet(parseData(dPath("delux")), deluxCost)
	total := processFlatSheet(parseData(dPath("total")), 1)

	fmt.Printf("How much money did I make last year?\n"+
		"Basic\tDelux\tTotal\n"+
		"---------------------\n"+
		"%v\t%v\t%v\n"+
		"", basic["year"][0], delux["year"][0], total["year"][0])

	fmt.Printf("\n\n")

	fmt.Printf("How much money do I make in a typical month? (Months Average)\n"+
		"Basic\tDelux\tTotal\n"+
		"---------------------\n"+
		"%v\t%v\t%v\n"+
		"", dataAvg(basic["month"]), dataAvg(delux["month"]), dataAvg(total["month"]))

	p, err := plot.New()
	if err != nil {
		log.Fatalf("Could not create the plot: %v", err)
	}

	p.Title.Text = "Profit per month in last 12 months"
	p.Title.Font.Size = 15
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Profit"

	//  Plot per moth profit of last year -> [0:11]
	err = plotutil.AddLinePoints(p,
		"Basic", sliceToPlotter(basic["month"][0:11]),
		"Delux", sliceToPlotter(delux["month"][0:11]),
		"Total", sliceToPlotter(total["month"][0:11]),
	)
	if err != nil {
		panic(err)
	}

	if err := p.Save(512, 512, "points.png"); err != nil {
		panic(err)
	}

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
	week, month, year := []float64{0.}, []float64{0.}, []float64{0.}

	for i := len(arr) - 1; i >= 0; i-- {
		currVal := arr[i] * mult

		week[currWeek] += currVal
		daysWeek++
		month[currMonth] += currVal
		daysMonth++
		year[currYear] += currVal
		daysYear++

		// Week reset
		if daysWeek >= 7 {
			daysWeek = 0
			currWeek++
			week = append(week, 0.)
		}
		// Month reset
		if daysMonth >= 30 {
			daysMonth = 0
			currMonth++
			month = append(month, 0.)
		}
		// Year reset
		if daysYear >= 365 {
			daysYear = 0
			currYear++
			year = append(year, 0.)
		}
	}

	return map[string][]float64{
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

func sliceToPlotter(data []float64) plotter.XYs {
	pts := make(plotter.XYs, len(data))
	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = data[i]
	}
	return pts
}
