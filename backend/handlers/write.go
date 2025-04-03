package handlers

import (
	"encoding/json"
	"net/http"
)

func (ctx *Context) Write(w http.ResponseWriter, bytes []byte, status int) {
	w.WriteHeader(status)
	n, err := w.Write(bytes)
	if err != nil {
		ctx.Logger.Warn("Failed to write response", "length", n, "error", err)
	}
}

func (ctx *Context) WriteJSON(w http.ResponseWriter, response any, status int) {
	responseBodyBuff, err := json.Marshal(&response)
	if err != nil {
		ctx.Logger.Warn("Failed to marshal response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.Write(w, responseBodyBuff, status)
}
