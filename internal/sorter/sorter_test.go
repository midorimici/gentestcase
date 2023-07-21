package sorter_test

import (
	"reflect"
	"testing"

	"github.com/midorimici/gentestcase/internal/model"
	"github.com/midorimici/gentestcase/internal/sorter"
)

func Test_sorter_Sort(t *testing.T) {
	elems := model.Factors{
		"e1": {Levels: map[string]string{"a": "a", "b": "b"}},
		"e2": {Levels: map[string]string{"d": "d", "e": "e", "f": "f"}},
		"e3": {Levels: map[string]string{"g": "g", "h": "h"}},
	}
	orderedFactors := []string{"e1", "e2", "e3"}
	levelOrders := map[string]map[string]int{
		"e1": {"a": 0, "b": 1},
		"e2": {"d": 2, "e": 3, "f": 4},
		"e3": {"g": 5, "h": 6},
	}
	want := []model.Combination{
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
	}

	type args struct {
		c []model.Combination
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
				{"e1": "b", "e2": "f", "e3": "h"},
				{"e1": "a", "e2": "d", "e3": "h"},
				{"e1": "a", "e2": "f", "e3": "g"},
				{"e1": "b", "e2": "e", "e3": "h"},
				{"e1": "b", "e2": "f", "e3": "g"},
				{"e1": "a", "e2": "d", "e3": "g"},
				{"e1": "b", "e2": "d", "e3": "h"},
				{"e1": "a", "e2": "e", "e3": "g"},
				{"e1": "a", "e2": "e", "e3": "h"},
				{"e1": "b", "e2": "d", "e3": "g"},
				{"e1": "a", "e2": "f", "e3": "h"},
				{"e1": "b", "e2": "e", "e3": "g"},
			}},
		},
		{
			name: "returns result as expected",
			args: args{[]model.Combination{
				{"e3": "h", "e2": "f", "e1": "b"},
				{"e3": "h", "e2": "d", "e1": "a"},
				{"e3": "g", "e2": "f", "e1": "a"},
				{"e3": "h", "e2": "e", "e1": "b"},
				{"e3": "g", "e2": "f", "e1": "b"},
				{"e3": "g", "e2": "d", "e1": "a"},
				{"e3": "h", "e2": "d", "e1": "b"},
				{"e3": "g", "e2": "e", "e1": "a"},
				{"e3": "h", "e2": "e", "e1": "a"},
				{"e3": "g", "e2": "d", "e1": "b"},
				{"e3": "h", "e2": "f", "e1": "a"},
				{"e3": "g", "e2": "e", "e1": "b"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sorter.New(elems, orderedFactors, levelOrders)
			if got := s.Sort(tt.args.c); !reflect.DeepEqual(got, want) {
				t.Errorf("sorter.Sort() = %v, want %v", got, want)
			}
		})
	}
}
