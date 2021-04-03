// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	RoundTime     time.Duration `mapstructure:"round_time"`
	Elasticsearch string        `mapstructure:"elasticsearch"`
	Username      string        `mapstructure:"username"`
	Password      string        `mapstructure:"password"`
	VerifyCerts   bool          `mapstructure:"verify_certs"`
	Teams         []Team        `mapstructure:"teams"`
	Setup         struct {
		Kibana   string `mapstructure:"kibana"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"setup"`
	Log struct {
		Verbose bool `mapstructure:"verbose"`
		Level   int8 `mapstructure:"level"`
		NoColor bool `mapstructure:"no_color"`
	} `mapstructure:"log"`
}

type Team struct {
	Name      string            `mapstructure:"name"`
	Overrides map[string]string `mapstructure:"overrides"`
}

func Get() Config {
	var c Config
	err := viper.Unmarshal(&c)
	cobra.CheckErr(err)
	return c
}
