package setup

import (
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checksource"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/config"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/esclient"
)

func Run() error {
	c := config.Get()

	err := Kibana(c.Setup.Kibana, c.Setup.Username, c.Setup.Password, c.VerifyCerts)
	if err != nil {
		return err
	}

	es, err := esclient.New(c.Elasticsearch, c.Setup.Username, c.Setup.Password, c.VerifyCerts)
	if err != nil {
		return err
	}

	err = Elasticsearch(es)
	if err != nil {
		return err
	}

	f := &checksource.Filesystem{
		Path:  c.Setup.CheckFolder,
		Teams: c.Teams,
	}

	return Checks(es, f)
}
