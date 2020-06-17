package keyrock

import (
	"fmt"
	"net/http"
)

func (c client) DeleteApplication(id ID) error {
	if err := c.validateToken(); err != nil {
		return err
	}
	if len(id.Value) == 0 {
		return fmt.Errorf("id can not be empty")
	}
	req, err := http.NewRequest("DELETE", c.getURL("/v1/applications/") + id.Value,nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-token", c.credentials.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not delete application")
	}
	return nil
}
