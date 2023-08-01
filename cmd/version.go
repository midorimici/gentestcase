package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "1.4.2"

// versionCmd represents the version command
func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show gentestcase version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
	return cmd
}
