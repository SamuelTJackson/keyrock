package keyrock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)


// Creates an application in Keyrock
//
// The following properties are required:
//
// RedirectUI, Name, Description, URL
// The following characteristics must not be specified:
//
// Secret, JwtSecret, Image, ID
func (c Client) CreateApplication(app *application) error {
	if err := c.validateToken(); err != nil {
		return err
	}
	if len(app.RedirectURI) == 0 {
		return fmt.Errorf("redirect URIS are required")
	}
	if len(app.Name) == 0 {
		return fmt.Errorf("name is required")
	}
	if len(app.Description) == 0 {
		return fmt.Errorf("description is required")
	}
	if len(app.URL) == 0 {
		return fmt.Errorf("URL is required")
	}
	if len(app.Secret) != 0 || len(app.JwtSecret) != 0 || len(app.Image) != 0 || len(app.ID.Value) != 0{
		return fmt.Errorf("only set name, description, redirect uri, grant types and Token types")
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
	req.Header.Set("X-Auth-Token", c.credentials.Token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("could not create application - response code %d\n",resp.StatusCode)
	}
	type applicationResponse struct {
		Application application `json:"application"`
	}

	var appResponse applicationResponse
	if err := json.NewDecoder(resp.Body).Decode(&appResponse); err != nil {
		return err
	}
	app.GrantType = appResponse.Application.GrantType
	app.ID = appResponse.Application.ID
	app.TokenTypes = appResponse.Application.TokenTypes
	app.Image = appResponse.Application.Image
	app.JwtSecret = appResponse.Application.JwtSecret
	app.Secret = appResponse.Application.Secret
	return nil
}

// Create a pep-proxy for the given application
//
// Returns the proxy id and password. Don't lose the password you
// can't get it again. If you forget it, you have to reset the password
func (c Client) CreatePepProxy(id ID) (*PepProxy, error) {
	if err := c.validateToken(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST",
		c.getURL(fmt.Sprintf("/v1/applications/%s/pep_proxies", id.Value)),nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("X-Auth-Token", c.credentials.Token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	type pepResponse struct {
		PepProxy `json:"pep_proxy"`
	}
	var pepProxy pepResponse
	if err := json.NewDecoder(resp.Body).Decode(&pepProxy); err != nil {
		return nil, err
	}

	return &pepProxy.PepProxy, nil

}

// Creates a new user
//
// You have to set the username, email and password other values are not allowed
func (c Client) CreateUser(newUser *user) error {
	if len(newUser.Username) == 0 {
		return fmt.Errorf("username is required")
	}
	if len(newUser.Email) == 0 {
		return fmt.Errorf("email is required")
	}
	if len(newUser.Password) == 0 {
		return fmt.Errorf("password is required")
	}
	if (len(newUser.ID.Value) + len(newUser.Image)) != 0 {
		return fmt.Errorf("only username, email and password are allowed")
	}
	type tempUser struct {
		User *user `json:"user"`
	}
	body, err := json.Marshal(&tempUser{User: newUser})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.getURL("/v1/users"),bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("X-Auth-Token", c.credentials.Token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("could not create user - status code: %d", resp.StatusCode)
	}
	var createdUser tempUser
	if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
		return err
	}
	newUser.Admin = createdUser.User.Admin
	newUser.DatePassword = createdUser.User.DatePassword
	newUser.Image = createdUser.User.Image
	newUser.Enabled = createdUser.User.Enabled
	newUser.Gravatar = createdUser.User.Gravatar
	newUser.ID = createdUser.User.ID
	return  nil
}
