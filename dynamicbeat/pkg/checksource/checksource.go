package checksource

import "github.com/scorestack/scorestack/dynamicbeat/pkg/check"

type CheckSource interface {
	LoadAll() ([]check.Config, error)
	LoadCheck(id string) (check.Config, error)
}
