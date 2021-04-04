package setup

import (
	"github.com/scorestack/scorestack/dynamicbeat/pkg/config"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/esclient"
)

func Run() error {
	c := config.Get()

	err := Kibana(c.Setup.Kibana, c.Setup.Username, c.Setup.Password, c.VerifyCerts, c.Teams)
	if err != nil {
		return err
	}

	es, err := esclient.New(c.Elasticsearch, c.Setup.Username, c.Setup.Password, c.VerifyCerts)
	if err != nil {
		return err
	}

	return Elasticsearch(es, c.Teams)
}
