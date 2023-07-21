package converter

import "github.com/midorimici/gentestcase/internal/model"

type Converter interface {
	ConvertCombinationMapsToTable(maps []model.Combination) [][]string
}

type converter struct {
	factors        model.Factors
	orderedFactors []string
}

func New(factors model.Factors, orderedFactors []string) Converter {
	return &converter{factors, orderedFactors}
}

func (c *converter) ConvertCombinationMapsToTable(maps []model.Combination) [][]string {
	table := [][]string{}
	for _, m := range maps {
		row := []string{}
		for _, e := range c.orderedFactors {
			op := m[e]
			row = append(row, c.factors[e].Levels[op])
		}
		table = append(table, row)
	}
	return table
}
