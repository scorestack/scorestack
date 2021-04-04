package checksource

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
	"go.uber.org/zap"
)

func NewElasticsearch(host string, username string, password string, verify bool, index string) (*Elasticsearch, error) {
	c := elasticsearch.Config{
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

	es, err := elasticsearch.NewClient(c)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %s", err)
	}

	return &Elasticsearch{Client: *es, Index: index}, nil
}

type Elasticsearch struct {
	elasticsearch.Client
	Index string
}

// The Document struct is used to parse Elasticsearch's JSON representation of
// a document.
type Document struct {
	Source map[string]interface{} `json:"_source"`
	ID     string                 `json:"_id"`
	Index  string                 `json:"_index"`
}

// GetAllDocuments finds and returns all the documents in the configured index.
// Any wildcards in the index name will be expanded.
func (e *Elasticsearch) GetAllDocuments() ([]Document, error) {
	return e.GetAllDocumentsFrom(e.Index)
}

// GetAllDocumentsFrom finds and returns all the documents in the specified
// index. Any wildcards in the index name will be expanded.
func (e *Elasticsearch) GetAllDocumentsFrom(index string) ([]Document, error) {
	// Check how many documents there are in the index
	resp, err := e.Count(e.Count.WithIndex(index), e.Count.WithExpandWildcards("all"))
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
	resp, err = e.Search(e.Search.WithIndex(index), e.Search.WithExpandWildcards("all"), e.Search.WithSize(count.Count))
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

// GetDocument finds a single document from the configured index using the
// document's ID.
func (e *Elasticsearch) GetDocument(id string) (*Document, error) {
	return e.GetDocumentFrom(id, e.Index)
}

// GetDocumentFrom finds a single document from the specified index using the
// document's ID.
func (e *Elasticsearch) GetDocumentFrom(id string, index string) (*Document, error) {
	resp, err := e.Get(index, id)
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

func (e *Elasticsearch) GetIndices(pattern string) ([]string, error) {
	resp, err := e.Indices.Get([]string{pattern})
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

func (e *Elasticsearch) GetAllAttributes(pattern string) (map[string]map[string]string, error) {
	docs, err := e.GetAllDocumentsFrom(pattern)
	if err != nil {
		return nil, err
	}

	// Organize attributes by check ID
	attributes := make(map[string]map[string]string)
	for _, doc := range docs {
		// Decode each attribute in the document
		attrs := make(map[string]string)
		for k, v := range doc.Source {
			// Decode the value of the attribute
			attrs[k] = v.(string)
		}
		attributes[doc.ID] = attrs
	}

	return attributes, nil
}

func (e *Elasticsearch) GetAttributes(id string, index string) (map[string]string, error) {
	doc, err := e.GetDocumentFrom(id, index)
	if err != nil {
		return nil, err
	}

	// Decode each attribute in the document
	attrs := make(map[string]string)
	for k, v := range doc.Source {
		// Decode the value of the attribute
		attrs[k] = v.(string)
	}

	return attrs, nil
}

func (e *Elasticsearch) LoadAll() ([]check.Config, error) {
	// Track how long it takes to update check definitions
	start := time.Now()

	results := make([]check.Config, 0)

	// Get list of checks
	checks, err := e.GetAllDocuments()
	if err != nil {
		return nil, err
	}

	// Get admin and user attributes
	admin, err := e.GetAllAttributes("attrib_admin_*")
	if err != nil {
		return nil, err
	}
	user, err := e.GetAllAttributes("attrib_user_*")
	if err != nil {
		return nil, err
	}

	// Iterate over each check
	for _, doc := range checks {
		result, err := buildCheckConfig(&doc, admin[doc.ID], user[doc.ID])
		if err != nil {
			return nil, err
		}
		results = append(results, *result)
	}

	zap.S().Infof("loaded %d check definitions in %.2f seconds", len(results), time.Since(start).Seconds())
	return results, nil
}

func (e *Elasticsearch) LoadCheck(id string) (*check.Config, error) {
	// Get check document
	check, err := e.GetDocument(id)
	if err != nil {
		return nil, err
	}

	// Parse team ID from check ID
	s := strings.Split(id, "-")
	team := s[len(s)-1]

	// Get attribute documents
	admin, err := e.GetAttributes(id, fmt.Sprintf("attrib_admin_%s", team))
	if err != nil {
		return nil, err
	}
	user, err := e.GetAttributes(id, fmt.Sprintf("admin_user_%s", team))
	if err != nil {
		return nil, err
	}
	indices, err := e.GetIndices(fmt.Sprintf("attrib_*_%s", team))
	if err != nil {
		return nil, err
	}
	attributes := make([]Document, len(indices))
	for i, name := range indices {
		doc, err := e.GetDocumentFrom(id, name)
		if err != nil {
			return nil, err
		}
		attributes[i] = *doc
	}

	return buildCheckConfig(check, admin, user)
}

func buildCheckConfig(doc *Document, admin map[string]string, user map[string]string) (*check.Config, error) {
	// Encode the definition as a JSON string
	def, err := json.Marshal(doc.Source["definition"])
	if err != nil {
		return nil, fmt.Errorf("Error encoding definition for %s to JSON string: %s", doc.ID, err)
	}

	// Unpack check definition into CheckConfig struct
	c := &check.Config{
		Metadata: check.Metadata{
			ID:          doc.Source["id"].(string),
			Name:        doc.Source["name"].(string),
			Type:        doc.Source["type"].(string),
			Group:       doc.Source["group"].(string),
			ScoreWeight: int64(doc.Source["score_weight"].(float64)),
		},
		Definition: def,
		Attributes: check.Attributes{
			Admin: admin,
			User:  user,
		},
	}

	return c, nil
}

func read(r io.Reader) string {
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}
