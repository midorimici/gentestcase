package filterer

import (
	"fmt"

	"github.com/midorimici/gentestcase/internal/condition"
	"github.com/midorimici/gentestcase/internal/model"
)

// Filterer is a filter which removes combinations which do not satisfy given constraints.
type Filterer interface {
	// Filter returns filtered combinations and all combinations with debug information.
	Filter() ([]model.Combination, []model.Combination, error)
}

type filterer struct {
	constraints  model.Constraints
	parser       condition.Parser
	combinations []model.Combination
	isDebug      bool
}

// New returns a new Filterer for given constraints and combinations.
//
// It outputs additional information for debugging when isDebug is true.
func New(constraints model.Constraints, p condition.Parser, comb []model.Combination, isDebug bool) Filterer {
	return &filterer{constraints, p, comb, isDebug}
}

func (f *filterer) Filter() ([]model.Combination, []model.Combination, error) {
	const funcName = "filterer.Filter"

	combs := []model.Combination{}
	debugCombs := []model.Combination{}
	for _, comb := range f.combinations {
		var debugComb model.Combination
		if f.isDebug {
			debugComb = copiedMap(comb)
		}

		isValidComb := true
		for _, c := range f.constraints {
			ok, resMap, err := isConstraintSatisfied(f.parser.Parse, comb, c)
			if err != nil {
				return nil, nil, fmt.Errorf("%s: %w", funcName, err)
			}

			if f.isDebug {
				for k, v := range resMap {
					debugComb[fmt.Sprintf("[debug] %s", k)] = debugBoolStr(v)
				}
			}

			if !ok {
				isValidComb = false
				break
			}
		}

		if isValidComb {
			combs = append(combs, comb)
		}

		if f.isDebug {
			debugComb["[debug] valid"] = debugBoolStr(isValidComb)
			debugCombs = append(debugCombs, debugComb)
		}
	}
	return combs, debugCombs, nil
}

func copiedMap(m map[string]string) map[string]string {
	dest := map[string]string{}
	for k, v := range m {
		dest[k] = v
	}
	return dest
}

func debugBoolStr(b bool) string {
	if b {
		return "o"
	}

	return "x"
}

func isConstraintSatisfied(parse func(model.Combination, string) (bool, error), cmb model.Combination, cns model.Constraint) (bool, map[string]bool, error) {
	const funcName = "isConstraintSatisfied"

	if cns.If == "" && cns.OnlyIf == "" {
		return false, nil, fmt.Errorf("%s: either if or only_if property is required", funcName)
	}

	if cns.OnlyIf == "" {
		return isIfConditionSatisfied(parse, cmb, cns)
	}

	if cns.If == "" {
		return isOnlyIfConditionSatisfied(parse, cmb, cns)
	}

	return false, nil, fmt.Errorf("%s: cannot specify both if and only_if", funcName)
}

func isIfConditionSatisfied(parse func(model.Combination, string) (bool, error), cmb model.Combination, cns model.Constraint) (bool, map[string]bool, error) {
	const funcName = "isIfConditionSatisfied"

	constraintID := cns.ID
	if constraintID == "" {
		constraintID = fmt.Sprintf("condition of if: %s", cns.If)
	}

	// Check if the combination is related to the condition
	isSatisfied, err := parse(cmb, cns.If)
	if err != nil {
		return false, nil, fmt.Errorf("%s: %w", funcName, err)
	}

	if isSatisfied {
		// Check if the combination satisfies then condition
		ok, err := parse(cmb, cns.Then)
		if err != nil {
			return false, nil, fmt.Errorf("%s: %w", funcName, err)
		}

		resMap := map[string]bool{
			fmt.Sprintf("%s\n-> if: %s", constraintID, cns.If):     isSatisfied,
			fmt.Sprintf("%s\n-> then: %s", constraintID, cns.Then): ok,
			constraintID: ok,
		}

		if ok {
			return true, resMap, nil
		}

		return false, resMap, nil
	}

	// Irrelevant condition is skipped
	if cns.Else == "" {
		resMap := map[string]bool{
			fmt.Sprintf("%s\n-> if: %s", constraintID, cns.If): isSatisfied,
			constraintID: true,
		}
		return true, resMap, nil
	}

	// Check else condition
	isElseSatisfied, err := parse(cmb, cns.Else)
	if err != nil {
		return false, nil, fmt.Errorf("%s: %w", funcName, err)
	}

	resMap := map[string]bool{
		fmt.Sprintf("%s\n-> if: %s", constraintID, cns.If):     isSatisfied,
		fmt.Sprintf("%s\n-> else: %s", constraintID, cns.Then): isElseSatisfied,
		constraintID: isElseSatisfied,
	}

	if isElseSatisfied {
		return true, resMap, nil
	}

	return false, resMap, nil
}

func isOnlyIfConditionSatisfied(parse func(model.Combination, string) (bool, error), cmb model.Combination, cns model.Constraint) (bool, map[string]bool, error) {
	const funcName = "isOnlyIfConditionSatisfied"

	constraintID := cns.ID
	if constraintID == "" {
		constraintID = fmt.Sprintf("condition of only_if: %s", cns.OnlyIf)
	}

	// Check if the combination is related to the condition
	isSatisfied, err := parse(cmb, cns.Then)
	if err != nil {
		return false, nil, fmt.Errorf("%s: %w", funcName, err)
	}

	if !isSatisfied {
		// Irrelevant condition is skipped
		if cns.Else == "" {
			resMap := map[string]bool{
				fmt.Sprintf("%s\n-> then: %s", constraintID, cns.Then): isSatisfied,
				constraintID: true,
			}
			return true, resMap, nil
		}

		shouldCheck, err := parse(cmb, cns.Else)
		if err != nil {
			return false, nil, fmt.Errorf("%s: %w", funcName, err)
		}

		// Irrelevant condition is skipped
		if !shouldCheck {
			resMap := map[string]bool{
				fmt.Sprintf("%s\n-> then: %s", constraintID, cns.Then): isSatisfied,
				fmt.Sprintf("%s\n-> else: %s", constraintID, cns.Else): shouldCheck,
				constraintID: true,
			}
			return true, resMap, nil
		}
	}

	// Check if the combination satisfies the condition
	ok, err := parse(cmb, cns.OnlyIf)
	if err != nil {
		return false, nil, fmt.Errorf("%s: %w", funcName, err)
	}

	isValid := (isSatisfied && ok) || (!isSatisfied && !ok)

	resMap := map[string]bool{
		fmt.Sprintf("%s\n-> then: %s", constraintID, cns.Then):      isSatisfied,
		fmt.Sprintf("%s\n-> only_if: %s", constraintID, cns.OnlyIf): ok,
		constraintID: isValid,
	}

	if isValid {
		return true, resMap, nil
	}

	return false, resMap, nil
}
