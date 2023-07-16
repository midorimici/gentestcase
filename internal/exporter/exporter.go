package exporter

import (
	"encoding/csv"
	"fmt"
	"os"
)

const CSVFileName = "data.csv"

type Exporter interface {
	ExportCSV(table [][]string) error
}

type exporter struct {
	headers []string
}

func New(headers []string) Exporter {
	return &exporter{headers}
}

func (e *exporter) ExportCSV(table [][]string) error {
	const funcName = "ExportCSV"

	f, err := os.Create(CSVFileName)
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	err = w.Write(e.headers)
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}

	err = w.WriteAll(table)
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}
	return nil
}
