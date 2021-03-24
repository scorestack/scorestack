package check

import "context"

type Check interface {
	GetConfig() Config
	SetConfig(c Config)
	Run(ctx context.Context) Result
}

type Metadata struct {
	ID          string
	Name        string
	Type        string
	Group       string
	ScoreWeight float64
}

type Config struct {
	Meta       Metadata
	Definition []byte
	Attribs    map[string]string
}

type Result struct {
	Meta    Metadata
	Passed  bool
	Message string
	Details map[string]string
}
