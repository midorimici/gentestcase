package main

import (
	"encoding/json"
	"flag"
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

	if err := export(s); err != nil {
		log.Fatal(err)
	}
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
