package model

type Data struct {
	Factors    Factors    `json:"factors" jsonschema:"title=Factors"`
	Conditions Conditions `json:"conditions,omitempty" jsonschema:"title=Conditions,description=Condition variables,example=is-bear-active: 'place.ground && !season.winter'"`
}

type Factors map[string]Factor

type Factor struct {
	Name   string           `json:"name" jsonschema:"title=Name,description=Outputted factor name"`
	Levels map[string]Level `json:"levels" jsonschema:"title=Levels"`
}

type Level struct {
	Name string `json:"name" jsonschema:"title=Name,description=Outputted level name"`

	// If represents condition to output the level.
	// example: factor1.option_id && !factor2.option_id
	If string `json:"if,omitempty" jsonschema:"title=Condition,description=Condition to output the level,example=factor1.option_id && !factor2.option_id"`
}

type Conditions map[string]string

type Combination map[string]string
