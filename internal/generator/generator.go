package generator

import "integtest/internal/model"

type Generator interface {
	Generate() ([][]string, error)
}

type generator struct {
	cases model.Cases
}

func New(c model.Cases) Generator {
	return &generator{c}
}

func (g *generator) Generate() ([][]string, error) {
	const funcName = "Generate"
	return nil, nil
}
