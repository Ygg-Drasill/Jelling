package handlers

import "net/http"

func NewJellingMux(ctx *Context) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /api/v1/health", ctx.HandleHealth())

	mux.Handle("POST /api/v1/account/register", ctx.HandleAccountRegister())
	mux.Handle("POST /api/v1/account/auth", ctx.HandleAccountLogin())

	mux.Handle("POST /api/v1/article/upload", ctx.WithAuthentication(ctx.HandleRunestoneUpload()))
	mux.Handle("GET /api/v1/article/search", ctx.HandleArticleSearch())

	return mux
}
