package keyrock

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
