package get

import (
	"net/http"
	"password-saver/internal/lib/response"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Key string `json:"key"`
}

type Response struct {
	response.Body
	Password string `json:"password"`
}

type PasswordGetter interface {
	GetPassword(string) (string, error)
}

func New(logger *slog.Logger, passwordGetter PasswordGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("Failed to decode request body", slog.String("Error", err.Error()))
			render.JSON(w, r, response.Error("Failed to decode request body"))
			return
		}

		if req.Key == "" {
			logger.Error("The key field is empty, although it should contain an alias or url")
			render.JSON(w, r, response.Error("The key field is empty, although it should contain an alias or url"))
			return
		}

		password, err := passwordGetter.GetPassword(req.Key)
		if err != nil {
			logger.Error("Failed to get password", slog.String("Error", err.Error()))
			render.JSON(w, r, response.Error("Failed to get password"))
			return
		}

		logger.Info("Password get successfully")

		render.JSON(w, r, Response{
			Body:     response.OK(),
			Password: password,
		})

	}
}
