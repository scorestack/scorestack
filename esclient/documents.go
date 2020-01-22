package esclient

import (
	"bytes"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch"
	"github.com/tidwall/gjson"
)

// GetAllDocuments : Returns a list of all documents in the specified index.
func GetAllDocuments(c *elasticsearch.Client, i string) ([]map[string]gjson.Result, error) {
	resp, err := c.Search(c.Search.WithIndex(i))
	if err != nil {
		return nil, fmt.Errorf("Error searching for documents for index %s: %s", i, err)
	}
	defer resp.Body.Close()
	docs := gjson.Get(read(resp.Body), "hits.hits.#._source").Array()

	// Unpack results
	out := make([]map[string]gjson.Result, 0)
	for _, doc := range docs {
		out = append(out, doc.Map())
	}

	return out, nil
}

// GetDocument : Returns a single specified document.
func GetDocument(c *elasticsearch.Client, idx string, id string) (map[string]gjson.Result, error) {
	resp, err := c.Get(idx, id)
	if err != nil {
		return nil, fmt.Errorf("Error getting document from index %s of id %s: %s", idx, id, err)
	}
	defer resp.Body.Close()
	doc := gjson.Get(read(resp.Body), "_source").Map()

	return doc, nil
}

func read(r io.Reader) string {
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}
