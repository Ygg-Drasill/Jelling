package handlers

import (
	"fmt"
	"github.com/Ygg-Drasill/Jelling/common/contentType"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

func (ctx *Context) HandleRunestoneUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorId := r.Context().Value("userId").(int)
		articleType := contentType.PLAIN
		var articleTitle string
		var articleSummary string
		var articleData []byte

		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			ctx.Logger.Debug("Bad request received, expected multipart/form-data")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		partReader, err := r.MultipartReader()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Debug("Failed to create multipart form reader", "error", err)
			return
		}

		var formPart *multipart.Part
		totalRowsAffected := int64(0)
		for formPart, err = partReader.NextPart(); err != io.EOF; formPart, err = partReader.NextPart() {
			formPartField := formPart.FormName()
			switch formPartField {
			case "data":
				articleData, err = io.ReadAll(formPart)
				break
			case "title":
				value, err := io.ReadAll(formPart)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					ctx.Logger.Error("Failed to read article title", "error", err)
					break
				}
				articleTitle = string(value)
				break
			case "summary":
				value, err := io.ReadAll(formPart)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					ctx.Logger.Error("Failed to read article summary", "error", err)
					break
				}
				articleSummary = string(value)
				break
			}

			err = formPart.Close()
			if err != nil {
				ctx.Logger.Error(err.Error())
			}
		}

		if len(articleData) == 0 {
			ctx.Logger.Warn("No bytes read from article")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var articleId int
		err = ctx.Db.QueryRow(`INSERT INTO articles (author_id, title, summary, content_type)
						VALUES (?,?,?,?) RETURNING id`, authorId, articleTitle, articleSummary, articleType).Scan(&articleId)

		if err != nil {
			ctx.Logger.Debug("Failed to insert into articles", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		res, err := ctx.Db.Exec(
			"INSERT INTO article_blobs (article_id, data) VALUES (?, ?)",
			articleId, articleData)
		if err != nil {
			ctx.Logger.Error("Failed to insert blob", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		blobId, err := res.LastInsertId()
		if err != nil {
			ctx.Logger.Error("Failed to get last blob id", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rows, err := res.RowsAffected()
		if err != nil {
			ctx.Logger.Error("Failed to get rows affected", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if rows == 0 {
			ctx.Logger.Info("No rows affected on blob upload", "rows", rows)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		totalRowsAffected += rows
		ctx.Logger.Debug("Uploaded blob", "blobId", blobId)

		response := fmt.Sprintf("{\"blobs\": %d}", totalRowsAffected)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(response))
		if err != nil {
			ctx.Logger.Error("Failed to write response", "error", err.Error())
		}
	}
}
