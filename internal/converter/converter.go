package converter

import (
	"fmt"
	"sort"

	"github.com/midorimici/gentestcase/internal/model"
)

type Converter interface {
	ConvertCombinationMapsToTable(maps []model.Combination) [][]string
}

type converter struct {
	factors        model.Factors
	orderedFactors []string
	isDebug        bool
}

func New(factors model.Factors, orderedFactors []string, isDebug bool) Converter {
	return &converter{factors, orderedFactors, isDebug}
}

func (c *converter) ConvertCombinationMapsToTable(maps []model.Combination) [][]string {
	table := [][]string{}

	debugKeys := []string{}
	if c.isDebug {
		for k := range maps[0] {
			var isFound bool
			for _, f := range c.orderedFactors {
				if k == f {
					isFound = true
					break
				}
			}
			if !isFound && k != "[debug] valid" {
				debugKeys = append(debugKeys, k)
			}
		}
	}
	sort.Slice(debugKeys, func(i, j int) bool {
		return debugKeys[i] < debugKeys[j]
	})
	debugKeys = append([]string{"[debug] valid"}, debugKeys...)

	// Header
	header := []string{}
	for _, f := range c.orderedFactors {
		h := c.factors[f].Name
		if c.isDebug {
			h = fmt.Sprintf("%s: %s", f, h)
		}
		header = append(header, h)
	}
	if c.isDebug {
		header = append(header, debugKeys...)
	}
	table = append(table, header)

	// Body
	for _, m := range maps {
		row := []string{}
		for _, f := range c.orderedFactors {
			lv := m[f]
			text := c.factors[f].Levels[lv]
			if c.isDebug {
				text = fmt.Sprintf("%s: %s", lv, text)
			}
			row = append(row, text)
		}
		if c.isDebug {
			for _, k := range debugKeys {
				row = append(row, m[k])
			}
		}
		table = append(table, row)
	}
	return table
}
