package model

type Data struct {
	Elements   Elements   `json:"elements" jsonschema:"title=Elements"`
	Conditions Conditions `json:"conditions,omitempty" jsonschema:"title=Conditions,description=Condition variables,example=is-bear-active: 'place.ground && !season.winter'"`
}

type Elements map[string]Element

type Element struct {
	Name    string            `json:"name" jsonschema:"title=Name,description=Outputted element name"`
	Options map[string]Option `json:"options" jsonschema:"title=Options"`
}

type Option struct {
	Name string `json:"name" jsonschema:"title=Name,description=Outputted option name"`

	// If represents condition to output the option.
	// example: element1.option_id && !element2.option_id
	If string `json:"if,omitempty" jsonschema:"title=Condition,description=Condition to output the option,example=element1.option_id && !element2.option_id"`
}

type Conditions map[string]string

type Combination map[string]string
