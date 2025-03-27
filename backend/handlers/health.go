package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type JellingHealth struct {
	Status              int           `json:"status"`
	DbWaitTime          time.Duration `json:"dbWaitTime"`
	DbConnectionsActive int           `json:"dbConnectionsActive"`
	DbConnectionsIdle   int           `json:"dbConnectionsIdle"`
}

func (ctx *Context) HandleHealth() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stats := ctx.Db.Stats()

		response := JellingHealth{
			Status:              0,
			DbWaitTime:          stats.WaitDuration,
			DbConnectionsActive: stats.InUse,
			DbConnectionsIdle:   stats.Idle,
		}

		responseJson, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(responseJson)
	})
}
