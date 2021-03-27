package cmd

import (
	"github.com/scorestack/scorestack/dynamicbeat/pkg/config"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/esclient"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/setup"
	"github.com/spf13/cobra"
)

const elasticsearchShort = "Add users and indices to Elasticsearch."

// elasticsearchCmd represents the elasticsearch command
var elasticsearchCmd = &cobra.Command{
	Use:   "elasticsearch",
	Short: elasticsearchShort,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Get()

		es, err := esclient.New(c.Elasticsearch, c.Setup.Username, c.Setup.Password, c.VerifyCerts)
		cobra.CheckErr(err)
		cobra.CheckErr(setup.Elasticsearch(es, c.Teams))
	},
}

func init() {
	setupCmd.AddCommand(elasticsearchCmd)
}
