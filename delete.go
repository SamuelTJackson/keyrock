package keyrock

import (
	"fmt"
	"net/http"
)

// Delete an application in keyrock by the given id.
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

// Delete the stored token. You are not able to request data from keyrock without token
func (c *client) DeleteToken() error {
	if err := c.validateToken(); err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", c.getURL("/v1/auth/tokens"), nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-token",c.credentials.token)
	req.Header.Set("X-Subject-token",c.credentials.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not delete token")
	}
	c.credentials = nil

	return nil
}
