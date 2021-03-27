package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const rootShort = "A service health check utility."
const rootLong = rootShort + `

Dynamicbeat interacts with network services like file shares and webservers to
determine if they are up and running properly. Dynamicbeat is a component of
the Scorestack project.`

var cfgFile string

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

func init() {
	cobra.OnInitialize(initConfig)

	// Config file path
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to config file (default: ${PWD}/dynamicbeat.yaml)")

	// Config file contents
	addFlag("round_time", "r", "30s", "time to wait between rounds of checks")
	addFlag("elasticsearch", "e", "https://localhost:9200", "address of Elasticsearch host to pull checks from and store results in")
	addFlag("username", "u", "dynamicbeat", "username for authentication with Elasticsearch")
	addFlag("password", "p", "changeme", "password for authentication with Elasticsearch")
	addBoolFlag("verify_certs", "v", false, "whether to verify the Elasticsearch TLS certificates")

	// Configure five default teams
	teams := make([]config.Team, 5)
	for i := 0; i < len(teams); i++ {
		teams[i] = config.Team{Name: fmt.Sprintf("team%02d", i+1)}
	}
	viper.SetDefault("teams", teams)
}

func addFlag(name string, short string, value string, help string) {
	rootCmd.PersistentFlags().StringP(name, short, value, help)
	_ = viper.BindPFlag(name, rootCmd.PersistentFlags().Lookup(name))
}

func addBoolFlag(name string, short string, value bool, help string) {
	rootCmd.PersistentFlags().BoolP(name, short, value, help)
	_ = viper.BindPFlag(name, rootCmd.PersistentFlags().Lookup(name))
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

	// Make sure dot separators are replaced by underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read in any matching environment variables
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
