package esclient

import (
	"fmt"
	"io"
	"sync"

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
	var wg sync.WaitGroup
	for _, doc := range docs {
		wg.Add(1)
		go func(doc struct {
			string
			io.Reader
			error
		}, wg *sync.WaitGroup) {
			defer wg.Done()
			if doc.error != nil {
				fmt.Printf("failed to index result for %s: %s\n", result.ID, doc.error)
				return
			}

			res, err := c.Index(doc.string, doc.Reader)
			if err != nil {
				fmt.Printf("failed to index result document for %s: %s\n", result.ID, err)
				return
			}
			if res.IsError() {
				// TODO: better error message here. res.String() is for testing or
				// debugging only
				fmt.Printf("failed to index result document in elasticsearch for %s: %s\n", result.ID, res.String())
				return
			}
			defer res.Body.Close()
		}(doc, &wg)
	}
	wg.Wait()
	return nil
}
