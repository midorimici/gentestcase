package filterer_test

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/midorimici/gentestcase/internal/filterer"
	"github.com/midorimici/gentestcase/internal/model"
)

type mockParser struct{}

func (p *mockParser) Parse(combination model.Combination, text string) (bool, error) {
	re := regexp.MustCompile(`^(!)?(\w+)\.(\w+)$`)
	match := re.FindStringSubmatch(text)
	isNeg := match[1] == "!"
	factor := match[2]
	level := match[3]

	result := combination[factor] == level
	if isNeg {
		result = !result
	}

	return result, nil
}

func Test_filterer_Filter(t *testing.T) {
	constraints := model.Constraints{
		{
			OnlyIf: "!e1.a",
			Then:   "e2.e",
		},
		{
			OnlyIf: "e2.f",
			Then:   "e3.h",
			Else:   "e3.g",
		},
	}
	parser := &mockParser{}
	combinations := []model.Combination{
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

	tests := []struct {
		name    string
		want    []model.Combination
		wantErr bool
	}{
		{
			name: "returns result as expected",
			want: []model.Combination{
				{"e1": "a", "e2": "d", "e3": "g"},
				{"e1": "a", "e2": "f", "e3": "h"},
				{"e1": "b", "e2": "d", "e3": "g"},
				{"e1": "b", "e2": "e", "e3": "g"},
				{"e1": "b", "e2": "f", "e3": "h"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := filterer.New(constraints, parser, combinations)
			got, err := f.Filter()
			if (err != nil) != tt.wantErr {
				t.Errorf("filterer.Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterer.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
