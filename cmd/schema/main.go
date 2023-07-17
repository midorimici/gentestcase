package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/midorimici/gentestcase/internal/model"
)

const schemaFilename = "schema.json"

func main() {
	s := jsonschema.Reflect(&model.Data{})
	export(s)
}

func export(s *jsonschema.Schema) {
	d, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(schemaFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.Write(d)
}
