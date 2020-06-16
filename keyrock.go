package keyrock

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Options struct {
	baseURL 	string
	email		string
	password	string
}

type client struct {
	httpClient *http.Client
	options *Options
	token string
}

func NewClient(options *Options) (*client,error) {
	httpClient := &http.Client{}
	newClient := &client{
		httpClient: httpClient,
		options: options,
	}
	if err := newClient.getToken(); err != nil {
		return nil, err
	}
	return newClient, nil
}

func (c *client) getToken() error {
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
	c.token = resp.Header.Get("X-Subject-Token")
	return nil
}
