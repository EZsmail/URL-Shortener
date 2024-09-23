package get

import (
	"errors"
	"log/slog"
	"net/http"
	"restapi/URL-Shortener/internal/lib/api/response"
	"restapi/URL-Shortener/internal/lib/logger/sl"
	"restapi/URL-Shortener/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	Alias string `json:"alias"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
	Url   string `json:"url"`
}

type URLGetter interface {
	GetURL(alias string) (string, error)
}

// TODO: rewrite handler

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.get.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body")

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// TODO: Validator dont work
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		var alias = req.Alias

		url, err := urlGetter.GetURL(alias)

		if url == "" {
			log.Debug("url not found", slog.String("alias", alias), slog.String("url", url))
		}

		// TODO: Change alias=="" to ErrAlias not found
		if errors.Is(err, storage.ErrAliasNotFound) {
			log.Info("alias not found", slog.String("alias", alias))

			render.JSON(w, r, response.Error("alias not found"))

			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, response.Error("failed to get url"))

			return
		}

		log.Info("url founded", slog.String("url", url))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
			Url:      url,
		})
	}
}
