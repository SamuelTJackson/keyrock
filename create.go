package keyrock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)



func (c client) CreateApplication(app *application) error {
	if err := c.validateToken(); err != nil {
		return err
	}
	if len(app.RedirectURI) == 0 {
		return fmt.Errorf("redirect URIS are required")
	}
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
	req.Header.Set("X-Auth-token", c.credentials.token)
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


func (c client) CreatePepProxy(id ID) (*PepProxy, error) {
	if err := c.validateToken(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST",
		c.getURL(fmt.Sprintf("/v1/applications/%s/pep_proxies", id.Value)),nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("X-Auth-token", c.credentials.token)
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