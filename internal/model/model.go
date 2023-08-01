package model

// Data contains factors, conditions and constraints.
type Data struct {
	Factors     Factors     `json:"factors" jsonschema:"title=Factors"`
	Conditions  Conditions  `json:"conditions,omitempty" jsonschema:"title=Conditions,description=Condition variables,example=is-bear-active: 'place.ground && !season.winter'"`
	Constraints Constraints `json:"constraints,omitempty" jsonschema:"title=Constraints"`
}

// Factors is a map with factor IDs as keys and Factor structs as values.
type Factors map[string]Factor

// Factor represents a factor.
type Factor struct {
	// Name is a name of the factor.
	// It is outputted to a CSV.
	Name string `json:"name" jsonschema:"title=Name,description=Outputted factor name"`

	// Levels is a map with level IDs as keys and level names as values.
	Levels map[string]string `json:"levels" jsonschema:"title=Levels,description=Possible values of the factor"`
}

// Conditions is a map with condition IDs as keys and condition statements as values.
type Conditions map[string]string

// Constraint contains constraint data.
type Constraint struct {
	// ID is a constraint ID outputted to a CSV with debugging.
	ID     string `json:"id,omitempty" jsonschema:"title=ID,description=Constraint ID used for debugging"`
	If     string `json:"if,omitempty" jsonschema:"oneof_required=if,title=If,description=The condition in then should be satisfied if this condition is satisfied,example=factor1.level1 && !factor2.level2"`
	OnlyIf string `yaml:"only_if" json:"only_if,omitempty" jsonschema:"oneof_required=only_if,title=Only if,description=The condition in then is available only if this condition is satisfied,example=factor1.level1 && !factor2.level2"`
	Then   string `json:"then" jsonschema:"title=Then,example=factor1.level1 && !factor2.level2"`
	Else   string `json:"else,omitempty" jsonschema:"title=Else,example=factor1.level1 && !factor2.level2"`
}

// Constraints is a list of Constraint
type Constraints []Constraint

// Combination represents pairs of factors and levels.
//
// It is a map with factor IDs as keys and level IDs as values.
type Combination map[string]string
