package cmd

import (
	"fmt"
	"os"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/checksource"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/config"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/esclient"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/setup"
	"github.com/spf13/cobra"
)

const checksShort = "Add or update checks."

var team string

// checksCmd represents the checks command
var checksCmd = &cobra.Command{
	Use:   "checks [path to checks]",
	Short: checksShort,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Get()

		// Filter out teams if specified
		var teams []config.Team
		if team != "" {
			for _, t := range c.Teams {
				if t.Name == team {
					teams = append(teams, t)
				}
			}
		} else {
			teams = c.Teams
		}

		// Make sure at least one team exists
		if len(teams) == 0 {
			fmt.Printf("No teams found. If you passed -t/--team, make sure the team you specified actually exists.")
			os.Exit(1)
		}

		es, err := esclient.New(c.Elasticsearch, c.Setup.Username, c.Setup.Password, c.VerifyCerts)
		cobra.CheckErr(err)

		f := &checksource.Filesystem{
			Path:  args[0],
			Teams: teams,
		}

		cobra.CheckErr(setup.Checks(es, f))
	},
}

func init() {
	setupCmd.AddCommand(checksCmd)

	checksCmd.Flags().StringVarP(&team, "team", "t", "", "add checks for only the specified team")
}
