package cmd

import (
	"github.com/spf13/cobra"
)

const configShort = "View or save your current configuration."

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: configShort,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
