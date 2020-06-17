package keyrock

import (
	"encoding/json"
	"time"
)

type user struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type credentials struct {
	token   string
	valid   time.Time
	methods []string
}

type UnauthorizedError struct {
	error
}
type TokenExpired struct {
	error
}

type PepProxy struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type TokenInfo struct {
	AccessToken string    `json:"access_token"`
	Expires     time.Time `json:"expires"`
	Valid       bool      `json:"valid"`
	User        struct {
		Scope        []interface{} `json:"scope"`
		ID           string        `json:"id"`
		Username     string        `json:"username"`
		Email        string        `json:"email"`
		DatePassword time.Time     `json:"date_password"`
		Enabled      bool          `json:"enabled"`
		Admin        bool          `json:"admin"`
	} `json:"User"`
}

type ApplicationList struct {
	Applications []struct {
		ID           string      `json:"id"`
		Name         string      `json:"name"`
		Description  string      `json:"description"`
		Image        string      `json:"image"`
		URL          string      `json:"url"`
		RedirectURI  string      `json:"redirect_uri"`
		GrantType    string      `json:"grant_type"`
		ResponseType string      `json:"response_type"`
		TokenTypes   string      `json:"token_types"`
		ClientType   interface{} `json:"client_type"`
	} `json:"applications"`
}

type applicationResponse struct {
	Application struct {
		ID           string `json:"id"`
		Secret       string `json:"secret"`
		Image        string `json:"image"`
		JwtSecret    string `json:"jwt_secret"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		RedirectURI  string `json:"redirect_uri"`
		URL          string `json:"url"`
		GrantType    string `json:"grant_type"`
		TokenTypes   string `json:"token_types"`
		ResponseType string `json:"response_type"`
	} `json:"application"`
}

type application struct {
	ID          ID       `json:"id,omitempty"`
	Secret      string   `json:"secret,omitempty"`
	Image       string   `json:"image,omitempty"`
	JwtSecret   string   `json:"jwt_secret,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	RedirectURI string   `json:"redirect_uri,omitempty"`
	URL         string   `json:"url"`
	GrantType   []string `json:"grant_type,omitempty"`
	TokenTypes  []string `json:"token_types,omitempty"`
}

type ID struct {
	Value string `json:"id,omitempty"`
}

func (i ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Value)
}
func (i *ID) UnmarshalJSON(data []byte) error {

	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	i.Value, _ = v[0].(string)

	return nil
}
