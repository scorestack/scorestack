package esclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch"
)

// The SearchResults struct is used to parse search results from Elasticsearch
// and pull out the JSON strings of the documents that are returned.
type SearchResults struct {
	Hits struct {
		Hits []Document
	}
}

// The CountResult struct is used to parse the response from Elasticsearch to
// a count request and pull out the number of documents that matched the query.
type CountResult struct {
	Count int
}

// The Document struct is used to parse Elasticsearch's JSON representation of
// a document.
type Document struct {
	Source map[string]interface{} `json:"_source"`
	ID     string                 `json:"_id"`
	Index  string                 `json:"_index"`
}

// GetAllDocuments finds and returns all the documents in an index. Any
// wildcards in the index name will be expanded.
func GetAllDocuments(c *elasticsearch.Client, i string) ([]Document, error) {
	// Check how many documents there are in the index
	resp, err := c.Count(c.Count.WithIndex(i), c.Count.WithExpandWildcards("all"))
	if err != nil {
		return nil, fmt.Errorf("Error getting number of documents in index %s: %s", i, err)
	}
	defer resp.Body.Close()

	var count CountResult
	err = json.Unmarshal([]byte(read(resp.Body)), &count)
	if err != nil {
		return nil, fmt.Errorf("Error decoding count result JSON string: %s", err)
	}

	// Get all the documents in the index
	resp, err = c.Search(c.Search.WithIndex(i), c.Search.WithExpandWildcards("all"), c.Search.WithSize(count.Count))
	if err != nil {
		return nil, fmt.Errorf("Error searching for documents for index %s: %s", i, err)
	}
	defer resp.Body.Close()

	// Decode JSON response into struct
	var docs SearchResults
	err = json.Unmarshal([]byte(read(resp.Body)), &docs)
	if err != nil {
		return nil, fmt.Errorf("Error decoding search results JSON string: %s", err)
	}

	return docs.Hits.Hits, nil
}

// GetDocument finds a single document from an index by ID.
func GetDocument(c *elasticsearch.Client, idx string, id string) (*Document, error) {
	resp, err := c.Get(idx, id)
	if err != nil {
		return nil, fmt.Errorf("Error getting document from index %s of id %s: %s", idx, id, err)
	}
	defer resp.Body.Close()

	// Decode JSON response into struct
	var doc *Document
	err = json.Unmarshal([]byte(read(resp.Body)), &doc)
	if err != nil {
		return nil, fmt.Errorf("Error decoding document JSON string: %s", err)
	}

	return doc, nil
}

func read(r io.Reader) string {
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}
