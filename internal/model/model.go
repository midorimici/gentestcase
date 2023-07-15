package model

type Cases map[string]Element

type Element struct {
	Order int
	Name string
	Options map[string]Option
}

type Option struct {
	Name string
	
	// If represents condition to output the option.
	// example: element1.option_id && !element2.option_id
	If string
}

type Combination map[string]string
