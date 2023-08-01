package loader

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/midorimici/gentestcase/internal/model"
)

// Result contains read data from a YAML file.
type Result struct {
	// Data contains factors, conditions and constraints.
	Data *model.Data

	// OrderedFactors is a list of factors sorted by ascending order of appearance.
	OrderedFactors []string

	// LevelOrders is a map with factor IDs as keys and level and order pairs as values.
	LevelOrders map[string]map[string]int
}

// Loader is used to load data from a YAML file.
type Loader interface {
	// Load reads data from a YAML file.
	Load() (*Result, error)
}

type loader struct {
	in io.Reader
}

// New returns a new Loader for a given input source.
func New(in io.Reader) Loader {
	return &loader{in}
}

func (l *loader) Load() (*Result, error) {
	const funcName = "loader.Load"

	bytes, err := readInput(l.in)
	if err != nil {
		log.Fatal(err)
	}

	d := &model.Data{}

	if err := yaml.Unmarshal(bytes, &d); err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	fileStr := string(bytes)
	factors, err := orderedFactors(fileStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	opOrds, err := levelOrders(fileStr, factors)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	return &Result{d, factors, opOrds}, nil
}

func readInput(f io.Reader) ([]byte, error) {
	const funcName = "readInput"

	Result, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	return Result, nil
}

func orderedFactors(fileStr string) ([]string, error) {
	re := regexp.MustCompile(`\n  (\w+):\n`)
	matches := re.FindAllStringSubmatch(fileStr, -1)
	if matches == nil {
		return nil, fmt.Errorf("orderedFactors: failed")
	}

	orders := []string{}
	for _, m := range matches {
		orders = append(orders, m[1])
	}

	return orders, nil
}

func levelOrders(fileStr string, factors []string) (map[string]map[string]int, error) {
	re := regexp.MustCompile(`\n      (\w+):`)
	matches := re.FindAllStringSubmatch(fileStr, -1)
	if matches == nil {
		return nil, fmt.Errorf("levelOrders: failed")
	}

	opLocs := re.FindAllStringIndex(fileStr, -1)
	if opLocs == nil {
		return nil, fmt.Errorf("levelOrders: failed")
	}

	elemRe := regexp.MustCompile(fmt.Sprintf(`\n  (?:%s):\n`, strings.Join(factors, "|")))
	elemLocs := elemRe.FindAllStringIndex(fileStr, -1)
	if elemLocs == nil {
		return nil, fmt.Errorf("levelOrders: failed")
	}

	orders := map[string]map[string]int{}
	for _, e := range factors {
		orders[e] = map[string]int{}
	}

	j := 0
	for i, m := range matches {
		if l := opLocs[i][0]; elemLocs[j][0] < l {
			if j < len(factors)-1 && l > elemLocs[j+1][0] {
				j++
			}

			orders[factors[j]][m[1]] = i
		}
	}

	return orders, nil
}
