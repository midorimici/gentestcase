package sorter

import (
	"sort"

	"github.com/midorimici/gentestcase/internal/model"
)

// Sorter is used to sort combinations.
type Sorter interface {
	// Sort sorts given combinations in ascending order of appearance in the source YAML file.
	Sort(combinations []model.Combination) []model.Combination
}

type sorter struct {
	factors        model.Factors
	orderedFactors []string
	levelOrders    map[string]map[string]int
}

// New returns a new Sorter.
func New(factors model.Factors, orderedFactors []string, levelOrders map[string]map[string]int) Sorter {
	return &sorter{factors, orderedFactors, levelOrders}
}

func (s *sorter) Sort(c []model.Combination) []model.Combination {
	combs := append([]model.Combination{}, c...)
	sort.Slice(combs, func(i, j int) bool {
		ci := combs[i]
		cj := combs[j]

		for _, e := range s.orderedFactors {
			oi := ci[e]
			oj := cj[e]

			oiOrd := s.levelOrders[e][oi]
			ojOrd := s.levelOrders[e][oj]

			if oiOrd < ojOrd {
				return true
			} else if oiOrd > ojOrd {
				return false
			}
		}

		return false
	})
	return combs
}
