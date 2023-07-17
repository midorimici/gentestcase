package filterer

import (
	"github.com/midorimici/gentestcase/internal/condition"
	"github.com/midorimici/gentestcase/internal/model"
)

type Filterer interface {
	Filter() ([]model.Combination, error)
}

type filterer struct {
	elements     model.Elements
	parser       condition.Parser
	combinations []model.Combination
}

func New(elements model.Elements, p condition.Parser, comb []model.Combination) Filterer {
	return &filterer{elements, p, comb}
}

func (f *filterer) Filter() ([]model.Combination, error) {
	const funcName = "filterer.Filter"

	combs := []model.Combination{}
	for _, comb := range f.combinations {
		isValidComb := true
		for element, option := range comb {
			cond := f.elements[element].Options[option].If
			if cond == "" {
				continue
			}

			ok, err := f.parser.Parse(comb, cond)
			if err != nil {
				return nil, err
			}
			if ok {
				continue
			}
			isValidComb = false
			break
		}
		if isValidComb {
			combs = append(combs, comb)
		}
	}
	return combs, nil
}
