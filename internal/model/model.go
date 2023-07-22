package model

type Data struct {
	Factors     Factors     `json:"factors" jsonschema:"title=Factors"`
	Conditions  Conditions  `json:"conditions,omitempty" jsonschema:"title=Conditions,description=Condition variables,example=is-bear-active: 'place.ground && !season.winter'"`
	Constraints Constraints `json:"constraints,omitempty" jsonschema:"title=Constraints"`
}

type Factors map[string]Factor

type Factor struct {
	Name   string            `json:"name" jsonschema:"title=Name,description=Outputted factor name"`
	Levels map[string]string `json:"levels" jsonschema:"title=Levels,description=Possible values of the factor"`
}

type Conditions map[string]string

type Constraint struct {
	If     string `json:"if,omitempty" jsonschema:"oneof_required=if,title=If,description=The condition in then should be satisfied if this condition is satisfied,example=factor1.level1 && !factor2.level2"`
	OnlyIf string `yaml:"only_if" json:"only_if,omitempty" jsonschema:"oneof_required=only_if,title=Only if,description=The condition in then is available only if this condition is satisfied,example=factor1.level1 && !factor2.level2"`
	Then   string `json:"then" jsonschema:"title=Then,example=factor1.level1 && !factor2.level2"`
	Else   string `json:"else,omitempty" jsonschema:"title=Else,example=factor1.level1 && !factor2.level2"`
}

type Constraints []Constraint

type Combination map[string]string
