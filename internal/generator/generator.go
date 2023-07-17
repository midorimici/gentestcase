package generator

import "github.com/midorimici/gentestcase/internal/model"

type Generator interface {
	Generate() []model.Combination
}

type generator struct {
	elements model.Elements
}

func New(c model.Elements) Generator {
	return &generator{c}
}

func (g *generator) Generate() []model.Combination {
	return g.combinations()
}

func (g *generator) combinations() []model.Combination {
	maxIndices := []int{}
	options := [][]string{}
	indexKeyMap := map[int]string{}
	index := 0
	for k, v := range g.elements {
		maxIndex := len(v.Options) - 1
		maxIndices = append(maxIndices, maxIndex)
		ops := []string{}
		for id := range v.Options {
			ops = append(ops, id)
		}
		options = append(options, ops)

		indexKeyMap[index] = k
		index++
	}

	t := combTable(len(g.elements), maxIndices, options)

	return tableToMapSlice(t, indexKeyMap)
}

func sum(a []int) int {
	var sum int
	for _, v := range a {
		sum += v
	}
	return sum
}

func combTable(length int, maxIndices []int, options [][]string) [][]string {
	maxIndicesSum := sum(maxIndices)
	combs := [][]string{}
	counter := make([]int, length)
	for sum(counter) < maxIndicesSum {
		combs = append(combs, optionsByCounter(options, counter))

		i := length - 1
		counter[i]++
		if counter[i] > maxIndices[i] {
			for j := i; j >= 0; j-- {
				if j > 0 && counter[j] > maxIndices[j] {
					counter[j] = 0
					counter[j-1]++
				}
			}
		}
	}
	combs = append(combs, optionsByCounter(options, counter))

	return combs
}

func optionsByCounter(options [][]string, counter []int) []string {
	ops := []string{}
	for k, v := range counter {
		ops = append(ops, options[k][v])
	}
	return ops
}

func tableToMapSlice(t [][]string, indexKeyMap map[int]string) []model.Combination {
	maps := []model.Combination{}
	for _, r := range t {
		m := model.Combination{}
		for i := 0; i < len(r); i++ {
			k := indexKeyMap[i]
			m[k] = r[i]
		}
		maps = append(maps, m)
	}
	return maps
}
