package save

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"restapi/URL-Shortener/internal/lib/api/response"
	"restapi/URL-Shortener/internal/lib/logger/sl"
	"restapi/URL-Shortener/internal/lib/random"
	"restapi/URL-Shortener/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" valodate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

// cfg?
var aliasLength = 4

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		// TODO: Remove this log
		log.Debug("Debug request:", slog.String("Request", fmt.Sprint(req)))

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// Validation
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			// render.JSON(w, r, response.Error("invalide request"))
			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		// TODO: unique alias
		// TODO: Check if url is empty
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		// TODO: Change error
		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url aready exists", slog.String("url", req.URL))

			render.JSON(w, r, response.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, response.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		// TODO: Change to ResponseOK func
		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
