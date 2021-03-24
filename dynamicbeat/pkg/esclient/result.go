package esclient

import (
	"fmt"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/event"
)

func Index(c *elasticsearch.Client, check event.Event) error {
	index, reader, err := event.Admin(check)
	if err != nil {
		return err
	}
	resp, err := c.Index(index, reader)
	if err != nil {
		return fmt.Errorf("failed to index admin result for %s: %s", check.Id, err)
	}
	if resp.IsError() {
		return fmt.Errorf("error indexing admin result document for %s: %s", check.Id, resp.Status())
	}
	defer resp.Body.Close()

	index, reader, err = event.Team(check)
	if err != nil {
		return err
	}
	resp, err = c.Index(index, reader)
	if err != nil {
		return fmt.Errorf("failed to index team result for %s: %s", check.Id, err)
	}
	if resp.IsError() {
		return fmt.Errorf("error indexing team result document for %s: %s", check.Id, resp.Status())
	}
	defer resp.Body.Close()

	index, reader, err = event.Generic(check)
	if err != nil {
		return err
	}
	resp, err = c.Index(index, reader)
	if err != nil {
		return fmt.Errorf("failed to index generic result for %s: %s", check.Id, err)
	}
	if resp.IsError() {
		return fmt.Errorf("error indexing generic result document for %s: %s", check.Id, resp.Status())
	}
	defer resp.Body.Close()

	return nil
}
