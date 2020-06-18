package keyrock

import (
	"encoding/json"
	"strings"
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

type application struct {
	ID          ID       `json:"id,omitempty"`
	Secret      string   `json:"secret,omitempty"`
	Image       string   `json:"image,omitempty"`
	JwtSecret   string   `json:"jwt_secret,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	RedirectURI string   `json:"redirect_uri,omitempty"`
	URL         string   `json:"url"`
	GrantType   *GrantType `json:"grant_type,omitempty"`
	TokenTypes  *TokenTypes `json:"token_types,omitempty"`
}

type TokenTypes struct {
	Types []string
}


// bearer is always returned back but is not valid
func (g *TokenTypes) UnmarshalJSON(data []byte) error {
	noQuotes := strings.ReplaceAll(string(data),"\"","")
	noSpaces := strings.ReplaceAll(noQuotes," ","")
	splitted := strings.Split(noSpaces,",")
	for index, tokenType := range splitted {
		if tokenType == "bearer" {
			splitted = append(splitted[:index], splitted[index+1:]...)
			break
		}
	}
	g.Types = splitted
	return nil
}
func (g *TokenTypes) MarshalJSON() ([]byte, error)  {
	if len(g.Types) == 0 {
		return nil, nil
	}
	return json.Marshal(g.Types)
}

type GrantType struct {
	Types []string
}

type ID struct {
	Value string
}
func (g *GrantType) UnmarshalJSON(data []byte) error {
	g.Types = strings.Split(strings.ReplaceAll(string(data),"\"",""),",")
	return nil
}
func (g *GrantType) MarshalJSON() ([]byte, error)  {
	if len(g.Types) == 0 {
		return nil, nil
	}
	return json.Marshal(g.Types)
}

func (i ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Value)
}
func (i *ID) UnmarshalJSON(data []byte) error {
	i.Value = strings.ReplaceAll(string(data),"\"","")
	return nil
}
