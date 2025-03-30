package handlers

import (
	"net/http"
)

func (ctx *Context) Write(w http.ResponseWriter, bytes []byte, status int) {
	w.WriteHeader(status)
	n, err := w.Write(bytes)
	if err != nil {
		ctx.Logger.Warn("Failed to write response", "length", n, "error", err)
	}
}
