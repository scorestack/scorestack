package esclient

import (
	"fmt"
	"io"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
)

func Index(c *elasticsearch.Client, result check.Result) error {
	docs := make([]struct {
		string
		io.Reader
		error
	}, 0)

	// Create the documents
	index, reader, err := result.Admin()
	docs = append(docs, struct {
		string
		io.Reader
		error
	}{index, reader, err})
	index, reader, err = result.Team()
	docs = append(docs, struct {
		string
		io.Reader
		error
	}{index, reader, err})
	index, reader, err = result.Generic()
	docs = append(docs, struct {
		string
		io.Reader
		error
	}{index, reader, err})

	// Loop through the documents and index them
	for _, doc := range docs {
		if doc.error != nil {
			return fmt.Errorf("failed to index result for %s: %s", result.ID, doc.error)
		}

		res, err := c.Index(doc.string, doc.Reader)
		if err != nil {
			return fmt.Errorf("failed to index result document for %s: %s", result.ID, err)
		}
		if res.IsError() {
			return fmt.Errorf("failed to index result document for %s: %s", result.ID, res.Status())
		}
		defer res.Body.Close()
	}

	return nil
}
