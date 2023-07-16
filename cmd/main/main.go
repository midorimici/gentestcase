package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"

	"integtest/internal/condition"
	"integtest/internal/converter"
	"integtest/internal/exporter"
	"integtest/internal/filterer"
	"integtest/internal/generator"
	"integtest/internal/model"
	"integtest/internal/sorter"
)

const elementsFileName = "elements.example.yml"

func main() {
	cases, elements, opOrds, err := loadData()
	if err != nil {
		log.Fatal(err)
	}

	g := generator.New(cases)
	cs := g.Generate()
	p := condition.NewParser(cases)
	f := filterer.New(cases, p, cs)
	fcs, err := f.Filter()
	if err != nil {
		log.Fatal(err)
	}

	s := sorter.New(cases, elements, opOrds)
	scs := s.Sort(fcs)

	cv := converter.New(cases, elements)
	t := cv.ConvertCombinationMapsToTable(scs)

	e := exporter.New(elements)
	if err := e.ExportCSV(t); err != nil {
		log.Fatal(err)
	}
}

func loadData() (model.Cases, []string, map[string]int, error) {
	const funcName = "loadData"

	bytes, err := readFile()
	if err != nil {
		log.Fatal(err)
	}

	c := model.Cases{}

	if err := yaml.Unmarshal(bytes, c); err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", funcName, err)
	}

	fileStr := string(bytes)
	elements, err := orderedElements(fileStr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", funcName, err)
	}

	opOrds, err := optionOrders(fileStr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s: %w", funcName, err)
	}

	return c, elements, opOrds, nil
}

func readFile() ([]byte, error) {
	const funcName = "readFile"

	f, err := os.Open(elementsFileName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	return data, nil
}

func orderedElements(fileStr string) ([]string, error) {
	re := regexp.MustCompile(`([^ \n]+):\n  name:`)
	matches := re.FindAllStringSubmatch(fileStr, -1)
	if matches == nil {
		return nil, fmt.Errorf("orderedElements: failed")
	}

	orders := []string{}
	for _, m := range matches {
		orders = append(orders, m[1])
	}

	return orders, nil
}

func optionOrders(fileStr string) (map[string]int, error) {
	re := regexp.MustCompile(`\n    ([^ \n]+):\n`)
	matches := re.FindAllStringSubmatch(fileStr, -1)
	if matches == nil {
		return nil, fmt.Errorf("optionOrders: failed")
	}

	orders := map[string]int{}
	for i, m := range matches {
		orders[m[1]] = i
	}

	return orders, nil
}
