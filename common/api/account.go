package api

type AccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SessionTokenResponse struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
