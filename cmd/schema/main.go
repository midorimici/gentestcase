package main

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

var outputFilename = flag.String("o", "schema.json", "output schema JSON filename")

func main() {
	flag.Parse()

	s := jsonschema.Reflect(&model.Data{})
	restrictPropertyNames(s)
	if err := addConditionPattern(s); err != nil {
		log.Fatal(err)
	}

	if err := export(s); err != nil {
		log.Fatal(err)
	}
}

func restrictPropertyNames(s *jsonschema.Schema) {
	restrictPropertyName(s.Definitions["Conditions"])
	restrictPropertyName(s.Definitions["Elements"])
	op, ok := s.Definitions["Element"].Properties.Get("options")
	if !ok {
		log.Fatal("options property not found in Element")
	}
	restrictPropertyName(op.(*jsonschema.Schema))
}

func restrictPropertyName(s *jsonschema.Schema) {
	s.PatternProperties = map[string]*jsonschema.Schema{`^\w+$`: s.PatternProperties[".*"]}
	s.AdditionalProperties = jsonschema.FalseSchema
}

func addConditionPattern(s *jsonschema.Schema) error {
	const funcName = "addConditionPattern"
	const condRe = `^[$!(\w][$!().\w&| ]+[\w)]$`

	s.Definitions["Conditions"].PatternProperties[`^\w+$`].Pattern = condRe
	p, ok := s.Definitions["Option"].Properties.Get("if")
	if !ok {
		return fmt.Errorf("%s: if property not found in Option", funcName)
	}
	ifSchema, ok := p.(*jsonschema.Schema)
	if !ok {
		return fmt.Errorf("%s: if property cannot convert to schema", funcName)
	}
	ifSchema.Pattern = condRe

	return nil
}

func export(s *jsonschema.Schema) error {
	d, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Setup output writer
	var out io.Writer
	if *outputFilename == "" {
		out = os.Stdout
	} else {
		f, err := os.Create(*outputFilename)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}

	out.Write(d)

	return nil
}
