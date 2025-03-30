package api

type AccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
