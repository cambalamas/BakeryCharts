package main

import (
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// CPlotData is the expected format of data for this plotter
type CPlotData map[string][]float64

// CPlot is a custom object to make a plot
type CPlot struct {
	title  string
	labelX string
	labelY string
	sizeX  float64
	sizeY  float64
	data   CPlotData
}

// NewCPlot is a CPlot constructor
func NewCPlot(title, labelX, labelY string, sizeX, sizeY float64, data CPlotData) CPlot {
	return CPlot{
		title:  title,
		labelX: labelX,
		labelY: labelY,
		sizeX:  sizeX,
		sizeY:  sizeY,
		data:   data,
	}
}

// MakePNG generate a png image in gived path
func (cp CPlot) MakePNG(path string) {
	p, err := plot.New()
	if err != nil {
		log.Fatalf("Could not create the plot %q: %v", cp.title, err)
	}

	p.Title.Text = cp.title
	p.Title.Font.Size = 15
	p.X.Label.Text = cp.labelX
	p.Y.Label.Text = cp.labelY
	// TODO: add X and Y axis granularity in steps marks

	// Transform data to fit into gonum.plot way
	var dataToPlot []interface{}
	for k, v := range cp.data {
		dataToPlot = append(dataToPlot, k)
		dataToPlot = append(dataToPlot, sliceToPlotter(v))
	}

	// Add data to the plot requiere unpack previous interface slice
	err = plotutil.AddLinePoints(p, dataToPlot...)
	if err != nil {
		log.Fatalf("Could not add line and points: %v", err)
	}

	if err := p.Save(vg.Length(cp.sizeX), vg.Length(cp.sizeY), path); err != nil {
		log.Fatalf("Could not create the plot %q at %q: %v", cp.title, path, err)
	}
}

// HELPERS

func sliceToPlotter(data []float64) plotter.XYs {
	pts := make(plotter.XYs, len(data))
	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = data[i]
	}
	return pts
}
