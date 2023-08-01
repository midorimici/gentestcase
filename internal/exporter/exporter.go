package exporter

import (
	"encoding/csv"
	"fmt"
	"io"
)

// Exporter is used to export a table to CSV.
type Exporter interface {
	// ExportCSV exports a given table to a CSV file.
	ExportCSV(table [][]string) error
}

type exporter struct {
	out io.Writer
}

// New returns a new Exporter for a given output destination.
func New(out io.Writer) Exporter {
	return &exporter{out}
}

func (e *exporter) ExportCSV(table [][]string) error {
	const funcName = "ExportCSV"

	w := csv.NewWriter(e.out)
	if err := w.WriteAll(table); err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}
	return nil
}
