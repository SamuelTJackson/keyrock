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
		ID           string    `json:"id"`
		Username     string    `json:"username"`
		Email        string    `json:"email"`
		DatePassword time.Time `json:"date_password"`
		Enabled      bool      `json:"enabled"`
		Admin        bool      `json:"admin"`
	} `json:"User"`
}
