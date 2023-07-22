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
	c1 := model.Constraint{
		If:   "e1.a",
		Then: "e2.d",
	}
	c2 := model.Constraint{
		If:   "e2.e",
		Then: "e3.g",
		Else: "e3.h",
	}
	c3 := model.Constraint{
		OnlyIf: "!e1.a",
		Then:   "e2.e",
	}
	c4 := model.Constraint{
		OnlyIf: "e2.f",
		Then:   "e3.h",
		Else:   "e3.g",
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
		name        string
		constraints model.Constraints
		want        []model.Combination
		wantErr     bool
	}{
		// Normal
		{
			name:        "returns result as expected with single condition",
			constraints: model.Constraints{c1},
			want: []model.Combination{
				{"e1": "a", "e2": "d", "e3": "g"},
				{"e1": "a", "e2": "d", "e3": "h"},
				{"e1": "b", "e2": "d", "e3": "g"},
				{"e1": "b", "e2": "d", "e3": "h"},
				{"e1": "b", "e2": "e", "e3": "g"},
				{"e1": "b", "e2": "e", "e3": "h"},
				{"e1": "b", "e2": "f", "e3": "g"},
				{"e1": "b", "e2": "f", "e3": "h"},
			},
		},
		{
			name:        "returns result as expected with single condition",
			constraints: model.Constraints{c2},
			want: []model.Combination{
				{"e1": "a", "e2": "d", "e3": "h"},
				{"e1": "a", "e2": "e", "e3": "g"},
				{"e1": "a", "e2": "f", "e3": "h"},
				{"e1": "b", "e2": "d", "e3": "h"},
				{"e1": "b", "e2": "e", "e3": "g"},
				{"e1": "b", "e2": "f", "e3": "h"},
			},
		},
		{
			name:        "returns result as expected with single condition",
			constraints: model.Constraints{c3},
			want: []model.Combination{
				{"e1": "a", "e2": "d", "e3": "g"},
				{"e1": "a", "e2": "d", "e3": "h"},
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
		{
			name:        "returns result as expected with single condition",
			constraints: model.Constraints{c4},
			want: []model.Combination{
				{"e1": "a", "e2": "d", "e3": "g"},
				{"e1": "a", "e2": "e", "e3": "g"},
				{"e1": "a", "e2": "f", "e3": "h"},
				{"e1": "b", "e2": "d", "e3": "g"},
				{"e1": "b", "e2": "e", "e3": "g"},
				{"e1": "b", "e2": "f", "e3": "h"},
			},
		},
		{
			name:        "returns result as expected with multiple conditions",
			constraints: model.Constraints{c1, c2},
			want: []model.Combination{
				{"e1": "a", "e2": "d", "e3": "h"},
				{"e1": "b", "e2": "d", "e3": "h"},
				{"e1": "b", "e2": "e", "e3": "g"},
				{"e1": "b", "e2": "f", "e3": "h"},
			},
		},
		{
			name:        "returns result as expected with multiple conditions",
			constraints: model.Constraints{c1, c3},
			want: []model.Combination{
				{"e1": "a", "e2": "d", "e3": "g"},
				{"e1": "a", "e2": "d", "e3": "h"},
				{"e1": "b", "e2": "d", "e3": "g"},
				{"e1": "b", "e2": "d", "e3": "h"},
				{"e1": "b", "e2": "e", "e3": "g"},
				{"e1": "b", "e2": "e", "e3": "h"},
				{"e1": "b", "e2": "f", "e3": "g"},
				{"e1": "b", "e2": "f", "e3": "h"},
			},
		},
		{
			name:        "returns result as expected with multiple conditions",
			constraints: model.Constraints{c3, c4},
			want: []model.Combination{
				{"e1": "a", "e2": "d", "e3": "g"},
				{"e1": "a", "e2": "f", "e3": "h"},
				{"e1": "b", "e2": "d", "e3": "g"},
				{"e1": "b", "e2": "e", "e3": "g"},
				{"e1": "b", "e2": "f", "e3": "h"},
			},
		},
		{
			name:        "returns result as expected with multiple conditions",
			constraints: model.Constraints{c1, c2, c3, c4},
			want: []model.Combination{
				{"e1": "b", "e2": "e", "e3": "g"},
				{"e1": "b", "e2": "f", "e3": "h"},
			},
		},

		// Abnormal
		{
			name: "returns error when neither if nor only_if is specified",
			constraints: model.Constraints{
				{
					Then: "e3.h",
				},
			},
			wantErr: true,
		},
		{
			name: "returns error when both if and only_if are specified",
			constraints: model.Constraints{
				{
					If:     "e1.a",
					OnlyIf: "e1.a",
					Then:   "e3.h",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := filterer.New(tt.constraints, parser, combinations)
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
