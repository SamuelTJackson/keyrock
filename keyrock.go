package keyrock

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
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
	credentials *credentials
}
func validateOptions(options *Options) error {
	if len(options.Password) == 0 {
		return fmt.Errorf("password can not be empty")
	}
	if len(options.Email) == 0{
		return fmt.Errorf("email can not be empty")
	}
	if len(options.BaseURL) == 0 {
		return fmt.Errorf("base url can not be empty")
	}
	if options.BaseURL[len(options.BaseURL) - 1:] == "/" {
		options.BaseURL = options.BaseURL[:len(options.BaseURL) - 1]
	}
	return nil
}
func (c client) validateToken() error {
	if len(c.credentials.token) == 0 {
		return fmt.Errorf("no token available")
	}
	if time.Now().After(c.credentials.valid.Add(time.Second * 10)) {
		return TokenExpired{fmt.Errorf("token expired since %f minutes",
			time.Now().Sub(c.credentials.valid).Minutes())}
	}

	return nil
}

func NewClient(options *Options) (*client,error) {
	if err := validateOptions(options); err != nil {
		return nil, err
	}
	httpClient := &http.Client{}
	httpClient.Timeout = time.Second * 2
	newClient := &client{
		httpClient: httpClient,
		options: options,
		mutex: &sync.Mutex{},
	}
	return newClient, nil
}
func (c *client) SetTransport(transport *http.Transport)  {
	c.httpClient.Transport = transport
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


func (c client) Ping() error {
	req, err := http.NewRequest("GET", c.getURL(""),nil)
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




