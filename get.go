package keyrock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (c *client) GetToken() error {
	body, err := json.Marshal(&user{
		Name:     c.options.Email,
		Password: c.options.Password,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.getURL("/v1/auth/tokens"), bytes.NewBuffer(body))
	req.Header.Set("Content-Type","application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return UnauthorizedError{error: fmt.Errorf("check your mail and/or password")}
	}

	type TokenResponse struct {
		Token struct {
			Methods   []string  `json:"methods"`
			ExpiresAt time.Time `json:"expires_at"`
		} `json:"token"`
	}
	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return err
	}
	token := resp.Header.Get("X-Subject-Token")

	if len(token) == 0 {
		return fmt.Errorf("Could not get token. Keyrock responsed with %s - code: %d ",
			resp.Status, resp.StatusCode)
	}
	c.credentials = &credentials{
		token: token,
		valid: tokenResponse.Token.ExpiresAt,
		methods: tokenResponse.Token.Methods,
	}
	return nil
}

func (c client) GetTokenInfo() (*TokenInfo, error) {
	if len(c.credentials.token) == 0 {
		return nil, fmt.Errorf("no token available")
	}
	req, err := http.NewRequest("GET",c.getURL("/v1/auth/tokens"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-token", c.credentials.token)
	req.Header.Set("X-Subject-token", c.credentials.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get informations")
	}
	var tokenInfo TokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("could not decode keyrock response - %s", err.Error())
	}
	return &tokenInfo, nil
}

func (c client) GetApplications() ([]*application, error) {
	if err := c.validateToken(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", c.getURL("/v1/applications"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-token", c.credentials.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get application information - status code %d\n", resp.StatusCode)
	}
	type ApplicationList struct {
		Applications []*application `json:"applications"`
	}

	var appList ApplicationList
	if err := json.NewDecoder(resp.Body).Decode(&appList); err != nil {
		return nil, err
	}
	return appList.Applications, nil
}

func (c client) GetApplication(id ID) (*application, error) {
	if err := c.validateToken(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", c.getURL(fmt.Sprintf("/v1/applications/%s", id.Value)), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-token", c.credentials.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var app application
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, err
	}
	return &app, nil
}
