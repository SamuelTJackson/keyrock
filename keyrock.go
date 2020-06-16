package keyrock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	token string
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

func (c *client) ListApplications() {

}

func (c *client) GetToken() (*Token, error) {
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

