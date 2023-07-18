package filterer_test

import (
	"reflect"
	"testing"

	"github.com/midorimici/gentestcase/internal/filterer"
	"github.com/midorimici/gentestcase/internal/model"
)

type mockParser struct{}

func (p *mockParser) Parse(combination model.Combination, text string) (bool, error) {
	e1 := combination["e1"]
	e2 := combination["e2"]
	e3 := combination["e3"]
	if e2 == "e" && e1 == "a" {
		return false, nil
	}
	if e3 == "h" && e2 != "f" {
		return false, nil
	}
	return true, nil
}

func Test_filterer_Filter(t *testing.T) {
	elems := model.Factors{
		"e1": {Levels: map[string]model.Level{"a": {}, "b": {}}},
		"e2": {Levels: map[string]model.Level{"d": {}, "e": {If: "!e1.a"}, "f": {}}},
		"e3": {Levels: map[string]model.Level{"g": {}, "h": {If: "e2.f"}}},
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
				{"e1": "a", "e2": "f", "e3": "g"},
				{"e1": "a", "e2": "f", "e3": "h"},
				{"e1": "b", "e2": "d", "e3": "g"},
				{"e1": "b", "e2": "e", "e3": "g"},
				{"e1": "b", "e2": "f", "e3": "g"},
				{"e1": "b", "e2": "f", "e3": "h"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := filterer.New(elems, parser, combinations)
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
