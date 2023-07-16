package converter

import "integtest/internal/model"

type Converter interface {
	ConvertCombinationMapsToTable(maps []model.Combination) [][]string
}

type converter struct {
	cases           model.Cases
	orderedElements []string
}

func New(cases model.Cases, orderedElements []string) Converter {
	return &converter{cases, orderedElements}
}

func (c *converter) ConvertCombinationMapsToTable(maps []model.Combination) [][]string {
	table := [][]string{}
	for _, m := range maps {
		row := []string{}
		for _, e := range c.orderedElements {
			op := m[e]
			row = append(row, c.cases[e].Options[op].Name)
		}
		table = append(table, row)
	}
	return table
}
