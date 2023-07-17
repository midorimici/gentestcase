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

type data struct {
	Data            *model.Data
	OrderedElements []string
	OptionOrders    map[string]map[string]int
}

type Loader interface {
	Load() (*data, error)
}

type loader struct {
	in io.Reader
}

func New(in io.Reader) Loader {
	return &loader{in}
}

func (l *loader) Load() (*data, error) {
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
	elements, err := orderedElements(fileStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	opOrds, err := optionOrders(fileStr, elements)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	return &data{d, elements, opOrds}, nil
}

func readInput(f io.Reader) ([]byte, error) {
	const funcName = "readInput"

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	return data, nil
}

func orderedElements(fileStr string) ([]string, error) {
	re := regexp.MustCompile(`\n  (\w+):\n`)
	matches := re.FindAllStringSubmatch(fileStr, -1)
	if matches == nil {
		return nil, fmt.Errorf("orderedElements: failed")
	}

	orders := []string{}
	for _, m := range matches {
		orders = append(orders, m[1])
	}

	return orders, nil
}

func optionOrders(fileStr string, elements []string) (map[string]map[string]int, error) {
	re := regexp.MustCompile(`\n      (\w+):\n`)
	matches := re.FindAllStringSubmatch(fileStr, -1)
	if matches == nil {
		return nil, fmt.Errorf("optionOrders: failed")
	}

	opLocs := re.FindAllStringIndex(fileStr, -1)
	if opLocs == nil {
		return nil, fmt.Errorf("optionOrders: failed")
	}

	elemRe := regexp.MustCompile(fmt.Sprintf(`\n  (?:%s):\n`, strings.Join(elements, "|")))
	elemLocs := elemRe.FindAllStringIndex(fileStr, -1)
	if elemLocs == nil {
		return nil, fmt.Errorf("optionOrders: failed")
	}

	orders := map[string]map[string]int{}
	for _, e := range elements {
		orders[e] = map[string]int{}
	}

	j := 0
	for i, m := range matches {
		if l := opLocs[i][0]; elemLocs[j][0] < l {
			if j < len(elements)-1 && l > elemLocs[j+1][0] {
				j++
			}

			orders[elements[j]][m[1]] = i
		}
	}

	return orders, nil
}
