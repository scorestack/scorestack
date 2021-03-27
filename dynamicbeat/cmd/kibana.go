package cmd

import (
	"github.com/scorestack/scorestack/dynamicbeat/pkg/config"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/setup"
	"github.com/spf13/cobra"
)

const kibanaShort = "Add dashboards to Kibana and configure Elastic Stack roles."

// kibanaCmd represents the kibana command
var kibanaCmd = &cobra.Command{
	Use:   "kibana",
	Short: kibanaShort,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Get()

		cobra.CheckErr(setup.Kibana(c.Setup.Kibana, c.Setup.Username, c.Setup.Password, c.VerifyCerts, c.Teams))
	},
}

func init() {
	setupCmd.AddCommand(kibanaCmd)
}
