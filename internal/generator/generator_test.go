package generator_test

import (
	"reflect"
	"testing"

	"integtest/internal/generator"
	"integtest/internal/model"
)

func Test_generator_Generate(t *testing.T) {
	c := model.Cases{
		"e1": {Options: map[string]model.Option{"a": {}, "b": {}}},
		"e2": {Options: map[string]model.Option{"d": {}, "e": {}, "f": {}}},
		"e3": {Options: map[string]model.Option{"g": {}, "h": {}}},
	}

	type args struct {
		cases model.Cases
	}
	tests := []struct {
		name string
		args args
		want []model.Combination
	}{
		{
			name: "returns result as expected",
			args: args{c},
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
			g := generator.New(tt.args.cases)
			if got := g.Generate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generator.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
