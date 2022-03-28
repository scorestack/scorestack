package esclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"go.uber.org/zap"
)

type Client struct {
	*elasticsearch.Client
}

func New(host string, username string, password string, verify bool) (*Client, error) {
	clientConfig := elasticsearch.Config{
		Addresses: []string{host},
		Username:  username,
		Password:  password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
			DialContext:         (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !verify,
			},
		},
	}
	es, err := elasticsearch.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %s", err)
	}

	return &Client{es}, nil
}

// The Document struct is used to parse Elasticsearch's JSON representation of
// a document.
type Document struct {
	Source map[string]interface{} `json:"_source"`
	ID     string                 `json:"_id"`
	Index  string                 `json:"_index"`
}

// TryLoad: return the value of a field in the document as a string. If the
// field is not present in the document, an error will be logged and an empty
// string will be returned.
func (d *Document) TryLoad(key string) string {
	if v, exists := d.Source[key]; !exists {
		zap.S().Errorf("expected field '%s' was not found in document document_id=%s index=%s", key, d.ID, d.Index)
		return ""
	} else {
		return v.(string)
	}
}

// GetAllDocumentsFrom finds and returns all the documents in the specified
// index. Any wildcards in the index name will be expanded.
func (c *Client) GetAllDocumentsFrom(index string) ([]Document, error) {
	// Check how many documents there are in the index
	resp, err := c.Count(c.Count.WithIndex(index), c.Count.WithExpandWildcards("all"))
	if err != nil {
		return nil, fmt.Errorf("Error getting number of documents in index %s: %s", index, err)
	}
	defer resp.Body.Close()

	count := struct {
		Count int
	}{}
	err = json.Unmarshal([]byte(read(resp.Body)), &count)
	if err != nil {
		return nil, fmt.Errorf("Error decoding count result JSON string: %s", err)
	}

	// Get all the documents in the index
	resp, err = c.Search(c.Search.WithIndex(index), c.Search.WithExpandWildcards("all"), c.Search.WithSize(count.Count))
	if err != nil {
		return nil, fmt.Errorf("Error searching for documents for index %s: %s", index, err)
	}
	defer resp.Body.Close()

	// Decode JSON response into struct
	docs := struct {
		Hits struct {
			Hits []Document
		}
	}{}
	err = json.Unmarshal([]byte(read(resp.Body)), &docs)
	if err != nil {
		return nil, fmt.Errorf("Error decoding search results JSON string: %s", err)
	}

	return docs.Hits.Hits, nil
}

// GetDocumentFrom finds a single document from the specified index using the
// document's ID.
func (c *Client) GetDocumentFrom(id string, index string) (*Document, error) {
	resp, err := c.Get(index, id)
	if err != nil {
		return nil, fmt.Errorf("Error getting document from index %s of id %s: %s", index, id, err)
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

func (c *Client) GetIndices(pattern string) ([]string, error) {
	resp, err := c.Indices.Get([]string{pattern})
	if err != nil {
		return nil, fmt.Errorf("failed to list indicies for pattern %s: ", pattern)
	}

	// Each top-level key of the response is the name of a matching index
	indices := make(map[string]interface{})
	err = json.Unmarshal([]byte(read(resp.Body)), &indices)
	if err != nil {
		return nil, fmt.Errorf("failed to decode index response as JSON string: %s", err)
	}
	index_names := make([]string, len(indices))
	i := 0
	for k := range indices {
		index_names[i] = k
	}

	return index_names, nil
}

func read(r io.Reader) string {
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}
