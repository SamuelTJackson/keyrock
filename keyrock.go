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
	if len(options.password) == 0 {
		return fmt.Errorf("password can not be empty")
	}
	if len(options.email) == 0{
		return fmt.Errorf("email can not be empty")
	}
	if options.baseURL[len(options.baseURL) - 1:] == "/" {
		options.baseURL = options.baseURL[:len(options.baseURL) - 1]
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
	if err := newClient.refreshToken(); err != nil {
		return nil, err
	}
	return newClient, nil
}

func (c *client) refreshToken() error {
	body, err := json.Marshal(&user{
		Name:     c.options.email,
		Password: c.options.password,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.getURL("/v1/auth/tokens"), bytes.NewBuffer(body))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	c.mutex.Lock()
	c.token = resp.Header.Get("X-Subject-Token")
	c.mutex.Unlock()
	return nil
}

func (c *client) getToken() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.token
}
