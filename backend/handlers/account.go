package handlers

import (
	"fmt"
	"net/http"
)

func (ctx *Context) HandleAccountRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hello world")
		w.WriteHeader(http.StatusNoContent)
	}
}

func (ctx *Context) HandleAccountLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}
