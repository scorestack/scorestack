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
		chk, generic, admin, user, err := def.Documents()
		if err != nil {
			zap.S().Errorf("skipping check due to error - %s", err)
		}

		queueItem(indexer, "checkdef", def.ID, chk)
		queueItem(indexer, "checks", def.ID, generic)
		if admin != nil {
			queueItem(indexer, fmt.Sprintf("attrib_admin_%s", def.Group), def.ID, admin)
		}
		if user != nil {
			queueItem(indexer, fmt.Sprintf("attrib_user_%s", def.Group), def.ID, user)
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
