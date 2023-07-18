package converter

import "github.com/midorimici/gentestcase/internal/model"

type Converter interface {
	ConvertCombinationMapsToTable(maps []model.Combination) [][]string
}

type converter struct {
	elements        model.Factors
	orderedElements []string
}

func New(elements model.Factors, orderedElements []string) Converter {
	return &converter{elements, orderedElements}
}

func (c *converter) ConvertCombinationMapsToTable(maps []model.Combination) [][]string {
	table := [][]string{}
	for _, m := range maps {
		row := []string{}
		for _, e := range c.orderedElements {
			op := m[e]
			row = append(row, c.elements[e].Levels[op].Name)
		}
		table = append(table, row)
	}
	return table
}
