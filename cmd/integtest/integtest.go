package main

import (
	"flag"
	"io"
	"log"
	"os"

	"integtest/internal/condition"
	"integtest/internal/converter"
	"integtest/internal/exporter"
	"integtest/internal/filterer"
	"integtest/internal/generator"
	"integtest/internal/loader"
	"integtest/internal/sorter"
)

var (
	inputFilename  = flag.String("input", "elements.yml", "input YAML filename")
	outputFilename = flag.String("output", "data.csv", "output CSV filename")
)

func main() {
	flag.Parse()

	// Setup input reader
	var in io.Reader
	if *inputFilename == "" {
		in = os.Stdin
	} else {
		f, err := os.Open(*inputFilename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		in = f
	}

	// Load data from input
	l := loader.New(in)
	d, err := l.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Generate all combinations
	g := generator.New(d.Cases)
	cs := g.Generate()

	// Filter unnecessary cases
	p := condition.NewParser(d.Cases)
	f := filterer.New(d.Cases, p, cs)
	fcs, err := f.Filter()
	if err != nil {
		log.Fatal(err)
	}

	// Sort cases
	s := sorter.New(d.Cases, d.OrderedElements, d.OptionOrders)
	scs := s.Sort(fcs)

	// Convert cases to table
	cv := converter.New(d.Cases, d.OrderedElements)
	t := cv.ConvertCombinationMapsToTable(scs)

	// Setup output writer
	var out io.Writer
	if *outputFilename == "" {
		out = os.Stdout
	} else {
		f, err := os.Create(*outputFilename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		out = f
	}

	// Export to CSV
	e := exporter.New(out, d.OrderedElements)
	if err := e.ExportCSV(t); err != nil {
		log.Fatal(err)
	}
}
