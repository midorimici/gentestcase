package loader

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"

	"integtest/internal/model"
)

const elementsFileName = "elements.example.yml"

type data struct {
	Cases           model.Cases
	OrderedElements []string
	OptionOrders    map[string]int
}

type Loader interface {
	Load() (*data, error)
}

type loader struct{}

func New() Loader {
	return &loader{}
}

func (l *loader) Load() (*data, error) {
	const funcName = "loader.Load"

	bytes, err := readFile()
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

func readFile() ([]byte, error) {
	const funcName = "readFile"

	f, err := os.Open(elementsFileName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}
	defer f.Close()

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
