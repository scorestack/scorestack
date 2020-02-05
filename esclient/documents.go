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

// The Document struct is used to parse Elasticsearch's JSON representation of
// a document and pull out a JSON string that contains the document fields.
type Document struct {
	Source map[string]interface{} `json:"_source"`
}

// GetAllDocuments finds all the documents in an index and returns the JSON
// strings that represent them.
func GetAllDocuments(c *elasticsearch.Client, i string) ([][]byte, error) {
	resp, err := c.Search(c.Search.WithIndex(i))
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

	// Get the JSON strings for each document
	out := make([][]byte, 0)
	for _, val := range docs.Hits.Hits {
		doc, err := json.Marshal(val.Source)
		if err != nil {
			return nil, fmt.Errorf("Error decoding document JSON string: %s", err)
		}
		out = append(out, doc)
	}

	return out, nil
}

// GetDocument finds a single document from an index by ID and returns the JSON
// string that represents it.
func GetDocument(c *elasticsearch.Client, idx string, id string) ([]byte, error) {
	resp, err := c.Get(idx, id)
	if err != nil {
		return nil, fmt.Errorf("Error getting document from index %s of id %s: %s", idx, id, err)
	}
	defer resp.Body.Close()

	// Decode JSON response into struct
	var doc Document
	err = json.Unmarshal([]byte(read(resp.Body)), &doc)
	if err != nil {
		return nil, fmt.Errorf("Error decoding document JSON string: %s", err)
	}

	// Get the JSON string for the document
	out, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("Error decoding document JSON string: %s", err)
	}

	return out, nil
}

func read(r io.Reader) string {
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}
