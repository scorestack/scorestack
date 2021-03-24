package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// INDEX_DATE is a layout string for time.Time.Format that is compatible with
// the Elastic Stack convention of indexes named "index-YYYY.MM.DD".
const INDEX_DATE = "2006.01.02"

// TIMESTAMP_DATE is the layout string used by Elastic Beats for the @timestamp
// field.
// See for more details:
// https://www.elastic.co/guide/en/beats/filebeat/current/processor-timestamp.html
const TIMESTAMP_DATE = "Mon Jan 2 15:04:05 MST 2006"

type Event struct {
	Timestamp   time.Time
	Id          string
	Name        string
	CheckType   string
	Group       string
	ScoreWeight float64
	Passed      bool
	Message     string
	Details     map[string]string
}

type full struct {
	Timestamp   string            `json:"@timestamp"`
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	CheckType   string            `json:"type"`
	Group       string            `json:"group"`
	ScoreWeight float64           `json:"score_weight"`
	Passed      bool              `json:"passed"`
	PassedInt   uint8             `json:"passed_int"`
	Epoch       int64             `json:"epoch"`
	Message     string            `json:"message"`
	Details     map[string]string `json:"details"`
}

type generic struct {
	Timestamp   string  `json:"@timestamp"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	CheckType   string  `json:"type"`
	Group       string  `json:"group"`
	ScoreWeight float64 `json:"score_weight"`
	Passed      bool    `json:"passed"`
	PassedInt   uint8   `json:"passed_int"`
	Epoch       int64   `json:"epoch"`
}

func Admin(e Event) (string, io.Reader, error) {
	f := &full{
		Timestamp:   e.Timestamp.Format(TIMESTAMP_DATE),
		ID:          e.Id,
		Name:        e.Name,
		CheckType:   e.CheckType,
		Group:       e.Group,
		ScoreWeight: e.ScoreWeight,
		Passed:      e.Passed,
		Message:     e.Message,
		Details:     e.Details,
		Epoch:       e.Timestamp.Unix(),
	}

	if e.Passed {
		f.PassedInt = 1
	} else {
		f.PassedInt = 0
	}

	body, err := json.Marshal(f)
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal event to JSON: %s", err)
	}

	return fmt.Sprintf("results-admin-%s", e.Timestamp.Format(INDEX_DATE)), bytes.NewReader(body), nil
}

func Team(e Event) (string, io.Reader, error) {
	_, reader, err := Admin(e)
	if err != nil {
		return "", nil, err
	}

	return fmt.Sprintf("results-%s-%s", e.Group, e.Timestamp.Format(INDEX_DATE)), reader, nil
}

func Generic(e Event) (string, io.Reader, error) {
	g := &generic{
		Timestamp:   e.Timestamp.Format(TIMESTAMP_DATE),
		ID:          e.Id,
		Name:        e.Name,
		CheckType:   e.CheckType,
		Group:       e.Group,
		ScoreWeight: e.ScoreWeight,
		Passed:      e.Passed,
		Epoch:       e.Timestamp.Unix(),
	}

	if e.Passed {
		g.PassedInt = 1
	} else {
		g.PassedInt = 0
	}

	body, err := json.Marshal(g)
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal event to JSON: %s", err)
	}

	return fmt.Sprintf("results-all-%s", e.Timestamp.Format(INDEX_DATE)), bytes.NewReader(body), nil
}
