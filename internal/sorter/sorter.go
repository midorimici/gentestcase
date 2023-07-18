package sorter

import (
	"sort"

	"github.com/midorimici/gentestcase/internal/model"
)

type Sorter interface {
	Sort(combinations []model.Combination) []model.Combination
}

type sorter struct {
	elements        model.Factors
	orderedElements []string
	optionOrders    map[string]map[string]int
}

func New(elements model.Factors, orderedElements []string, optionOrders map[string]map[string]int) Sorter {
	return &sorter{elements, orderedElements, optionOrders}
}

func (s *sorter) Sort(c []model.Combination) []model.Combination {
	combs := append([]model.Combination{}, c...)
	sort.Slice(combs, func(i, j int) bool {
		ci := combs[i]
		cj := combs[j]

		for _, e := range s.orderedElements {
			oi := ci[e]
			oj := cj[e]

			oiOrd := s.optionOrders[e][oi]
			ojOrd := s.optionOrders[e][oj]

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
