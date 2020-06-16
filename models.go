package keyrock

import "time"

type user struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Token struct {
	Token string
}

type UnauthorizedError struct {
	error
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
	Applications []Application `json:"applications"`
}

type Application struct {
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
}
