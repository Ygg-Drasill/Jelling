package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (ctx *Context) WithAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userId int
		var userName string

		//Try to authenticate with token
		var token string
		var withToken bool
		for _, cookie := range r.Cookies() {
			if cookie.Name == "session" {
				token = cookie.Value
				withToken = true
				fmt.Println(token)
			}
		}

		if !withToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var err error
		err = ctx.Db.QueryRow(`SELECT users.id as user_id, name
					FROM user_sessions INNER JOIN users ON user_sessions.user_id = users.id
					WHERE token = ? AND expiry_date > ?`,
			token, time.Now().UnixMilli()).Scan(&userId, &userName)

		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.Error("Failed to check user session", "error", err)
			return
		}

		//TODO: set context deadline
		authContext := r.Context()
		authContext = context.WithValue(authContext, "userId", userId)
		authContext = context.WithValue(authContext, "userName", userName)
		next.ServeHTTP(w, r.WithContext(authContext))
	})
}
