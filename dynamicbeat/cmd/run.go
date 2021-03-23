package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const runShort = `Interact with an instance of Scorestack to run checks and store results.`
const runLong = runShort + `

Dynamicbeat will pull check configurations from the configured Scorestack
instance, execute the checks at regular intervals, and store the results in
Scorestack. This process will be repeated until Dynamicbeat is terminated.`

var cfgFile string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: runShort,
	Long:  runLong,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	cobra.OnInitialize(initConfig)

	runCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to config file (default is $PWD/dynamicbeat.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find current directory
		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		// Search config in home directory with name "dynamicbeat" (without extension).
		viper.AddConfigPath(cwd)
		viper.SetConfigName("dynamicbeat")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
