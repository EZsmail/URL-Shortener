package delete

import (
	"log/slog"
	"net/http"
	"restapi/URL-Shortener/internal/lib/api/response"
	"restapi/URL-Shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	Alias string `json:"alias" validate:"omitempty"`
	Id    int    `json:"id" validate:"omitempty"`
}

type Response struct {
	response.Response
}

type URLdeleter interface {
	DelURLByAlias(alias string) (int, error)
}

func New(log *slog.Logger, urldeleter URLdeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.delete.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body")

			render.JSON(w, r, response.Error("failed to decode request body"))

		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		alias := req.Alias

		_, err = urldeleter.DelURLByAlias(alias)

		if err != nil {
			log.Error("failed to del url", sl.Err(err))

			render.JSON(w, r, response.Error("failed to del url"))

			return
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
		})
	}
}
