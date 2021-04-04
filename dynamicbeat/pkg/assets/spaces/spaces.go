package spaces

import (
	"io"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets"
)

func Scorestack() io.Reader {
	return assets.Read("spaces/scorestack.json")
}
