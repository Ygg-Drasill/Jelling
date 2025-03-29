package handlers

import "net/http"

func NewJellingMux(ctx *Context) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /api/v1/health", ctx.HandleHealth())

	mux.HandleFunc("POST /api/v1/account/register", ctx.HandleAccountRegister())
	mux.HandleFunc("POST /api/v1/account/auth", ctx.HandleAccountLogin())

	return mux
}
