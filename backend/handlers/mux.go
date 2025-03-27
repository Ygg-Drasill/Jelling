package handlers

import "net/http"

func NewJellingMux(ctx *Context) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /api/v1/health", ctx.HandleHealth())

	return mux
}
