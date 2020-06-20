package keyrock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (c *Client) GetTokenWithPassword() error {
	type user struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
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
		} `json:"Token"`
	}
	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return err
	}
	token := resp.Header.Get("X-Subject-Token")

	if len(token) == 0 {
		return fmt.Errorf("Could not get Token. Keyrock responsed with %s - code: %d ",
			resp.Status, resp.StatusCode)
	}
	c.credentials = &Credentials{
		Token:   token,
		Valid:   tokenResponse.Token.ExpiresAt,
		methods: tokenResponse.Token.Methods,
	}
	return nil
}

func (c *Client) GetTokenWithToken(token string) error {
	type tokenRequest struct {
		Token     string `json:"token"`
	}
	body, err := json.Marshal(&tokenRequest{
		Token: token,
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
	if err := checkError(resp); err != nil {
		return err
	}

	type TokenResponse struct {
		Token struct {
			Methods   []string  `json:"methods"`
			ExpiresAt time.Time `json:"expires_at"`
		} `json:"Token"`
	}
	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return err
	}
	authToken := resp.Header.Get("X-Subject-Token")

	if len(authToken) == 0 {
		return fmt.Errorf("Could not get Token. Keyrock responsed with %s - code: %d ",
			resp.Status, resp.StatusCode)
	}
	c.credentials = &Credentials{
		Token:   authToken,
		Valid:   tokenResponse.Token.ExpiresAt,
		methods: tokenResponse.Token.Methods,
	}
	return nil
}

func (c Client) GetTokenInfo() (*TokenInfo, error) {
	if c.credentials == nil || len(c.credentials.Token) == 0 {
		return nil, fmt.Errorf("no Token available")
	}
	req, err := http.NewRequest("GET",c.getURL("/v1/auth/tokens"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", c.credentials.Token)
	req.Header.Set("X-Subject-Token", c.credentials.Token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkError(resp); err != nil {
		return nil, err
	}

	var tokenInfo TokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("could not decode keyrock response - %s", err.Error())
	}
	return &tokenInfo, nil
}

func (c Client) GetApplications() ([]*application, error) {
	if err := c.validateToken(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", c.getURL("/v1/applications"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", c.credentials.Token)
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

func (c Client) GetApplication(id ID) (*application, error) {
	if err := c.validateToken(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", c.getURL(fmt.Sprintf("/v1/applications/%s", id.Value)), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", c.credentials.Token)
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

func (c Client) ListRoles(id ID) ([]*Role, error) {
	req, err := http.NewRequest("GET",c.getURL(fmt.Sprintf("/v1/applications/%s/roles",id.Value)),nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", c.credentials.Token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	type roleList struct {
		Roles []*Role `json:"roles"`
	}
	var roles roleList
	if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get roles - status code: %d", resp.StatusCode)
	}

	return roles.Roles,nil
}
func (c HTTPClient) GetUserInformation(token string) (*Userinformation, error){
	resp, err := http.Get(fmt.Sprintf("%s/user?access_token=%s",c.KeyrockBaseURL, token))
	if err != nil {
		return nil, fmt.Errorf("could not get token from keyrock - error: %s" , err)
	}
	defer resp.Body.Close()
	if err := checkError(resp); err != nil {
		return nil, err
	}
	var userInformation Userinformation
	if err := json.NewDecoder(resp.Body).Decode(&userInformation); err != nil {
		return nil, err
	}

	return &userInformation, nil
}