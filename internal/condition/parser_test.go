package condition_test

import (
	"testing"

	"github.com/midorimici/gentestcase/internal/condition"
	"github.com/midorimici/gentestcase/internal/model"
)

func Test_parser_Parse(t *testing.T) {
	elems := model.Elements{
		"e1": {Options: map[string]model.Option{"a": {}, "b": {}}},
		"e2": {Options: map[string]model.Option{"d": {}, "e": {}, "f": {}}},
		"e3": {Options: map[string]model.Option{"g": {}, "h": {}}},
	}
	combination := model.Combination{"e1": "a", "e2": "d", "e3": "g"}

	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// Normal
		{
			name: "returns result as expected with matching single condition",
			args: args{text: "e1.a"},
			want: true,
		},
		{
			name: "returns result as expected with unmatching single condition",
			args: args{text: "e1.b"},
			want: false,
		},

		{
			name: "returns result as expected with matching single AND condition",
			args: args{text: "e1.a && e2.d"},
			want: true,
		},
		{
			name: "returns result as expected with unmatching single AND condition",
			args: args{text: "e1.a && e2.e"},
			want: false,
		},
		{
			name: "returns result as expected with matching multiple AND condition",
			args: args{text: "e1.a && e2.d && e3.g"},
			want: true,
		},
		{
			name: "returns result as expected with unmatching multiple AND condition",
			args: args{text: "e1.a && e2.d && e3.h"},
			want: false,
		},

		{
			name: "returns result as expected with matching single OR condition",
			args: args{text: "e1.b || e2.d"},
			want: true,
		},
		{
			name: "returns result as expected with unmatching single OR condition",
			args: args{text: "e1.b || e2.e"},
			want: false,
		},
		{
			name: "returns result as expected with matching multiple OR condition",
			args: args{text: "e1.b || e2.e || e3.g"},
			want: true,
		},
		{
			name: "returns result as expected with unmatching multiple OR condition",
			args: args{text: "e1.b || e2.e || e3.h"},
			want: false,
		},

		{
			name: "returns result as expected with matching single condition with negation",
			args: args{text: "!e1.b"},
			want: true,
		},
		{
			name: "returns result as expected with unmatching single condition with negation",
			args: args{text: "!e1.a"},
			want: false,
		},
		{
			name: "returns result as expected with matching single AND condition with negation",
			args: args{text: "!e1.b && !e2.e"},
			want: true,
		},
		{
			name: "returns result as expected with unmatching single AND condition with negation",
			args: args{text: "!e1.b && !e2.d"},
			want: false,
		},
		{
			name: "returns result as expected with matching single OR condition with negation",
			args: args{text: "!e1.b || !e2.d"},
			want: true,
		},
		{
			name: "returns result as expected with unmatching single OR condition with negation",
			args: args{text: "!e1.a || !e2.d"},
			want: false,
		},

		{
			name: "returns result as expected with grouped condition",
			args: args{text: "e1.a && (e2.e || e3.g)"},
			want: true,
		},
		{
			name: "returns result as expected with grouped condition",
			args: args{text: "e1.a && (e2.e || e3.h)"},
			want: false,
		},
		{
			name: "returns result as expected with grouped condition",
			args: args{text: "e1.a || e1.b"},
			want: true,
		},
		{
			name: "returns result as expected with grouped condition",
			args: args{text: "!((!e2.d && !e2.e) || e3.h)"},
			want: true,
		},
		{
			name: "returns result as expected with grouped condition",
			args: args{text: "(e1.a || e1.b) && !((!e2.d && !e2.e) || e3.h)"},
			want: true,
		},

		// Abnormal
		{
			name:    "returns error with element which does not exist",
			args:    args{text: "e.a"},
			wantErr: true,
		},
		{
			name:    "returns error with option which does not exist",
			args:    args{text: "e1.c"},
			wantErr: true,
		},
		{
			name:    "returns error with invalid expression",
			args:    args{text: "e1.a&&e2.c"},
			wantErr: true,
		},
		{
			name:    "returns error with invalid logical operator",
			args:    args{text: "e1.a & e2.c"},
			wantErr: true,
		},
		{
			name:    "returns error with unpaired parenthesis",
			args:    args{text: "e1.a && e2.c)"},
			wantErr: true,
		},
		{
			name:    "returns error with unpaired parenthesis",
			args:    args{text: "!((e1.a && e2.c)"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := condition.NewParser(elems)
			got, err := p.Parse(combination, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
