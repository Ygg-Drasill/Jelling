package handlers

import (
	"encoding/json"
	"github.com/Ygg-Drasill/Jelling/common/api"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type SessionTokenResponse struct {
	Token string `json:"token"`
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

		sessionToken := uuid.New().String()
		res, err := ctx.Db.Exec(
			"INSERT INTO user_sessions (user_id, token, expiry_date) VALUES (?, ?, ?)",
			userId, sessionToken, time.Now().Add(time.Hour*24*7))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.Error("Failed to insert user session", "error", err)
			return
		}

		if n, err := res.RowsAffected(); err != nil || n == 0 {
			ctx.Logger.Error("Failed to insert user session", "rows", n, "error", err)
		}

		response.Token = sessionToken
		responseBodyBuff, err := json.Marshal(&response)
		ctx.Write(w, responseBodyBuff, http.StatusNoContent)
	}
}

func (ctx *Context) HandleAccountLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}
