package keyrock

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)


const (
	TokenTypePermanent = "permanent"
	TokenTypeJWT       = "jwt"
	GrantTypePassword  = "password"
	GrantTypeAuthCode  = "authorization_code"
	GrantTypeImplicit  = "implicit"
)

// Password must be set
// Email must be set
// Base URL must be set
// Remove trailing slash from Base URL if exists
func validateOptions(options *Options) error {
	if len(options.Password) == 0 {
		return fmt.Errorf("password can not be empty")
	}
	if len(options.Email) == 0 {
		return fmt.Errorf("email can not be empty")
	}
	if len(options.BaseURL) == 0 {
		return fmt.Errorf("base url can not be empty")
	}
	if options.BaseURL[len(options.BaseURL)-1:] == "/" {
		options.BaseURL = options.BaseURL[:len(options.BaseURL)-1]
	}
	return nil
}

func (c *Client) WithToken(token string) *Client {
	if c.credentials == nil {
		c.credentials = &Credentials{
			Token: token,
		}
	} else {
		c.credentials.Token = token
	}
	return c
}

func (c Client) validateToken() error {
	if len(c.credentials.Token) == 0 {
		return fmt.Errorf("no Token available")
	}
	if time.Now().After(c.credentials.Valid.Add(time.Second * 10)) {
		if c.options.AutomaticTokenRefresh {
			return c.GetTokenWithPassword()
		}
		return TokenExpired{fmt.Errorf("Token expired since %f minutes",
			time.Now().Sub(c.credentials.Valid).Minutes())}
	}

	return nil
}
// Creates a new Keyrock Client
// Default timeout is 2 sec
func NewClient(options *Options) (*Client, error) {
	if err := validateOptions(options); err != nil {
		return nil, err
	}
	httpClient := &http.Client{}
	httpClient.Timeout = time.Second * 2
	newClient := &Client{
		httpClient: httpClient,
		options:    options,
		mutex:      &sync.Mutex{},
	}
	return newClient, nil
}
func (c *Client) SetTransport(transport *http.Transport) {
	c.httpClient.Transport = transport
}

func NewUser() *user {
	return &user{}
}
func NewApplication() *application {
	return &application{}
}


func (c Client) Ping() error {
	req, err := http.NewRequest("GET", c.getURL(""), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("keyrock not correct installed")
	}
	return nil
}
