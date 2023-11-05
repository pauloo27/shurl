package link

import (
	"log/slog"
	"net/http"

	"github.com/pauloo27/shurl/internal/ctx"
)

type CreateLinkBody struct {
	Slug        string `json:"slug" validate:"required,min=3,max=20"`
	OriginalURL string `json:"original_url" validate:"required,http_url"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	body, ok := ctx.MustGetBody[CreateLinkBody](w, r)
	if !ok {
		return
	}

	slog.Debug("create api", "body", body)

	w.Write([]byte("created?"))
}
