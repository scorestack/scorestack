package esclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esutil"
	"go.uber.org/zap"
)

type Indexable interface {
	GetIndex() string
	GetDocumentId() string
	GetBody() map[string]interface{}
}

func indexableErr(prefix string, i Indexable, err error) string {
	return indexErr(prefix, i.GetDocumentId(), i.GetIndex(), err)
}

func indexErr(prefix string, documentId string, index string, err error) string {
	return fmt.Sprintf("%s document_id='%s' index='%s': %s", prefix, documentId, index, err)
}

func Queue(idxr esutil.BulkIndexer, i Indexable) {
	item := esutil.BulkIndexerItem{
		Index:      i.GetIndex(),
		Action:     "index",
		DocumentID: i.GetDocumentId(),
		OnFailure:  onQueueFailure,
	}

	document := i.GetBody()
	buf, err := json.Marshal(document)
	if err != nil {
		zap.S().Error(indexableErr("failed to get body", i, err))
	}

	item.Body = bytes.NewReader(buf)

	err = idxr.Add(context.Background(), item)
	if err != nil {
		zap.S().Error(indexableErr("failed to index document", i, err))
	}
}

func onQueueFailure(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
	if err != nil {
		zap.S().Error(indexErr("failed to index document", item.DocumentID, item.Index, err))
	} else {
		zap.S().Error(indexErr("failed to index document", item.DocumentID, item.Index, err))
	}
}
