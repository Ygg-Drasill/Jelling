package api

import (
	"net/http"
	"time"
)

func NewSessionCookie(sessionToken string, expiryDate time.Time) *http.Cookie {
	return &http.Cookie{
		Path:     "/api/v1",
		Name:     "session",
		Value:    sessionToken,
		Expires:  expiryDate,
		HttpOnly: true,
		Secure:   false,
	}
}
