package sorter

import (
	"sort"

	"integtest/internal/model"
)

type Sorter interface {
	Sort(combinations []model.Combination) []model.Combination
}

type sorter struct {
	cases           model.Cases
	orderedElements []string
	optionOrders    map[string]int
}

func New(cases model.Cases, orderedElements []string, optionOrders map[string]int) Sorter {
	return &sorter{cases, orderedElements, optionOrders}
}

func (s *sorter) Sort(c []model.Combination) []model.Combination {
	combs := append([]model.Combination{}, c...)
	sort.Slice(combs, func(i, j int) bool {
		ci := combs[i]
		cj := combs[j]

		for _, e := range s.orderedElements {
			oi := ci[e]
			oj := cj[e]

			oiOrd := s.optionOrders[oi]
			ojOrd := s.optionOrders[oj]

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
