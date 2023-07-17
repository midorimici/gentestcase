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
	s := jsonschema.Reflect(&model.Cases{})
	d, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(schemaFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.Write(d)
}
