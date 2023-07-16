package main

import (
	"log"

	"integtest/internal/condition"
	"integtest/internal/converter"
	"integtest/internal/exporter"
	"integtest/internal/filterer"
	"integtest/internal/generator"
	"integtest/internal/loader"
	"integtest/internal/sorter"
)

func main() {
	l := loader.New()
	d, err := l.Load()
	if err != nil {
		log.Fatal(err)
	}

	g := generator.New(d.Cases)
	cs := g.Generate()
	p := condition.NewParser(d.Cases)
	f := filterer.New(d.Cases, p, cs)
	fcs, err := f.Filter()
	if err != nil {
		log.Fatal(err)
	}

	s := sorter.New(d.Cases, d.OrderedElements, d.OptionOrders)
	scs := s.Sort(fcs)

	cv := converter.New(d.Cases, d.OrderedElements)
	t := cv.ConvertCombinationMapsToTable(scs)

	e := exporter.New(d.OrderedElements)
	if err := e.ExportCSV(t); err != nil {
		log.Fatal(err)
	}
}
