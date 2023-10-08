package save

import (
	"log/slog"
	"net/http"
	"password-saver/internal/lib/random"
	"password-saver/internal/lib/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	URL      string `json:"url"`
	Alias    string `json:"alias"`
	Password string `json:"password"`
}

type Response struct {
	response.Body
	Alias    string `json:"alias"`
	Password string `json:"password"`
}

type PasswordSaver interface {
	SavePassword(string, string, string) error
}

const RandomPasswordLenght = 10

func New(logger *slog.Logger, passwordSaver PasswordSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = slog.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logger.Error("Failed to decode request body: %w", err)
			render.JSON(w, r, response.Error("Failed to decode request body"))
			return
		}

		pasword := req.Password
		alias := req.Alias
		if pasword == "" {
			pasword = random.GeneratePassword(RandomPasswordLenght)
		}
		if alias == "" {
			generatedAlias, err := random.GenerateAlias(req.URL)
			if err != nil {
				logger.Error("Failed to generate alias", slog.String("Error", err.Error()))
				alias = req.URL
			} else {
				alias = generatedAlias
			}
		}

		if err := passwordSaver.SavePassword(req.URL, alias, pasword); err != nil {
			logger.Error("Failed to save password", slog.String("err", err.Error()))
			render.JSON(w, r, response.Error("Failed to save password"))
			return
		}

		logger.Info("Password saved successfully")

		render.JSON(w, r, Response{
			Body:     response.OK(),
			Alias:    alias,
			Password: pasword,
		})
	}
}
