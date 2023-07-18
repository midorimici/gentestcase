package exporter

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/midorimici/gentestcase/internal/model"
)

type Exporter interface {
	ExportCSV(table [][]string) error
}

type exporter struct {
	out      io.Writer
	elements model.Factors
	headers  []string
}

func New(out io.Writer, elements model.Factors, headers []string) Exporter {
	return &exporter{out, elements, headers}
}

func (e *exporter) ExportCSV(table [][]string) error {
	const funcName = "ExportCSV"

	headers := []string{}
	for _, h := range e.headers {
		headers = append(headers, e.elements[h].Name)
	}

	w := csv.NewWriter(e.out)
	if err := w.Write(headers); err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}

	if err := w.WriteAll(table); err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}
	return nil
}
