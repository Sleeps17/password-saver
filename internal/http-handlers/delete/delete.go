package delete

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
}

type PasswordDeleter interface {
	DeletePassword(string) error
}

func New(logger *slog.Logger, passwordDeleter PasswordDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logger = logger.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logger.Error("Failed to decode request body", slog.String("Error", err.Error()))
			render.JSON(w, r, response.Error("Failed to decode request body"))
			return
		}

		if req.Key == "" {
			logger.Error("The key field is empty, but must contain an alias or url")
			render.JSON(w, r, response.Error("The key field is empty, but must contain an alias or url"))
			return
		}

		if err := passwordDeleter.DeletePassword(req.Key); err != nil {
			logger.Error("Failed to delete password", slog.String("Error", err.Error()))
			render.JSON(w, r, response.Error("Failed to delete password"))
			return
		}

		logger.Info("Password deleted succsecfully")

		render.JSON(w, r, Response{
			Body: response.OK(),
		})
	}
}
