package generator

import "github.com/midorimici/gentestcase/internal/model"

// Generator is used to generate all possible combinations.
type Generator interface {
	// Generate returns all possible combinations.
	Generate() []model.Combination
}

type generator struct {
	factors model.Factors
}

// New returns a new Generator for given factor and levels.
func New(c model.Factors) Generator {
	return &generator{c}
}

func (g *generator) Generate() []model.Combination {
	return g.combinations()
}

func (g *generator) combinations() []model.Combination {
	maxIndices := []int{}
	levels := [][]string{}
	indexKeyMap := map[int]string{}
	index := 0
	for k, v := range g.factors {
		maxIndex := len(v.Levels) - 1
		maxIndices = append(maxIndices, maxIndex)
		lvs := []string{}
		for id := range v.Levels {
			lvs = append(lvs, id)
		}
		levels = append(levels, lvs)

		indexKeyMap[index] = k
		index++
	}

	t := combTable(len(g.factors), maxIndices, levels)

	return tableToMapSlice(t, indexKeyMap)
}

func sum(a []int) int {
	var sum int
	for _, v := range a {
		sum += v
	}
	return sum
}

func combTable(length int, maxIndices []int, levels [][]string) [][]string {
	maxIndicesSum := sum(maxIndices)
	combs := [][]string{}
	counter := make([]int, length)
	for sum(counter) < maxIndicesSum {
		combs = append(combs, levelsByCounter(levels, counter))

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
	combs = append(combs, levelsByCounter(levels, counter))

	return combs
}

func levelsByCounter(levels [][]string, counter []int) []string {
	lvs := []string{}
	for k, v := range counter {
		lvs = append(lvs, levels[k][v])
	}
	return lvs
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
