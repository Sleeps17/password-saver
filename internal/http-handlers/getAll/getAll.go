package getAll

import (
	"log/slog"
	"net/http"
	"password-saver/internal/lib/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Body
	Rows []string `json:"rows"`
}

type AllPasswordsGetter interface {
	GetAllPasswords() ([]string, error)
}

func New(logger *slog.Logger, worker AllPasswordsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		rows, err := worker.GetAllPasswords()
		if err != nil {
			logger.Error("Failed to get all passwords", slog.String("Error", err.Error()))
			render.JSON(w, r, response.Error("Failed to get all passwords"))
			return
		}

		logger.Info("All passwords got succesfully")
		render.JSON(w, r, Response{
			Body: response.OK(),
			Rows: rows,
		})
	}
}
