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
			ok, err := isConstraintSatisfied(f.parser.Parse, comb, c)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", funcName, err)
			}
			if !ok {
				isValidComb = false
				break
			}
		}
		if isValidComb {
			combs = append(combs, comb)
		}
	}
	return combs, nil
}

func isConstraintSatisfied(parse func(model.Combination, string) (bool, error), cmb model.Combination, cns model.Constraint) (bool, error) {
	const funcName = "isConstraintSatisfied"

	if cns.If == "" && cns.OnlyIf == "" {
		return false, fmt.Errorf("%s: either if or only_if property is required", funcName)
	}

	if cns.If == "" {
		return isOnlyIfConditionSatisfied(parse, cmb, cns)
	}

	// TODO: if condition
	return false, nil
}

func isOnlyIfConditionSatisfied(parse func(model.Combination, string) (bool, error), cmb model.Combination, cns model.Constraint) (bool, error) {
	const funcName = "isOnlyIfConditionSatisfied"

	// Check if the combination is related to the condition
	isSatisfied, err := parse(cmb, cns.Then)
	if err != nil {
		return false, fmt.Errorf("%s: %w", funcName, err)
	}

	if !isSatisfied {
		if cns.Else == "" {
			return true, nil
		}

		shouldCheck, err := parse(cmb, cns.Else)
		if err != nil {
			return false, fmt.Errorf("%s: %w", funcName, err)
		}

		// Irrelevant condition is skipped
		if !shouldCheck {
			return true, nil
		}
	}

	// Check if the combination satisfies the condition
	ok, err := parse(cmb, cns.OnlyIf)
	if err != nil {
		return false, fmt.Errorf("%s: %w", funcName, err)
	}

	if (isSatisfied && ok) || (!isSatisfied && !ok) {
		return true, nil
	}

	return false, nil
}
