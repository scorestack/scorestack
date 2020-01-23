// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period       time.Duration     `config:"period"`
	UpdatePeriod time.Duration     `config:"update_period"`
	CheckSource  CheckSourceConfig `config:"check_source"`
}

type CheckSourceConfig struct {
	Hosts       []string `config:"hosts"`
	Username    string   `config:"username"`
	Password    string   `config:"password"`
	VerifyCerts bool     `config:"verify_certs"`
}

var DefaultConfig = Config{
	Period:       1 * time.Second,
	UpdatePeriod: 1 * time.Minute,
	CheckSource: CheckSourceConfig{
		Hosts:       []string{"http://localhost:9200"},
		Username:    "dynamicbeat",
		Password:    "changeme",
		VerifyCerts: true,
	},
}
