package exporter

import (
	"encoding/csv"
	"fmt"
	"io"
)

type Exporter interface {
	ExportCSV(table [][]string) error
}

type exporter struct {
	out     io.Writer
	headers []string
}

func New(out io.Writer, headers []string) Exporter {
	return &exporter{out, headers}
}

func (e *exporter) ExportCSV(table [][]string) error {
	const funcName = "ExportCSV"

	w := csv.NewWriter(e.out)
	if err := w.Write(e.headers); err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}

	if err := w.WriteAll(table); err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}
	return nil
}
