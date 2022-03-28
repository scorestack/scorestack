package checksource

import (
	"github.com/scorestack/scorestack/dynamicbeat/pkg/models"
)

type CheckSource interface {
	LoadAll() ([]models.CheckConfig, error)
}
