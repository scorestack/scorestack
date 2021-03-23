package cmd

import (
	"github.com/spf13/cobra"
)

const rootShort = "A service health check utility."

var rootLong = rootShort + `

Dynamicbeat interacts with network services like file shares and webservers to
determine if they are up and running properly. Dynamicbeat is a component of
the Scorestack project.`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dynamicbeat [command]",
	Short: rootShort,
	Long:  rootLong,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
