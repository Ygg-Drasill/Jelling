package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Ygg-Drasill/Jelling/common/api"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type SessionTokenResponse struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

func (ctx *Context) HandleAccountRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody api.AccountRequest
		var response SessionTokenResponse

		err := json.NewDecoder(r.Body).Decode(&requestBody)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Warn("Failed to read request body", "error", err)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
		var userId int
		err = ctx.Db.QueryRow(
			"INSERT INTO users (name, hash) VALUES (?, ?) RETURNING id",
			requestBody.Username, hash).Scan(&userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sessionToken, err := ctx.NewToken(w, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Token = sessionToken
		response.UserId = userId
		response.Username = requestBody.Username
		ctx.Logger.Info("User registered successfully", "userId", userId)
		ctx.WriteJSON(w, response, http.StatusOK)
	}
}

func (ctx *Context) HandleAccountLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody api.AccountRequest
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Warn("Failed to read request body", "error", err)
			return
		}

		var valid bool

		//Try to authenticate with token
		token, withToken := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if withToken {
			passwordIncluded := len(requestBody.Password) > 0
			username := requestBody.Username
			err = ctx.Db.QueryRow(`WITH sessions as (
					SELECT name, token, expiry_date
					FROM user_sessions INNER JOIN users ON user_sessions.user_id = users.id
					WHERE name = ?) SELECT EXISTS (SELECT 1 FROM sessions where token = ? AND expiry_date > ?)`,
				username, token, time.Now().UnixMilli()).Scan(&valid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.Error("Failed to check user session", "error", err)
				return
			}

			if valid {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			if !valid && !passwordIncluded {
				ctx.Write(w, []byte("Invalid credentials, token may have expired"), http.StatusUnauthorized)
				return
			}
		}

		//Authenticate without token
		var storedHash string
		var userId int
		err = ctx.Db.QueryRow("SELECT hash, id FROM users WHERE name = ?", requestBody.Username).Scan(&storedHash, &userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.Error("Failed to authenticate user", "error", err)
		}

		hashError := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(requestBody.Password))
		valid = hashError == nil
		if hashError != nil {
			w.WriteHeader(http.StatusUnauthorized)
			ctx.Logger.Warn("Password did not match", "error", hashError)
			return
		}

		var response SessionTokenResponse
		response.Token, err = ctx.NewToken(w, userId)
		response.UserId = userId
		response.Username = requestBody.Username
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.Error("Failed to create session token", "error", err)
			return
		}
		ctx.WriteJSON(w, response, http.StatusOK)
	}
}

func (ctx *Context) NewToken(w http.ResponseWriter, userId int) (string, error) {
	sessionToken := uuid.New().String()
	res, err := ctx.Db.Exec(
		"INSERT INTO user_sessions (user_id, token, expiry_date) VALUES (?, ?, ?)",
		userId, sessionToken, time.Now().Add(time.Hour*24*7).UnixMilli())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "", err
	}
	n, err := res.RowsAffected()
	if err != nil || n == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return "", fmt.Errorf("failed to insert session token (%d rows affacted)", n)
	}

	return sessionToken, nil
}
