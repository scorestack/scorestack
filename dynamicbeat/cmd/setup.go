package cmd

import (
	"fmt"
	"strings"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/setup"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const setupShort = "Add dashboards, users, roles, and indices to a Scorestack instance."
const setupLong = setupShort + `

Adds the Scoreboard dashboard, users and roles, and indexes to Scorestack, as
well as configuring some Kibana settings. Additionally, the Team Overview
dashboard, team user, and team results index will be added for each configured
team.`

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: setupShort,
	Long:  setupLong,
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(setup.Run())
	},
}

func setupFlag(name string, short string, value string, help string) {
	setupCmd.PersistentFlags().StringP(name, short, value, help)
	name = strings.TrimPrefix(name, "setup-") // Remove the setup- part of the setup username and password
	_ = viper.BindPFlag(fmt.Sprintf("setup.%s", name), setupCmd.Flags().Lookup(name))
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Config file contents
	setupFlag("kibana", "k", "https://localhost:5601", "address of Kibana host to set up")
	setupFlag("setup-username", "U", "elastic", "username of Elasticsearch superuser to use for setup")
	setupFlag("setup-password", "P", "changeme", "password of Elasticsearch superuser to use for setup")
}
