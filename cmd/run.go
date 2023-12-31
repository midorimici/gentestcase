package cmd

import (
	"github.com/midorimici/gentestcase/internal/runner"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
func runCmd() *cobra.Command {
	var (
		inputFilename  string
		outputFilename string
		isWatching     bool
		debugFilename  string
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Generate combination test cases",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := runner.New(inputFilename, outputFilename, isWatching, debugFilename)
			return r.Run()
		},
	}

	cmd.Flags().StringVarP(&inputFilename, "input", "i", "cases.yml", "input YAML filename")
	cmd.Flags().StringVarP(&outputFilename, "output", "o", "data.csv", "output CSV filename")
	cmd.Flags().BoolVarP(&isWatching, "watch", "w", false, "watch input file change")
	cmd.Flags().StringVarP(&debugFilename, "debug", "d", "", "output debug CSV with debug information to the specified filename")

	return cmd
}
