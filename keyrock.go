package keyrock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type Options struct {
	BaseURL 	string
	Email		string
	Password	string
}

type client struct {
	httpClient *http.Client
	options *Options
	mutex *sync.Mutex
}



func validateOptions(options *Options) error {
	if len(options.Password) == 0 {
		return fmt.Errorf("password can not be empty")
	}
	if len(options.Email) == 0{
		return fmt.Errorf("email can not be empty")
	}
	if options.BaseURL[len(options.BaseURL) - 1:] == "/" {
		options.BaseURL = options.BaseURL[:len(options.BaseURL) - 1]
	}
	return nil
}

func NewClient(options *Options) (*client,error) {
	if err := validateOptions(options); err != nil {
		return nil, err
	}
	httpClient := &http.Client{}
	newClient := &client{
		httpClient: httpClient,
		options: options,
		mutex: &sync.Mutex{},
	}
	return newClient, nil
}

func (c client) ListApplications(token *Token) (*ApplicationList, error) {
	req, err := http.NewRequest("GET", c.getURL("/v1/applications"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-token", token.Token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var appList ApplicationList
	if err := json.NewDecoder(resp.Body).Decode(&appList); err != nil {
		return nil, err
	}
	return &appList, nil
}
func NewApplication(name string, description string, url string) *application{
	return &application{
		Name:        name,
		Description: description,
		URL:         url,
	}
}
func (a *application) WithRedirectURIS(uri ...string) *application {
	a.RedirectURI = strings.Join(uri,",")
	return a
}
func (a *application) WithGrantTypes(types ...string) *application {
	a.GrantType = types
	return a
}
func (a *application) WithTokenTypes(types ...string) *application {
	a.TokenTypes = types
	return a
}


func (c client) CreateApplication(token *Token, app *application)  error {
	if len(app.Secret) != 0 || len(app.JwtSecret) != 0 || len(app.Image) != 0 || len(app.ID.Value) != 0{
		return fmt.Errorf("only set name, description, redirect uri, grant types and token types")
	}
	body, err := json.Marshal(struct {
		*application `json:"application"`
	}{app})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.getURL("/v1/applications"), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-token", token.Token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var appResponse applicationResponse
	if err := json.NewDecoder(resp.Body).Decode(&appResponse); err != nil {
		return err
	}
	app.GrantType = strings.Split(appResponse.Application.GrantType,",")
	app.ID.Value = appResponse.Application.ID
	app.TokenTypes = strings.Split(appResponse.Application.TokenTypes,",")
	app.Image = appResponse.Application.Image
	app.JwtSecret = appResponse.Application.JwtSecret
	app.Secret = appResponse.Application.Secret
	return nil
}

func (c client) DeleteApplication(token *Token, id ID) error {
	if len(id.Value) == 0 {
		return fmt.Errorf("id can not be empty")
	}
	req, err := http.NewRequest("DELETE", c.getURL("/v1/applications/") + id.Value,nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-token", token.Token)
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

func (c client) GetTokenInfo(token *Token) (*TokenInfo, error) {
	req, err := http.NewRequest("GET",c.getURL("/v1/auth/tokens"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-token", token.Token)
	req.Header.Set("X-Subject-token", token.Token)
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

func (c client) GetToken() (*Token, error) {
	body, err := json.Marshal(&user{
		Name:     c.options.Email,
		Password: c.options.Password,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.getURL("/v1/auth/tokens"), bytes.NewBuffer(body))
	req.Header.Set("Content-Type","application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, UnauthorizedError{error: fmt.Errorf("check your mail and/or password")}
	}
	if token := resp.Header.Get("X-Subject-Token"); len(token) != 0 {
		return &Token{
			Token: token,
		}, nil
	}
	return nil, fmt.Errorf("Could not get token. Keyrock responsed with %s - code: %d ",
		resp.Status, resp.StatusCode)
}

