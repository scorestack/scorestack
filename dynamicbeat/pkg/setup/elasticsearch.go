package setup

import (
	"strings"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets/indices"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets/users"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/esclient"
	"go.uber.org/zap"
)

func Elasticsearch(c *esclient.Client) error {
	zap.S().Info("checking if Elasticsearch is up")
	err := c.Wait()
	if err != nil {
		return err
	}

	err = c.AddUser("dynamicbeat", users.Dynamicbeat())
	if err != nil {
		return err
	}

	// Add default index template
	zap.S().Info("adding default index template")
	idx := strings.NewReader(`{"index_patterns":["check*","attrib_*","results*"],"settings":{"number_of_replicas":"0"}}`)
	res, err := c.Indices.PutTemplate("default", idx)
	if err != nil {
		return err
	}
	err = c.CloseAndCheck(res)
	if err != nil {
		return err
	}

	// Create results indices
	err = c.AddIndex("results-admin", indices.ResultsAdmin())
	if err != nil {
		return err
	}
	err = c.AddIndex("results-all", indices.ResultsAll())
	if err != nil {
		return err
	}

	return nil
}
