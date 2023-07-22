package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "gentestcase",
		Short: "CLI tool to generate combination test cases",
	}
	setSubcommands(rootCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var subcommands = []func() *cobra.Command{
	versionCmd,
	runCmd,
}

func setSubcommands(rootCmd *cobra.Command) {
	for _, c := range subcommands {
		rootCmd.AddCommand(c())
	}
}
