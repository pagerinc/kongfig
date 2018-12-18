package cmd

import (
	"github.com/pagerinc/kongfig/api"
	"github.com/spf13/cobra"
)

var (
	fileVar   string
	dryRunVar bool
)

func init() {
	const (
		defaultConfig = "config.json"
		configUsage   = "Filename that contains the configuration to apply"
		defaultDryRun = false
		dryRunUsage   = "simulate an install"
	)

	applyCmd.Flags().StringVarP(&fileVar, "file", "f", defaultConfig, configUsage)
	applyCmd.Flags().BoolVar(&dryRunVar, "dry-run", defaultDryRun, dryRunUsage)
	kongfig.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration to a Kong instance",
	Long:  `Use apply to restore your settings into an existing Kong instance.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient(fileVar)
		if err != nil {
			return err
		}
		err = client.UpdateAllRecursively()

		return err
	},
}
