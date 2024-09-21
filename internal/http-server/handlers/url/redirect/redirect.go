package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"restapi/URL-Shortener/internal/lib/api/response"
	"restapi/URL-Shortener/internal/lib/logger/sl"
	"restapi/URL-Shortener/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("reqquest_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Error("alias is empty")

			render.JSON(w, r, response.Error("not request"))

			return
		}

		url, err := urlGetter.GetURL(alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Error("url not found", "alias", alias)

			render.JSON(w, r, response.Error("not found"))
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, response.Error("failed to get url"))

			return
		}

		log.Info("got url", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
