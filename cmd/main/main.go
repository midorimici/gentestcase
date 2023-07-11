package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"integtest/internal/model"
)

const elementsFileName = "elements.yml"

func main() {
	d, err := loadData()
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("%#v\n", len(d))
}

func loadData() (model.Cases, error) {
	const funcName = "loadData"

	bytes, err := readFile();
	if err != nil {
		log.Fatal(err)
	}
	
	c := model.Cases{}

	if err := yaml.Unmarshal(bytes, c); err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}
	
	return c, nil
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
