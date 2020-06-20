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
// name, description, redirect URI, URL, grant types and Token types
func (c Client) UpdateApplication(app *application) error {
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
	req.Header.Set("X-Auth-Token", c.credentials.Token)
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

// Assign a role to a user in an application
func (c Client) AssignRoleToUserInApp(roleID ID, userID ID, appID ID) error {
	uri := fmt.Sprintf("/v1/applications/%s/users/%s/roles/%s",appID.Value, userID.Value, roleID.Value)
	req, err := http.NewRequest("POST",c.getURL(uri), nil) // POST is used instead of PUT - POST gives error
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-Token", c.credentials.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("could not assign user to role")
	}
	return nil
}