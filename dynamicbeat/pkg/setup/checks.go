package setup

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checksource"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/esclient"
	"go.uber.org/zap"
)

func Checks(c *esclient.Client, f *checksource.Filesystem) error {
	defs, err := f.LoadAll()
	if err != nil {
		return err
	}
	zap.S().Infof("loaded %d check(s)", len(defs))

	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: c.Client,
	})
	if err != nil {
		return fmt.Errorf("failed to build bulk indexer: %s", err)
	}

	for _, def := range defs {
		esclient.Queue(indexer, &def)

		for _, attribute := range def.Attributes {
			esclient.Queue(indexer, &attribute)

			for _, value := range attribute.Values {
				esclient.Queue(indexer, &value)
			}
		}
	}

	zap.S().Info("waiting for checks to finish indexing...")
	return indexer.Close(context.Background())
}
