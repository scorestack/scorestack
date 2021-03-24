package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const viewShort = "View your current configuration."

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: viewShort,
	Run: func(cmd *cobra.Command, args []string) {
		c := viper.AllSettings()
		bytes, err := yaml.Marshal(&c)
		cobra.CheckErr(err)
		fmt.Print(string(bytes))
	},
}

func init() {
	configCmd.AddCommand(viewCmd)
}
