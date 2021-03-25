package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const setupShort = "Add dashboards, users, and checks to a Scorestack instance."
const setupLong = setupShort + `

Adds the Scoreboard dashboard, users and roles, and indexes to Scorestack, as
well as configuring some Kibana settings. Additionally, the Team Overview
dashboard, team user, team results index, and team check definitions and
attributes will be added for each configured team.`

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: setupShort,
	Long:  setupLong,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setup called")
	},
}

func setupFlag(name string, short string, value string, help string) {
	setupCmd.Flags().StringP(name, short, value, help)
	_ = viper.BindPFlag(fmt.Sprintf("setup.%s", name), setupCmd.Flags().Lookup(name))
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Config file contents
	setupFlag("kibana", "k", "https://localhost:5601", "address of Kibana host to set up")
	setupFlag("check_folder", "f", "./checks", "path to the folder that contains check files")
}
