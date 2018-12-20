package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var kongfig = &cobra.Command{
	Use:   "kongfig",
	Short: "Kongfig is a configuration management tool for Kong API gateway",
	Long: `Kongfig is a configuration management tool for the Kong API gateway.

Find more information at https://github.com/pagerinc/kongfig`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute runs the konfig cli
func Execute() {
	if err := kongfig.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
