package keyrock

import (
	"fmt"
	"net/http"
)

// Delete an application in keyrock by the given id.
func (c Client) DeleteApplication(id ID) error {
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
	req.Header.Set("X-Auth-Token", c.credentials.Token)
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

// Delete the stored Token. You are not able to request data from keyrock without Token
func (c *Client) DeleteToken() error {
	if err := c.validateToken(); err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", c.getURL("/v1/auth/tokens"), nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-Token",c.credentials.Token)
	req.Header.Set("X-Subject-Token",c.credentials.Token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not delete Token")
	}
	c.credentials = nil

	return nil
}
