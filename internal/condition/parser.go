package condition

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/midorimici/gentestcase/internal/model"
)

var (
	errParseFailed         = errors.New("condition parse failed")
	errInvalidOperator     = errors.New("invalid logical operator")
	errParenthesisUnpaired = errors.New("parenthesis is not paired")
)

var (
	parRe   *regexp.Regexp
	valueRe *regexp.Regexp
	refRe   *regexp.Regexp
)

func init() {
	parRe = regexp.MustCompile(`^(!?)\((.+)\)$`)
	valueRe = regexp.MustCompile(`(!?)(\w+)\.(\w+)`)
	refRe = regexp.MustCompile(`(!?)\$(\w+)`)
}

// Parser is used to get parsed results of condition statement texts.
type Parser interface {
	// Parse returns whether a given combination satisfies a condition
	// represented by a given condition statement text.
	Parse(combination model.Combination, text string) (bool, error)
}

type parser struct {
	data *model.Data
}

// NewParser returns a new Parser for given data.
func NewParser(data *model.Data) Parser {
	return &parser{data}
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

	exp, err := parseExpression(parRemoved)
	if err != nil {
		return false, err
	}
	if exp == nil {
		return p.valueCondition(combination, parRemoved)
	}

	result, err := p.expCondition(combination, exp)
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

type expression struct {
	left     string
	right    string
	operator string
}

func parseExpression(text string) (*expression, error) {
	var exp *expression
	var depth int
	var operatorIndex int
	for i, r := range text {
		switch r {
		case '(':
			depth++

		case ')':
			depth--

		case '&', '|':
			if depth == 0 {
				operatorIndex = i
				exp = &expression{operator: fmt.Sprintf("%s%s", string(r), string(r))}
			}
		}

		if exp != nil {
			break
		}
	}

	// Just a value without logical operator
	if exp == nil {
		return nil, nil
	}

	leftEndIndex := operatorIndex - 2
	rightStartIndex := operatorIndex + 3
	if text[leftEndIndex+1:rightStartIndex] != fmt.Sprintf(" %s ", exp.operator) {
		return nil, fmt.Errorf(`parsing %s: %w`, text, errParseFailed)
	}

	exp.left = text[:leftEndIndex+1]
	exp.right = text[rightStartIndex:]
	return exp, nil
}

func (p *parser) valueCondition(combination model.Combination, text string) (bool, error) {
	v := valueRe.FindStringSubmatch(text)
	if len(v) == 0 {
		r := refRe.FindStringSubmatch(text)
		if len(r) == 0 {
			return false, fmt.Errorf(`parsing "%s": %w`, text, errParseFailed)
		}
		ref := r[2]
		result, err := p.refCondition(combination, ref)
		if err != nil {
			return false, fmt.Errorf(`parsing "%s": %w`, text, err)
		}
		if r[1] == "!" {
			return !result, nil
		}
		return result, nil
	}

	factor := v[2]
	level := v[3]

	e, ok := p.data.Factors[factor]
	if !ok {
		return false, fmt.Errorf(`parsing "%s": %w`, text, fmt.Errorf(`factor "%s" is not defined`, factor))
	}

	_, ok = e.Levels[level]
	if !ok {
		return false, fmt.Errorf(`parsing "%s": %w`, text, fmt.Errorf(`level "%s" in factor %s is not defined`, level, factor))
	}

	result := level == combination[factor]
	if v[1] == "!" {
		return !result, nil
	}
	return result, nil
}

func (p *parser) refCondition(combination model.Combination, ref string) (bool, error) {
	const funcName = "refCondition"

	exp, ok := p.data.Conditions[ref]
	if !ok {
		return false, fmt.Errorf(`%s: condition "%s" is not defined`, funcName, ref)
	}
	result, err := p.condition(combination, exp)
	if err != nil {
		return false, fmt.Errorf("%s: %w", funcName, err)
	}
	return result, nil
}

func (p *parser) expCondition(combination model.Combination, exp *expression) (bool, error) {
	leftResult, err := p.condition(combination, exp.left)
	if err != nil {
		return false, err
	}

	rightResult, err := p.condition(combination, exp.right)
	if err != nil {
		return false, err
	}

	var result bool
	switch exp.operator {
	case "&&":
		result = leftResult && rightResult

	case "||":
		result = leftResult || rightResult

	default:
		return false, fmt.Errorf(`parsing "%s": %w`, exp.operator, errInvalidOperator)
	}
	return result, nil
}
