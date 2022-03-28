package setup

import (
	"context"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checksource"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/esclient"
	"go.uber.org/zap"
)

func Checks(c *esclient.Client, f *checksource.Filesystem) error {
	zap.S().Infof("loading checks from %s", f.Path)
	defs, err := f.LoadAll()
	if err != nil {
		return err
	}

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

func queueItem(i esutil.BulkIndexer, index string, id string, body io.Reader) {
	err := i.Add(
		context.Background(),
		esutil.BulkIndexerItem{
			Index:      index,
			Action:     "index",
			DocumentID: id,
			Body:       body,
			OnFailure: func(
				ctx context.Context,
				item esutil.BulkIndexerItem,
				res esutil.BulkIndexerResponseItem,
				err error,
			) {
				if err != nil {
					zap.S().Errorf("failed to add document of id '%s' to index '%s': %s", id, index, err)
				} else {
					zap.S().Errorf("failed to add document of id '%s' to index '%s' due to %s error: %s", id, index, res.Error.Type, res.Error.Reason)
				}
			},
		},
	)
	if err != nil {
		zap.S().Errorf("failed to add document of id '%s' and index '%s' to bulk index queue: %s", err)
	}
}
