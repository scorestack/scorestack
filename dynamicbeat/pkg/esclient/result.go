package esclient

import (
	"fmt"
	"io"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
)

func (c *Client) AddResult(result check.Result) error {
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
			return fmt.Errorf("failed to index result for %s: %s", result.CheckId, doc.error)
		}

		res, err := c.Index(doc.string, doc.Reader)
		if err != nil {
			return fmt.Errorf("failed to index result document for %s: %s", result.CheckId, err)
		}
		if res.IsError() {
			// TODO: better error message here. res.String() is for testing or
			// debugging only
			return fmt.Errorf("failed to index result document for %s: %s", result.CheckId, res.String())
		}
		defer res.Body.Close()
	}

	return nil
}
