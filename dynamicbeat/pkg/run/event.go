package run

import "time"

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
