package schema

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/midorimici/gentestcase/internal/model"
)

type SchemaSaver interface {
	Save() error
}

type schemaSaver struct {
	outputFilename string
}

func New(outputFilename string) SchemaSaver {
	return &schemaSaver{outputFilename: outputFilename}
}

func (s *schemaSaver) Save() error {
	flag.Parse()

	schema := jsonschema.Reflect(&model.Data{})
	restrictPropertyNames(schema)
	if err := addConditionPatterns(schema); err != nil {
		return err
	}

	if err := s.export(schema); err != nil {
		return err
	}

	return nil
}

func restrictPropertyNames(s *jsonschema.Schema) {
	restrictPropertyName(s.Definitions["Conditions"])
	restrictPropertyName(s.Definitions["Factors"])
	op, ok := s.Definitions["Factor"].Properties.Get("levels")
	if !ok {
		log.Fatal("levels property not found in Factor")
	}
	restrictPropertyName(op.(*jsonschema.Schema))
}

func restrictPropertyName(s *jsonschema.Schema) {
	s.PatternProperties = map[string]*jsonschema.Schema{`^\w+$`: s.PatternProperties[".*"]}
	s.AdditionalProperties = jsonschema.FalseSchema
}

func addConditionPatterns(s *jsonschema.Schema) error {
	const funcName = "addConditionPattern"
	const condRe = `^[$!(\w][$!().\w&| ]+[\w)]$`

	s.Definitions["Conditions"].PatternProperties[`^\w+$`].Pattern = condRe

	condProps := []string{"if", "only_if", "then", "else"}
	constraintProps := s.Definitions["Constraint"].Properties
	for _, prop := range condProps {
		p, ok := constraintProps.Get(prop)
		if !ok {
			return fmt.Errorf("%s: if property not found in Level", funcName)
		}
		ifSchema, ok := p.(*jsonschema.Schema)
		if !ok {
			return fmt.Errorf("%s: if property cannot convert to schema", funcName)
		}
		ifSchema.Pattern = condRe
	}

	return nil
}

func (s *schemaSaver) export(schema *jsonschema.Schema) error {
	d, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Setup output writer
	var out io.Writer
	if s.outputFilename == "" {
		out = os.Stdout
	} else {
		f, err := os.Create(s.outputFilename)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}

	out.Write(d)

	return nil
}
