package condition

import (
	"errors"
	"fmt"
	"regexp"

	"integtest/internal/model"
)

var (
	errParseFailed         = errors.New("condition parse failed")
	errInvalidOperator     = errors.New("invalid logical operator")
	errParenthesisUnpaired = errors.New("parenthesis is not paired")
)

var (
	parRe   *regexp.Regexp
	expRe   *regexp.Regexp
	valueRe *regexp.Regexp
)

func init() {
	parRe = regexp.MustCompile(`^(!?)\((.+)\)$`)
	s := `(!?\(.+\)|.+?)`
	expRe = regexp.MustCompile(fmt.Sprintf(`^%s ([&|]{2}) %s$`, s, s))
	valueRe = regexp.MustCompile(`(!?)(.+)\.(.+)`)
}

type Parser interface {
	Parse(combination model.Combination, text string) (bool, error)
}

type parser struct {
	cases model.Cases
}

func NewParser(cs model.Cases) Parser {
	return &parser{cases: cs}
}

func (p *parser) Parse(combination model.Combination, text string) (bool, error) {
	const funcName = "parser.Parse"

	result, err := p.condition(combination, text)
	if err != nil {
		return false, fmt.Errorf("%s: %w", funcName, err)
	}
	return result, nil
}

func (p *parser) condition(combination model.Combination, text string) (bool, error) {
	isGrouped, err := isGrouped(text)
	if err != nil {
		return false, err
	}

	parRemoved := text
	if isGrouped {
		parRemoved = parRe.ReplaceAllString(text, "$2")
	}

	matches := expRe.FindStringSubmatch(parRemoved)
	if len(matches) == 0 {
		return p.valueCondition(combination, parRemoved)
	}

	result, err := p.expCondition(combination, matches)
	if err != nil {
		return false, err
	}

	if isGrouped {
		if parMatches := parRe.FindStringSubmatch(text); len(parMatches) > 0 && parMatches[1] == "!" {
			return !result, nil
		}
	}
	return result, nil
}

func isGrouped(text string) (bool, error) {
	var depth int
	for i, r := range text {
		if i == 0 && r == '!' {
			continue
		}

		switch r {
		case '(':
			depth++

		case ')':
			depth--

		default:
			if depth == 0 {
				return false, nil
			}
		}
	}

	if depth != 0 {
		return false, fmt.Errorf(`parsing "%s": %w`, text, errParenthesisUnpaired)
	}

	return true, nil
}

func (p *parser) valueCondition(combination model.Combination, text string) (bool, error) {
	v := valueRe.FindStringSubmatch(text)
	if len(v) == 0 {
		return false, fmt.Errorf(`parsing "%s": %w`, text, errParseFailed)
	}

	element := v[2]
	option := v[3]

	e, ok := p.cases[element]
	if !ok {
		return false, fmt.Errorf(`parsing "%s": %w`, text, fmt.Errorf(`element "%s" is not defined`, element))
	}

	_, ok = e.Options[option]
	if !ok {
		return false, fmt.Errorf(`parsing "%s": %w`, text, fmt.Errorf(`option "%s" in element %s is not defined`, option, element))
	}

	result := option == combination[element]
	if v[1] == "!" {
		return !result, nil
	}
	return result, nil
}

func (p *parser) expCondition(combination model.Combination, matches []string) (bool, error) {
	left := matches[1]
	operator := matches[2]
	right := matches[3]

	leftResult, err := p.condition(combination, left)
	if err != nil {
		return false, err
	}

	rightResult, err := p.condition(combination, right)
	if err != nil {
		return false, err
	}

	var result bool
	switch operator {
	case "&&":
		result = leftResult && rightResult

	case "||":
		result = leftResult || rightResult

	default:
		return false, fmt.Errorf(`parsing "%s": %w`, operator, errInvalidOperator)
	}
	return result, nil
}
