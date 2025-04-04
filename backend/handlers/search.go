package handlers

import (
	"fmt"
	"github.com/Ygg-Drasill/Jelling/common/api"
	"net/http"
	"strings"
)

func (ctx *Context) HandleArticleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has("title") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		params := api.SearchParameters{
			SortBy: r.URL.Query().Get("sortBy"),
			Title:  r.URL.Query().Get("title"),
			Topics: strings.Split(r.URL.Query().Get("topics"), ";"),
		}

		var articlesFound int
		rows, err := ctx.Db.Query(`SELECT id FROM articles WHERE title = ?`, params.Title)

		for rows.Next() {
			rows.Scan()
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.Error("Failed to search articles", "error", err)
			return
		}

		_, err = w.Write([]byte(fmt.Sprintf("%d", articlesFound)))
	}
}
