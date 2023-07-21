package condition_test

import (
	"testing"

	"github.com/midorimici/gentestcase/internal/condition"
	"github.com/midorimici/gentestcase/internal/model"
)

func Test_parser_Parse(t *testing.T) {
	elems := model.Factors{
		"e1": {Levels: map[string]string{"a": "a", "b": "b"}},
		"e2": {Levels: map[string]string{"d": "d", "e": "e", "f": "f"}},
		"e3": {Levels: map[string]string{"g": "g", "h": "h"}},
	}
	refs := model.Conditions{
		"c1": "e1.a",
		"c2": "!e2.e",
		"c3": "e3.h",
		"c4": "$c1 && $c2",
	}
	data := &model.Data{
		Factors:    elems,
		Conditions: refs,
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
		{
			name: "returns result as expected with matching single condition reference",
			args: args{text: "$c1"},
			want: true,
		},
		{
			name: "returns result as expected with matching single condition reference with negation",
			args: args{text: "!$c2"},
			want: false,
		},
		{
			name: "returns result as expected with unmatching single condition reference",
			args: args{text: "$c3"},
			want: false,
		},
		{
			name: "returns result as expected with matching single condition reference which refers to other condition references",
			args: args{text: "$c4"},
			want: true,
		},
		{
			name: "returns result as expected with multiple condition reference",
			args: args{text: "$c1 && $c2"},
			want: true,
		},
		{
			name: "returns result as expected with condition reference and normal condition",
			args: args{text: "e2.e || $c1"},
			want: true,
		},

		// Abnormal
		{
			name:    "returns error with factor which does not exist",
			args:    args{text: "e.a"},
			wantErr: true,
		},
		{
			name:    "returns error with level which does not exist",
			args:    args{text: "e1.c"},
			wantErr: true,
		},
		{
			name:    "returns error with condition reference which does not exist",
			args:    args{text: "$ref"},
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
			p := condition.NewParser(data)
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
