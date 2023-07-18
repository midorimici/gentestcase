package generator_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/midorimici/gentestcase/internal/generator"
	"github.com/midorimici/gentestcase/internal/model"
)

func sorted(s []model.Combination) []model.Combination {
	sort.Slice(s, func(i, j int) bool {
		si := s[i]
		sj := s[j]
		ti := si["e1"] + si["e2"] + si["e3"]
		tj := sj["e1"] + sj["e2"] + sj["e3"]

		return ti < tj
	})
	return s
}

func Test_generator_Generate(t *testing.T) {
	elems := model.Factors{
		"e1": {Levels: map[string]model.Level{"a": {}, "b": {}}},
		"e2": {Levels: map[string]model.Level{"d": {}, "e": {}, "f": {}}},
		"e3": {Levels: map[string]model.Level{"g": {}, "h": {}}},
	}

	type args struct {
		elements model.Factors
	}
	tests := []struct {
		name string
		args args
		want []model.Combination
	}{
		{
			name: "returns result as expected",
			args: args{elems},
			want: []model.Combination{
				{"e1": "a", "e2": "d", "e3": "g"},
				{"e1": "a", "e2": "d", "e3": "h"},
				{"e1": "a", "e2": "e", "e3": "g"},
				{"e1": "a", "e2": "e", "e3": "h"},
				{"e1": "a", "e2": "f", "e3": "g"},
				{"e1": "a", "e2": "f", "e3": "h"},
				{"e1": "b", "e2": "d", "e3": "g"},
				{"e1": "b", "e2": "d", "e3": "h"},
				{"e1": "b", "e2": "e", "e3": "g"},
				{"e1": "b", "e2": "e", "e3": "h"},
				{"e1": "b", "e2": "f", "e3": "g"},
				{"e1": "b", "e2": "f", "e3": "h"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := generator.New(tt.args.elements)
			got := g.Generate()
			got = sorted(got)
			want := sorted(tt.want)

			if !reflect.DeepEqual(got, want) {
				t.Errorf("generator.Generate() = %v, want %v", got, want)
			}
		})
	}
}
