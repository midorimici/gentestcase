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
	out io.Writer
}

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
