package check

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Result struct {
	Metadata
	Timestamp time.Time
	Passed    bool
	Message   string
	Details   map[string]string
}

type generic struct {
	Metadata
	Timestamp string `json:"@timestamp"`
	Passed    bool   `json:"passed"`
	PassedInt uint8  `json:"passed_int"`
	Epoch     int64  `json:"epoch"`
}

func newGeneric(r *Result) generic {
	out := generic{
		Metadata:  r.Metadata,
		Timestamp: r.Timestamp.Format(time.RFC3339),
		Passed:    r.Passed,
		PassedInt: 0,
		Epoch:     r.Timestamp.Unix(),
	}

	if r.Passed {
		out.PassedInt = 1
	}

	return out
}

type full struct {
	generic
	Message string            `json:"message"`
	Details map[string]string `json:"details"`
}

func newFull(r *Result) full {
	return full{
		generic: newGeneric(r),
		Message: r.Message,
		Details: r.Details,
	}
}

func marshalError(err error) (string, io.Reader, error) {
	return "", nil, fmt.Errorf("failed to marshal event to JSON: %s", err)
}

func ok(index string, body []byte) (string, io.Reader, error) {
	return index, bytes.NewReader(body), nil
}

// Generic creates a JSON blob containing a check result without an error
// message or details field, as well as a destination index name for the check
// result document. The check results generated by this function can be used
// for visualizations that all teams will be able to see.
func (r *Result) Generic() (string, io.Reader, error) {
	body, err := json.Marshal(newGeneric(r))
	if err != nil {
		return marshalError(err)
	}

	return ok("results-all", body)
}

// Team creates a JSON blob containing a check result and the destination index
// name for the check result document. The check results generated by this
// function are for use by a single team only, and can therefore be used to
// provide more in-depth feedback to a team about why their checks may be
// failing.
func (r *Result) Team() (string, io.Reader, error) {
	body, err := json.Marshal(newFull(r))
	if err != nil {
		return marshalError(err)
	}

	return ok(fmt.Sprintf("results-%s", r.Group), body)
}

// Admin creates a JSON blob containing a check result and the destination
// index name for the check result document. The check results generated by
// this function are identical to the results generated for each team, but they
// are aggregated in a single index. This is intended to make it easier for
// Scorestack administrators to quickly view all available check result
// information to debug check definition problems, infrastructure issues, or
// anything else that might go awry during a competition.
func (r *Result) Admin() (string, io.Reader, error) {
	body, err := json.Marshal(newFull(r))
	if err != nil {
		return marshalError(err)
	}

	return ok("results-admin", body)
}
