package setup

import (
	"fmt"

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
	defs, err := f.LoadAll()
	if err != nil {
		return err
	}

	for _, def := range defs {
		fmt.Printf("%s: %#v", def.ID, def)
	}

	return nil
}
