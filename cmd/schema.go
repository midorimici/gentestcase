package cmd

import (
	"github.com/midorimici/gentestcase/internal/schema"
	"github.com/spf13/cobra"
)

// schemaCmd represents the run command
func schemaCmd() *cobra.Command {
	var outputFilename string

	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Generate JSON schema file to validate definition YAML files",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := schema.New(outputFilename)
			return s.Save()
		},
	}

	cmd.Flags().StringVarP(&outputFilename, "output", "o", "schema.json", "output schema JSON filename")

	return cmd
}
