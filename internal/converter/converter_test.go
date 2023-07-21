package converter_test

import (
	"reflect"
	"testing"

	"github.com/midorimici/gentestcase/internal/converter"
	"github.com/midorimici/gentestcase/internal/model"
)

func Test_converter_ConvertCombinationMapsToTable(t *testing.T) {
	elems := model.Factors{
		"e1": {Levels: map[string]string{"a": "A", "b": "B"}},
		"e2": {Levels: map[string]string{"d": "D", "e": "E", "f": "F"}},
		"e3": {Levels: map[string]string{"g": "G", "h": "H"}},
	}
	orderedFactors := []string{"e1", "e2", "e3"}
	want := [][]string{
		{"A", "D", "G"},
		{"A", "D", "H"},
		{"A", "E", "G"},
		{"A", "E", "H"},
		{"A", "F", "G"},
		{"A", "F", "H"},
		{"B", "D", "G"},
		{"B", "D", "H"},
		{"B", "E", "G"},
		{"B", "E", "H"},
		{"B", "F", "G"},
		{"B", "F", "H"},
	}

	type args struct {
		maps []model.Combination
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "returns result as expected",
			args: args{[]model.Combination{
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
			}},
		},
		{
			name: "returns result as expected",
			args: args{[]model.Combination{
				{"e3": "g", "e2": "d", "e1": "a"},
				{"e3": "h", "e2": "d", "e1": "a"},
				{"e3": "g", "e2": "e", "e1": "a"},
				{"e3": "h", "e2": "e", "e1": "a"},
				{"e3": "g", "e2": "f", "e1": "a"},
				{"e3": "h", "e2": "f", "e1": "a"},
				{"e3": "g", "e2": "d", "e1": "b"},
				{"e3": "h", "e2": "d", "e1": "b"},
				{"e3": "g", "e2": "e", "e1": "b"},
				{"e3": "h", "e2": "e", "e1": "b"},
				{"e3": "g", "e2": "f", "e1": "b"},
				{"e3": "h", "e2": "f", "e1": "b"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := converter.New(elems, orderedFactors)
			if got := c.ConvertCombinationMapsToTable(tt.args.maps); !reflect.DeepEqual(got, want) {
				t.Errorf("converter.ConvertCombinationMapsToTable() = %v, want %v", got, want)
			}
		})
	}
}
