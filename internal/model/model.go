package model

type Cases map[string]Element

type Element struct {
	Name    string            `json:"name" jsonschema:"required"`
	Options map[string]Option `json:"options" jsonschema:"required"`
}

type Option struct {
	Name string `json:"name" jsonschema:"required"`

	// If represents condition to output the option.
	// example: element1.option_id && !element2.option_id
	If string `json:"if" jsonschema:"title=Condition,description=Condition to output the option,example=element1.option_id && !element2.option_id"`
}

type Combination map[string]string
