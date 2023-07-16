package loader

import (
	"fmt"
	"io"
	"log"
	"regexp"

	"gopkg.in/yaml.v3"

	"github.com/midorimici/gentestcase/internal/model"
)

type data struct {
	Cases           model.Cases
	OrderedElements []string
	OptionOrders    map[string]int
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

	c := model.Cases{}

	if err := yaml.Unmarshal(bytes, c); err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	fileStr := string(bytes)
	elements, err := orderedElements(fileStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	opOrds, err := optionOrders(fileStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	return &data{c, elements, opOrds}, nil
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
	re := regexp.MustCompile(`([^ \n]+):\n  name:`)
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

func optionOrders(fileStr string) (map[string]int, error) {
	re := regexp.MustCompile(`\n    ([^ \n]+):\n`)
	matches := re.FindAllStringSubmatch(fileStr, -1)
	if matches == nil {
		return nil, fmt.Errorf("optionOrders: failed")
	}

	orders := map[string]int{}
	for i, m := range matches {
		orders[m[1]] = i
	}

	return orders, nil
}
