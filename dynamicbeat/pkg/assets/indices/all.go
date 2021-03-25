package indices

import (
	"io"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets"
)

func ResultsAdmin() io.Reader {
	return assets.Read("indices/results-admin.json")
}

func ResultsAll() io.Reader {
	return assets.Read("indices/results-all.json")
}

func ResultsTeam() io.Reader {
	return assets.Read("indices/results-team.json")
}
