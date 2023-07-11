package model

type Cases map[string]Element

type Element struct {
	Name string
	Options map[string]Option
}

type Option struct {
	Name string
	Only []string
	Except []string
}
