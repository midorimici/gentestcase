package filterer

import (
	"fmt"

	"github.com/midorimici/gentestcase/internal/condition"
	"github.com/midorimici/gentestcase/internal/model"
)

type Filterer interface {
	Filter() ([]model.Combination, error)
}

type filterer struct {
	constraints  model.Constraints
	parser       condition.Parser
	combinations []model.Combination
}

func New(constraints model.Constraints, p condition.Parser, comb []model.Combination) Filterer {
	return &filterer{constraints, p, comb}
}

func (f *filterer) Filter() ([]model.Combination, error) {
	const funcName = "filterer.Filter"

	combs := []model.Combination{}
	for _, comb := range f.combinations {
		isValidComb := true
		for _, c := range f.constraints {
			// Check if the combination is related to the condition
			isSatisfied, err := f.parser.Parse(comb, c.Then)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", funcName, err)
			}

			if !isSatisfied {
				if c.Else == "" {
					continue
				}

				shouldCheck, err := f.parser.Parse(comb, c.Else)
				if err != nil {
					return nil, fmt.Errorf("%s: %w", funcName, err)
				}

				// Irrelevant condition is skipped
				if !shouldCheck {
					continue
				}
			}

			// Check if the combination satisfies the condition
			ok, err := f.parser.Parse(comb, c.OnlyIf)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", funcName, err)
			}

			if (isSatisfied && ok) || (!isSatisfied && !ok) {
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
