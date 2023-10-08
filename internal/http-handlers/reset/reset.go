package reset

import (
	"log/slog"
	"net/http"
	"password-saver/internal/lib/random"
	"password-saver/internal/lib/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Key         string `json:"key"`
	NewPassword string `json:"new_password"`
}

type Response struct {
	response.Body
	Key      string `json:"key"`
	Password string `json:"password"`
}

type PasswordReseter interface {
	ResetPassword(string, string) error
}

func New(logger *slog.Logger, passwordReseter PasswordReseter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logger = logger.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logger.Error("Failed to decode request body", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("Failed to decode request body"))
			return
		}

		if req.Key == "" {
			logger.Error("The key field is empty, but must contain an alias or url")
			render.JSON(w, r, response.Error("The key field is empty, but must contain an alias or url"))
			return
		}

		password := req.NewPassword
		if password == "" {
			password = random.GeneratePassword(RandomPasswordLenght)
		}

		if err := passwordReseter.ResetPassword(req.Key, password); err != nil {
			logger.Error("Failed to change password", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("Failed to change password"))
			return
		}

		logger.Info("Password changed successfully")
		render.JSON(w, r, Response{
			Body:     response.OK(),
			Key:      req.Key,
			Password: password,
		})
	}
}

const RandomPasswordLenght = 10
