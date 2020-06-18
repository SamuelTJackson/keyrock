package keyrock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Updates the given application
//
// Only the following properties are usable for an update:
//
// name, description, redirect URI, URL, grant types and token types
func (c client) UpdateApplication(app *application) error {
	if len(app.ID.Value)  == 0 {
		return fmt.Errorf("id can not be empty")
	}

	body, err := json.Marshal(struct {
		Application struct {
			Name        string   `json:"name"`
			Description string   `json:"description"`
			RedirectURI string   `json:"redirect_uri,omitempty"`
			URL         string   `json:"url"`
			GrantType   *GrantType `json:"grant_type,omitempty"`
			TokenTypes  *TokenTypes `json:"token_types,omitempty"`
		}`json:"application"`

	}{
		Application: struct {
			Name        string   `json:"name"`
			Description string   `json:"description"`
			RedirectURI string   `json:"redirect_uri,omitempty"`
			URL         string   `json:"url"`
			GrantType   *GrantType `json:"grant_type,omitempty"`
			TokenTypes  *TokenTypes `json:"token_types,omitempty"`
		}{
			Name: app.Name,
			Description: app.Description,
			RedirectURI: app.RedirectURI,
			URL: app.URL,
			GrantType: app.GrantType,
			TokenTypes: app.TokenTypes,
		},
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH",c.getURL(fmt.Sprintf("/v1/applications/%s", app.ID.Value)),
		bytes.NewBuffer(body))
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
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not update values")
	}
	return nil
}
